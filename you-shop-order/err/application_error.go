package err

type ApplicationError interface {
	error
	ErrCode() int
	ErrTitle() string
	ErrDetail() string
	ErrData() interface{}
}
