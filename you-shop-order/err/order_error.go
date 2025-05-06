package err

type OrderError struct {
	Code   int         `json:"code"`
	Title  string      `json:"title"`
	Data   interface{} `json:"data"`
	Detail string      `json:"detail"`
}

func (o *OrderError) Error() string {
	return o.Title
}

func NewOrderError(code int, title, detail string, data interface{}) *OrderError {
	return &OrderError{
		Code:   code,
		Title:  title,
		Data:   data,
		Detail: detail,
	}
}

// Create OrderError object with Code - 500, Title and Detail - "Error Occurred"
func NewOrderDefaultError(data interface{}) *OrderError {
	return &OrderError{
		Code:   500,
		Title:  "Error Occurred",
		Detail: "Error Occurred",
		Data:   data,
	}
}

func (o *OrderError) ErrCode() int {
	return o.Code
}
func (o *OrderError) ErrTitle() string {
	return o.Title
}
func (o *OrderError) ErrDetail() string {
	return o.Detail
}
func (o *OrderError) ErrData() interface{} {
	return o.Data
}
