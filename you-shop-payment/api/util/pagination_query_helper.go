package util

import (
	"strconv"

	"github.com/TechwizsonORG/payment-service/api/constant"
	"github.com/gin-gonic/gin"
)

// Get page and page size with default values corresponding 1, 10
func GetPaginationQuery(c *gin.Context) (page int, pageSize int) {
	pageStr := c.Query(constant.PAGE_QUERY)
	pageSizeStr := c.Query(constant.PAGE_SIZE_QUERY)

	if pageStr == "" {
		page = constant.DEFAULT_PAGE
	} else {
		page, _ = strconv.Atoi(pageStr)
	}
	if pageSizeStr == "" {
		pageSize = constant.DEFAULT_PAGE_SIZE
	} else {
		pageSize, _ = strconv.Atoi(pageSizeStr)
	}

	if page < 1 {
		page = constant.DEFAULT_PAGE
	}

	if pageSize < 1 {
		pageSize = constant.DEFAULT_PAGE_SIZE
	}

	return page, pageSize
}
