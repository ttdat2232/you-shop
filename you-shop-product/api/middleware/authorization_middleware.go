package middleware

import (
	"strings"

	"github.com/TechwizsonORG/product-service/api/model"
	"github.com/TechwizsonORG/product-service/util"
	"github.com/gin-gonic/gin"
)

func AuthorizationMiddleware(requireRole, requireScope []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if userId := c.GetHeader("userId"); strings.Compare("", userId) == 0 {
			c.JSON(401, model.NewApiResponse(401, "Unauthorized", false, nil))
			c.Abort()
			return
		}

		isRoleMatched := true
		isScopeMatched := true

		if len(requireRole) > 0 {
			isRoleMatched = false
			userRoleArray := strings.Split(c.Request.Header.Get("role"), ",")
			if len(userRoleArray) > 0 {
				isRoleMatched = util.ContainsAny(requireRole, userRoleArray)
			}
		}

		if len(requireScope) > 0 {
			isScopeMatched = false
			userScopeArray := strings.Split(c.Request.Header.Get("scope"), ",")
			if len(userScopeArray) > 0 {
				isScopeMatched = util.ContainsAny(requireScope, userScopeArray)
			}
		}

		if !isRoleMatched || !isScopeMatched {
			c.JSON(401, model.NewApiResponse(401, "Not allowed", false, nil))
			c.Abort()
			return
		}
		c.Next()
	}
}
