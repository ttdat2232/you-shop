package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	configModel "github.com/TechwizsonORG/shipment-service/config/model"
	"github.com/TechwizsonORG/shipment-service/infrastructure/ghn/model"
	"github.com/rs/zerolog"
)

type GhnService struct {
	ghnConfig configModel.GhnConfig
	logger    zerolog.Logger
}

func NewGhnService(ghnConfig configModel.GhnConfig, logger zerolog.Logger) *GhnService {
	logger = logger.With().Str("infrastructure", "ghn").Logger()
	return &GhnService{
		ghnConfig: ghnConfig,
		logger:    logger,
	}
}

func (g *GhnService) GetProvinces() []model.GhnProvinceAddress {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", g.ghnConfig.BaseUrl, "master-data/province"), nil)
	if err != nil {
		g.logger.Error().Err(err).Msg("")
		return make([]model.GhnProvinceAddress, 0)
	}
	req.Header.Add("Token", g.ghnConfig.Token)
	req.Header.Add("ShopId", g.ghnConfig.ShopId)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		g.logger.Error().Err(err).Msg("")
		return make([]model.GhnProvinceAddress, 0)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		g.logger.Error().Str("failed_code", res.Status).Msg("fetch province from Giao Hang Nhanh was failed")
		return make([]model.GhnProvinceAddress, 0)
	}

	var response model.GhnReponse[[]model.GhnProvinceAddress]
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		g.logger.Error().Err(err).Msg("")
		return make([]model.GhnProvinceAddress, 0)
	}

	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		g.logger.Error().Err(err).Msg("")
		return make([]model.GhnProvinceAddress, 0)
	}
	return response.Data
}

func (g *GhnService) GetProviceByName(name string) (*model.GhnProvinceAddress, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", g.ghnConfig.BaseUrl, "master-data/province"), nil)
	if err != nil {
		g.logger.Error().Err(err).Msg("")
		return nil, err
	}
	req.Header.Add("Token", g.ghnConfig.Token)
	req.Header.Add("ShopId", g.ghnConfig.ShopId)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var response model.GhnReponse[[]model.GhnProvinceAddress]
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	for _, province := range response.Data {
		if strings.EqualFold(province.ProvinceName, name) {
			return &province, nil
		}
		for _, provinceName := range province.NameExtension {
			if strings.EqualFold(provinceName, name) {
				return &province, nil
			}
		}
	}
	return nil, fmt.Errorf("province with name %s not found", name)
}

func (g *GhnService) GetDistrictsByProvinceId(provinceId string) ([]model.GhnDistrictAddress, error) {
	provinceIdNumber, parseError := strconv.Atoi(provinceId)
	if parseError != nil {
		return nil, parseError
	}
	requestBody := map[string]any{}
	requestBody["province_id"] = provinceIdNumber
	jsonReqBody, parseError := json.Marshal(requestBody)
	if parseError != nil {
		return nil, parseError
	}
	bodyReader := bytes.NewReader(jsonReqBody)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", g.ghnConfig.BaseUrl, "master-data/district"), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch district: %w", err)
	}
	req.Header.Add("Token", g.ghnConfig.Token)
	req.Header.Add("ShopId", g.ghnConfig.ShopId)
	req.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch district: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var body model.GhnReponse[[]model.GhnDistrictAddress]
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return body.Data, nil
}

func (g *GhnService) GetDistrictByName(name, provinceId string) (*model.GhnDistrictAddress, error) {
	provinceIdNumber, parseError := strconv.Atoi(provinceId)
	if parseError != nil {
		return nil, parseError
	}
	requestBody := map[string]interface{}{}
	requestBody["provice_id"] = provinceIdNumber
	jsonReqBody, parseError := json.Marshal(requestBody)
	if parseError != nil {
		return nil, parseError
	}
	bodyReader := bytes.NewReader(jsonReqBody)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", g.ghnConfig.BaseUrl, "master-data/district"), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch district: %w", err)
	}
	defer req.Body.Close()

	req.Header.Add("Token", g.ghnConfig.Token)
	req.Header.Add("ShopId", g.ghnConfig.ShopId)
	req.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch district: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var body model.GhnReponse[[]model.GhnDistrictAddress]
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	for _, district := range body.Data {
		if strings.EqualFold(district.DistrictName, name) {
			return &district, nil
		}
		for _, districtName := range district.NameExtension {
			if strings.EqualFold(districtName, name) {
				return &district, nil
			}
		}
	}
	return nil, fmt.Errorf("province with name %s not found", name)
}

func (g *GhnService) GetWardByName(name, provinceId string) (*model.GhnWardAddress, error) {
	provinceIdNumber, parseError := strconv.Atoi(provinceId)
	if parseError != nil {
		return nil, parseError
	}
	requestBody := map[string]interface{}{}
	requestBody["district_id"] = provinceIdNumber
	jsonReqBody, parseError := json.Marshal(requestBody)
	if parseError != nil {
		return nil, parseError
	}
	bodyReader := bytes.NewReader(jsonReqBody)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", g.ghnConfig.BaseUrl, "master-data/ward?district_id"), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch district: %w", err)
	}
	defer req.Body.Close()

	req.Header.Add("Token", g.ghnConfig.Token)
	req.Header.Add("ShopId", g.ghnConfig.ShopId)
	req.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch district: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var body model.GhnReponse[[]model.GhnWardAddress]
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	for _, ward := range body.Data {
		if strings.EqualFold(ward.WardName, name) {
			return &ward, nil
		}
		for _, wardName := range ward.NameExtension {
			if strings.EqualFold(wardName, name) {
				return &ward, nil
			}
		}
	}
	return nil, fmt.Errorf("province with name %s not found", name)
}

func (g *GhnService) GetWardsByDisctrictId(districtId string) ([]model.GhnWardAddress, error) {
	provinceIdNumber, parseError := strconv.Atoi(districtId)
	if parseError != nil {
		return nil, parseError
	}
	requestBody := map[string]any{}
	requestBody["district_id"] = provinceIdNumber
	jsonReqBody, parseError := json.Marshal(requestBody)
	if parseError != nil {
		return nil, parseError
	}
	bodyReader := bytes.NewReader(jsonReqBody)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", g.ghnConfig.BaseUrl, "master-data/ward?district_id"), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ward: %w", err)
	}
	req.Header.Add("Token", g.ghnConfig.Token)
	req.Header.Add("ShopId", g.ghnConfig.ShopId)
	req.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch district: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	var body model.GhnReponse[[]model.GhnWardAddress]
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return body.Data, nil
}
