package err

import "net/http"

type ValidationError struct {
	Title  string
	Detail string
	Data   []ValidationErrorField
	Code   int
}

func (v *ValidationError) Error() string {
	return v.Title
}

type ValidationErrorField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewValidationError(title string, detail string, data []ValidationErrorField) *ValidationError {
	return &ValidationError{
		Title:  title,
		Detail: detail,
		Data:   data,
		Code:   http.StatusBadRequest,
	}
}
func (v *ValidationError) ErrCode() int {
	return v.Code
}
func (v *ValidationError) ErrTitle() string {
	return v.Title
}
func (v *ValidationError) ErrDetail() string {
	return v.Detail
}
func (v *ValidationError) ErrData() interface{} {
	return v.Data
}
