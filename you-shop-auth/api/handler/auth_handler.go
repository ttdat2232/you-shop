package handler

import (
	"github.com/TechwizsonORG/auth-service/api/model"
	authModel "github.com/TechwizsonORG/auth-service/api/model/auth"
	service "github.com/TechwizsonORG/auth-service/usecase/auth"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthInterface
}

func NewAuthHandler(authService service.AuthInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (a *AuthHandler) AuthRoutes(router *gin.RouterGroup) {
	router.POST("login", a.login)
	router.POST("register", a.Register)
}

// Login godoc
//
//	@Summary	Login
//	@Accept		json
//	@Tags		auth
//	@Produce	json
//	@Param		login	body		authModel.LoginRequest	true	"login model"
//	@Failure	401		{object}	model.ApiResponse
//	@Failure	500		{object}	model.ApiResponse
//	@Success	200		{object}	model.ApiResponse{data=authModel.LoginResponse}
//	@Router		/auth/login [post]
func (a *AuthHandler) login(c *gin.Context) {
	var loginReq authModel.LoginRequest
	if err := c.BindJSON(&loginReq); err != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: err})
		return
	}
	user, accessToken, refreshToken, e := a.authService.Login(loginReq.Email, loginReq.Username, loginReq.Password)
	if e != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: e})
		return
	}
	c.JSON(200, model.SuccessResponse(authModel.From(*user, accessToken, refreshToken)))
}

// Register godoc
//
//	@Summary	Register
//	@Accept		json
//	@Tags		auth
//	@Produce	json
//	@Param		register	body		authModel.RegisterRequest	true	"register model"
//	@Failure	401			{object}	model.ApiResponse
//	@Failure	500			{object}	model.ApiResponse
//	@Success	201
//	@Router		/auth/register [post]
func (a *AuthHandler) Register(c *gin.Context) {
	var registerReq authModel.RegisterRequest
	if err := c.BindJSON(&registerReq); err != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: err})
		return
	}
	_, e := a.authService.Register(registerReq.Email, registerReq.Username, registerReq.Password, registerReq.PhoneNumber)

	if e != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: e})
		return
	}

	c.Status(201)
}
