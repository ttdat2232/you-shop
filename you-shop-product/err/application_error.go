package err

type ApplicationError interface {
	error
	Code() int
	Title() string
	Data() interface{}
	Detail() string
}
