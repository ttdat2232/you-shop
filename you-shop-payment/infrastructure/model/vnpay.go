package model

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/TechwizsonORG/payment-service/config/model"
	"github.com/TechwizsonORG/payment-service/util"
	"github.com/google/uuid"
)

type VNPayReturn struct {
	Version       string    `json:"vnp_Version"`       // Version of the API
	Command       string    `json:"vnp_Command"`       // Command (e.g., "pay")
	TxnRef        string    `json:"vnp_TxnRef"`        // Merchant's transaction reference
	OrderInfo     uuid.UUID `json:"vnp_OrderInfo"`     // Description of the order
	Amount        int64     `json:"vnp_Amount"`        // Amount in the smallest currency unit (e.g., cents)
	ResponseCode  string    `json:"vnp_ResponseCode"`  // Response code (e.g., "00" for success)
	TransactionNo string    `json:"vnp_TransactionNo"` // VNPAY's transaction number
	BankCode      string    `json:"vnp_BankCode"`      // Bank code (if applicable)
	PayDate       time.Time `json:"vnp_PayDate"`       // Payment date and time
	SecureHash    string    `json:"vnp_SecureHash"`    // Secure hash for verification
}

type VNPayPaymentRequest struct {
	Version    string    `json:"vnp_Version"`    // API version (e.g., "2.1.0")
	Command    string    `json:"vnp_Command"`    // Command (e.g., "pay")
	TmnCode    string    `json:"vnp_TmnCode"`    // Merchant code
	Amount     int64     `json:"vnp_Amount"`     // Amount in the smallest currency unit (e.g., cents)
	Currency   string    `json:"vnp_CurrCode"`   // Currency code (e.g., "VND")
	TxnRef     string    `json:"vnp_TxnRef"`     // Merchant's transaction reference
	OrderInfo  uuid.UUID `json:"vnp_OrderInfo"`  // Description of the order - orderId
	Locale     string    `json:"vnp_Locale"`     // Language (e.g., "vn" for Vietnamese)
	ReturnURL  string    `json:"vnp_ReturnUrl"`  // URL to return after payment
	IpAddr     string    `json:"vnp_IpAddr"`     // Customer's IP address
	CreateDate string    `json:"vnp_CreateDate"` // Payment creation date (format: yyyyMMddHHmmss)
	BankCode   string    `json:"vnp_BankCode"`   // Bank code (optional)
	SecureHash string    `json:"vnp_SecureHash"` // Secure hash for verification
}

// CreatePaymentURL generates a payment URL for VNPAY
func CreatePaymentURL(request VNPayPaymentRequest, vnpayConfig model.VnpayConfig) (string, error) {
	// Set default values if not provided
	if request.Version == "" {
		request.Version = vnpayConfig.Version
	}
	if request.Command == "" {
		request.Command = vnpayConfig.Command
	}
	if request.TmnCode == "" {
		request.TmnCode = vnpayConfig.MerchantCode
	}
	if request.Currency == "" {
		request.Currency = vnpayConfig.Currency
	}
	if request.Locale == "" {
		request.Locale = vnpayConfig.Locale
	}
	if request.ReturnURL == "" {
		request.ReturnURL = vnpayConfig.ReturnUrl
	}
	if request.CreateDate == "" {
		request.CreateDate = util.GetCurrentUtcTime(7).Format("20060102150405") // yyyyMMddHHmmss
	}

	// Convert amount to string (in cents)
	amountStr := fmt.Sprintf("%d", request.Amount)

	// Prepare parameters
	params := map[string]string{
		"vnp_Version":    request.Version,
		"vnp_Command":    request.Command,
		"vnp_TmnCode":    request.TmnCode,
		"vnp_Amount":     amountStr,
		"vnp_CurrCode":   request.Currency,
		"vnp_TxnRef":     request.TxnRef,
		"vnp_OrderInfo":  request.OrderInfo.String(),
		"vnp_Locale":     request.Locale,
		"vnp_ReturnUrl":  request.ReturnURL,
		"vnp_IpAddr":     request.IpAddr,
		"vnp_CreateDate": request.CreateDate,
	}

	// Add optional bank code if provided
	if request.BankCode != "" {
		params["vnp_BankCode"] = request.BankCode
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
	hash := hmac.New(sha512.New, []byte(vnpayConfig.SecretKey))
	hash.Write([]byte(queryString))
	secureHash := hex.EncodeToString(hash.Sum(nil))

	// Append secure hash to query string
	paymentURL := fmt.Sprintf("%s?%s&vnp_SecureHash=%s", vnpayConfig.PageUrl, queryString, secureHash)
	return paymentURL, nil
}

func ParseVNPayReturn(queryParams url.Values) (*VNPayReturn, error) {
	// Parse the payment date
	payDate, err := time.Parse("20060102150405", queryParams.Get("vnp_PayDate"))
	if err != nil {
		return nil, fmt.Errorf("invalid payment date: %v", err)
	}

	// Parse the amount (convert from string to int64)
	amountStr := queryParams.Get("vnp_Amount")
	amount := int64(0)
	if amountStr != "" {
		_, err := fmt.Sscanf(amountStr, "%d", &amount)
		if err != nil {
			return nil, fmt.Errorf("invalid amount: %v", err)
		}
	}

	// Create and return the VNPayReturn struct
	return &VNPayReturn{
		Version:       queryParams.Get("vnp_Version"),
		Command:       queryParams.Get("vnp_Command"),
		TxnRef:        queryParams.Get("vnp_TxnRef"),
		OrderInfo:     uuid.MustParse(queryParams.Get("vnp_OrderInfo")),
		Amount:        amount,
		ResponseCode:  queryParams.Get("vnp_ResponseCode"),
		TransactionNo: queryParams.Get("vnp_TransactionNo"),
		BankCode:      queryParams.Get("vnp_BankCode"),
		PayDate:       payDate,
		SecureHash:    queryParams.Get("vnp_SecureHash"),
	}, nil
}
