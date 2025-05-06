package err

import "fmt"

type AppError struct {
	Code    int
	Message string
	Title   string
	Data    interface{}
}

func (ae *AppError) Error() string {
	return fmt.Sprintf("%s:\n%s", ae.Title, ae.Message)
}

func NewAppError(code int, message string, title string, data interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Title:   title,
		Data:    data,
	}
}

func NewValidationErr(message string, title string, data interface{}) *AppError {
	return &AppError{
		Code:    400,
		Message: message,
		Title:   title,
		Data:    data,
	}
}
