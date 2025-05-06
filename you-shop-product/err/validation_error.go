package err

import "net/http"

type ValidationError struct {
	title  string
	detail string
	data   []ValidationErrorField
}

func (p ValidationError) Error() string {
	return p.title
}

func (p ValidationError) Code() int {
	return http.StatusBadRequest
}

func (p ValidationError) Title() string {
	return p.title
}

func (p ValidationError) Detail() string {
	return p.detail
}

func (p ValidationError) Data() interface{} {
	return p.data
}

type ValidationErrorField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewValidationError(title string, detail string, data []ValidationErrorField) ValidationError {
	return ValidationError{
		title:  title,
		detail: detail,
		data:   data,
	}
}
