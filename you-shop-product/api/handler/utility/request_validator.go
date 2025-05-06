package utility

import (
	"strconv"

	"github.com/TechwizsonORG/product-service/api/constant"
	"github.com/TechwizsonORG/product-service/err"
	"github.com/gin-gonic/gin"
)

func PaginationValidator(c *gin.Context) (bool, err.ValidationError) {
	page := c.Query(constant.PAGE_QUERY)
	pageSize := c.Query(constant.PAGE_SIZE_QUERY)

	if page == "" && pageSize == "" {
		return true, err.ValidationError{}
	}

	pageInt, atoiErr := strconv.Atoi(page)
	if atoiErr != nil || pageInt < 1 {
		return false, err.NewValidationError("Invalid page", "Page must be a positive number", []err.ValidationErrorField{
			{
				Field:   constant.PAGE_QUERY,
				Message: "Page must be a positive number",
			},
		})
	}

	pageSizeInt, atoiErr := strconv.Atoi(pageSize)
	if atoiErr != nil || pageSizeInt < 1 {
		return false, err.NewValidationError("Invalid page size", "Page size must be a positive number", []err.ValidationErrorField{
			{
				Field:   constant.PAGE_SIZE_QUERY,
				Message: "Page size must be a positive number",
			},
		})
	}

	return true, err.ValidationError{}
}
