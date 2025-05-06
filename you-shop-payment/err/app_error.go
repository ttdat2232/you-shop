package err

type AppError interface {
	error
}

type PaymentError struct {
	Code    int
	Message string
	Title   string
	Data    any
}

func (s *PaymentError) Error() string {
	return s.Message
}

func NewAppError(code int, message, title string, data interface{}) *PaymentError {
	return &PaymentError{
		Code:    code,
		Message: message,
		Title:   title,
		Data:    data,
	}
}

func NewCommonErr() *PaymentError {
	return &PaymentError{
		Code:    500,
		Message: "Error occurred",
		Title:   "Error occurred",
		Data:    nil,
	}
}
