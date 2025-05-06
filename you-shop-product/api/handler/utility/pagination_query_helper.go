package utility

import (
	"strconv"

	"github.com/TechwizsonORG/product-service/api/constant"
	"github.com/gin-gonic/gin"
)

func GetPaginationQuery(c *gin.Context) (page int, pageSize int) {
	pageStr := c.Query(constant.PAGE_QUERY)
	pageSizeStr := c.Query(constant.PAGE_SIZE_QUERY)

	if pageStr == "" {
		page = 1
	} else {
		page, _ = strconv.Atoi(pageStr)
	}

	if pageSizeStr == "" {
		pageSize = 10
	} else {
		pageSize, _ = strconv.Atoi(pageSizeStr)
	}

	return page, pageSize
}
