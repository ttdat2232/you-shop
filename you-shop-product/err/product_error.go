package err

import (
	"fmt"
	"net/http"
)

type ProductError struct {
	code   int
	title  string
	detail string
	data   interface{}
}

func (p ProductError) Error() string {
	return p.detail
}

func (p ProductError) Code() int {
	return p.code
}

func (p ProductError) Title() string {
	return p.title
}

func (p ProductError) Detail() string {
	return p.detail
}

func (p ProductError) Data() interface{} {
	return p.data
}

func CommonError() ProductError {
	return ProductError{
		code:  http.StatusInternalServerError,
		title: "Internal server error",
	}
}

func NotFoundProductError(detail string) ProductError {
	return ProductError{
		code:   http.StatusNotFound,
		title:  "Not Found",
		detail: detail,
	}
}

func NotFoundProductErrorWithId(id string) ProductError {
	return ProductError{
		code:   http.StatusNotFound,
		title:  "Not Found",
		detail: fmt.Sprintf("Product with id %s not found", id),
	}
}

func NotFoundProductErrorWithData(data interface{}) ProductError {
	return ProductError{
		code:   http.StatusNotFound,
		title:  "Not Found",
		detail: fmt.Sprintf("Product not found %v", data),
	}
}

func NewProductError(code int, title string, detail string, data interface{}) ProductError {
	return ProductError{
		code:   code,
		title:  title,
		detail: detail,
		data:   data,
	}
}
