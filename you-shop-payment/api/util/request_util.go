package util

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserId(c *gin.Context) (uuid.UUID, error) {
	userIdStr := c.Request.Header.Get("userId")
	if strings.Compare("", userIdStr) == 0 { // empty string
		return uuid.UUID{}, errors.New("couldn't get user id from header")
	}
	if userId, err := uuid.Parse(userIdStr); err != nil {

		return uuid.UUID{}, errors.New("couldn't get user id from header")
	} else {
		return userId, nil
	}
}
