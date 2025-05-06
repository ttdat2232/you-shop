package handler

import (
	"github.com/TechwizsonORG/auth-service/api/model"
	requestModel "github.com/TechwizsonORG/auth-service/api/model/token"
	userResponse "github.com/TechwizsonORG/auth-service/api/model/user"
	"github.com/TechwizsonORG/auth-service/usecase/token"
	"github.com/gin-gonic/gin"
)

type TokenHandler struct {
	service token.TokenInterface
}

func NewTokenHandler(service token.TokenInterface) *TokenHandler {
	return &TokenHandler{
		service: service,
	}
}

func (t *TokenHandler) TokenRoutes(route *gin.RouterGroup) {
	tokenRoute := route.Group("/token")
	tokenRoute.POST("/validate", t.validateToken)
}

// ValidateToken godoc
//
//	@Summary	Validate token
//	@Accept		json
//	@Tags		token
//	@Produce	json
//	@Param		login	body		requestModel.ValidateTokenRequest	true	"login model"
//	@Failure	401		{object}	model.ApiResponse
//	@Failure	500		{object}	model.ApiResponse
//	@Success	200		{object}	model.ApiResponse{data=userResponse.UserResponse}
//	@Router		/auth/token/validate [post]
func (t *TokenHandler) validateToken(c *gin.Context) {
	var validateTokenRequest requestModel.ValidateTokenRequest
	if err := c.BindJSON(&validateTokenRequest); err != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: err})
		return
	}

	user, role, scope, validationError := t.service.ValidateTokenWithResponse(validateTokenRequest.Token, validateTokenRequest.Type)
	if validationError != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: validationError})
		return
	}

	c.JSON(200, model.SuccessResponse(userResponse.From(*user, role, scope)))
}
