package model

import "net/http"

type ApiResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	IsSuccess bool        `json:"is_success"`
	Data      interface{} `json:"data,omitempty"`
}

func NewApiResponse(code int, message string, isSuccess bool, data interface{}) *ApiResponse {
	return &ApiResponse{Code: code, Message: message, IsSuccess: isSuccess, Data: data}
}

func SuccessResponse(data interface{}) *ApiResponse {
	return &ApiResponse{
		Code:      http.StatusOK,
		Message:   "Success",
		IsSuccess: true,
		Data:      data,
	}
}
