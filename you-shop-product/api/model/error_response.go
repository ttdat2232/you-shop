package model

import (
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Type     string      `json:"type"`
	Title    string      `json:"title"`
	Detail   string      `json:"detail"`
	Instance string      `json:"instance,omitempty"`
	Status   int         `json:"status"`
	Data     interface{} `json:"data,omitempty"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", e.Title, e.Detail)
}

type FailField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewValidationErrorResponse(title string, detail string, fields []FailField) ErrorResponse {
	return ErrorResponse{
		Status: http.StatusBadRequest,
		Title:  title,
		Detail: detail,
		Data:   fields,
	}
}
