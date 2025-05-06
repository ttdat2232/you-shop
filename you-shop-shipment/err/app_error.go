package err

type AppError interface {
	error
}

type ShipmentError struct {
	Code    int
	Message string
	Title   string
	Data    any
}

func (s *ShipmentError) Error() string {
	return s.Message
}

func NewAppError(code int, message, title string, data interface{}) *ShipmentError {
	return &ShipmentError{
		Code:    code,
		Message: message,
		Title:   title,
		Data:    data,
	}
}
