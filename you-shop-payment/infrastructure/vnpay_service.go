package infrastructure

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/TechwizsonORG/payment-service/config/model"
	"github.com/TechwizsonORG/payment-service/entity"
	"github.com/TechwizsonORG/payment-service/usecase"
	"github.com/TechwizsonORG/payment-service/usecase/rpc"
	"github.com/TechwizsonORG/payment-service/util"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type VnpayService struct {
	config         model.VnpayConfig
	rpc            rpc.RpcInterface
	rpcEndpoint    model.RpcServerEndpoint
	logger         zerolog.Logger
	paymentService usecase.Service
}

func (v *VnpayService) CreatePaymentUrlByOrderId(orderId uuid.UUID, ipAddress string) (string, error) {
	payment, getErr := v.paymentService.GetPaymentByOrderId(orderId)
	if getErr != nil {
		v.logger.Error().Err(getErr).Msg("")
		return "", getErr
	}
	return v.CreatePaymentUrl(*payment, ipAddress)
}

func NewVnpayService(config model.VnpayConfig, rpc rpc.RpcInterface, rpcEndpoint model.RpcServerEndpoint, logger zerolog.Logger) *VnpayService {
	return &VnpayService{
		config:      config,
		rpc:         rpc,
		rpcEndpoint: rpcEndpoint,
		logger:      logger,
	}
}

func (v *VnpayService) CreatePaymentUrl(payment entity.Payment, ipAddr string) (string, error) {

	params := map[string]string{
		"vnp_Version":    v.config.Version,
		"vnp_Command":    v.config.Command,
		"vnp_TmnCode":    v.config.MerchantCode,
		"vnp_Amount":     fmt.Sprintf("%d", int(payment.Amount*100)),
		"vnp_CurrCode":   v.config.Currency,
		"vnp_TxnRef":     payment.Id.String(),
		"vnp_OrderInfo":  fmt.Sprintf("Thanh_toa_don_hang_%s", payment.Id.String()),
		"vnp_Locale":     v.config.Locale,
		"vnp_ReturnUrl":  v.config.ReturnUrl,
		"vnp_IpAddr":     ipAddr,
		"vnp_OrderType":  v.config.OrderType,
		"vnp_CreateDate": util.GetCurrentUtcTime(7).Format("20060102150405"),
		"vnp_ExpireDate": util.GetCurrentUtcTime(7).Add(1 * time.Hour).Format("20060102150405"),
	}

	// Sort parameters by key
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build query string
	var queryParams []string
	for _, key := range keys {
		queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, url.QueryEscape(params[key])))
	}
	queryString := strings.Join(queryParams, "&")

	// Generate secure hash
	hash := hmac.New(sha512.New, []byte(v.config.SecretKey))
	hash.Write([]byte(queryString))
	// secureHash := hex.EncodeToString(hash.Sum(nil))
	secureHash := hex.EncodeToString(hash.Sum(nil))
	v.logger.Debug().Str("query_string", queryString).Str("hash_value", secureHash).Msg("")
	paymentURL := fmt.Sprintf("%s?%s&vnp_SecureHash=%s", v.config.PageUrl, queryString, secureHash)
	v.logger.Debug().Str("Payment URL", paymentURL).Msg("")
	return paymentURL, nil
}

func (v *VnpayService) VnpPayHandleCallback(request *http.Request) (paymentId *uuid.UUID, transactionStatus entity.TransactionStatus, err error) {
	// Verify secure hash
	queryParams := request.URL.Query()
	secureHash := queryParams.Get("vnp_SecureHash")
	queryParams.Del("vnp_SecureHash")
	queryParams.Del("vnp_SecureHashType")

	// Sort and build query string
	var keys []string
	for k := range queryParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var queryParamBuilder strings.Builder
	for i, k := range keys {
		if i > 0 {
			queryParamBuilder.WriteByte('&')
		}
		queryParamBuilder.WriteString(fmt.Sprintf("%s=%s", k, queryParams.Get(k)))
	}
	params := queryParamBuilder.String()

	v.logger.Debug().Str("params", params).Msg("")
	// Generate expected secure hash
	hash := hmac.New(sha512.New, []byte(v.config.SecretKey))
	hash.Write([]byte(params))
	expectedSecureHash := hex.EncodeToString(hash.Sum(nil))

	if secureHash != expectedSecureHash {
		v.logger.Error().Str("secure_hash", secureHash).Str("expected_secure_hash", expectedSecureHash).Msg("Failed to verify secure hash")
		return nil, entity.TransciontFailed, errors.New("Secure hash failed")
	}

	parsedPaymentId, parseErr := uuid.Parse(queryParams.Get("vnp_TxnRef"))
	if parseErr != nil {
		v.logger.Error().Err(parseErr).Msg("Failed to parse payment id")
		return nil, entity.TransciontFailed, errors.New("Invalid payment id")
	}

	// Check response code
	responseCode := queryParams["vnp_ResponseCode"][0]
	if responseCode == "00" {
		return &parsedPaymentId, entity.TransactionSuccess, nil
	} else {
		v.logger.Error().Str("response_code", responseCode).Msg("Failed transaction")
		return nil, entity.TransciontFailed, errors.New("failed")
	}

}
