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

func NewAppError(code int, title string, message string, data interface{}) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Title:   title,
		Data:    data,
	}
}

func NewTokenValidationError(message string, data interface{}) *AppError {
	return NewAppError(401, "Token Validation Error", message, data)
}

func NewTokenGenerationError(message string, data interface{}) *AppError {
	return NewAppError(401, "Token Generation Error", message, data)
}

func NewCreateUserError(message string, data interface{}) *AppError {
	return NewAppError(401, "Creating User Error", message, data)
}

func NewWrongUserNameOrPasswordError() *AppError {
	return NewAppError(401, "Wrong username or password", "Wrong username or password", nil)
}

func NewUnhandledError() *AppError {
	return NewAppError(500, "Internal Server Error", "Internal Server Error", nil)
}
