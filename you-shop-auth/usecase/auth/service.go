package auth

import (
	"github.com/TechwizsonORG/auth-service/entity"
	"github.com/TechwizsonORG/auth-service/err"
	"github.com/TechwizsonORG/auth-service/usecase"
	"github.com/TechwizsonORG/auth-service/usecase/token"
	"github.com/TechwizsonORG/auth-service/util"
	"github.com/rs/zerolog"
)

type Service struct {
	logger       zerolog.Logger
	userRepo     usecase.UserRepository
	tokenService token.TokenInterface
}

func NewAuthService(logger zerolog.Logger, userRepo usecase.UserRepository, tokenService token.TokenInterface) *Service {
	logger = logger.
		With().
		Str("Service", "Auth").
		Logger()
	return &Service{
		logger:       logger,
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (s *Service) Login(email string, username string, password string) (u *entity.User, accessToken string, refreshToken string, e *err.AppError) {
	user, error := s.userRepo.GetUserByEmailOrUsername(email, username)
	if error != nil {
		s.logger.Err(error).Msg("")
		return nil, "", "", err.NewWrongUserNameOrPasswordError()
	}
	if !util.VerifyPassword(password, user.PasswordHash) {
		return nil, "", "", err.NewWrongUserNameOrPasswordError()
	}

	accessToken, e = s.tokenService.GenerateToken(user.Id, token.AccessToken)

	if e != nil {
		return nil, "", "", e
	}

	refreshToken, e = s.tokenService.GenerateToken(user.Id, token.RefreshToken)

	if e != nil {
		return nil, "", "", e
	}

	return user, accessToken, refreshToken, nil
}
func (s *Service) Register(email string, username string, password string, phoneNumber string) (*entity.User, *err.AppError) {

	user, error := s.userRepo.GetUserByEmailOrUsername(email, "")
	if error != nil {
		s.logger.Error().Err(error).Msg("")
		return nil, err.NewCreateUserError("Error occurred", nil)
	}

	if user != nil {
		return nil, err.NewCreateUserError("Email was already used", nil)
	}

	user, error = s.userRepo.GetUserByEmailOrUsername("", username)

	if error != nil {
		s.logger.Error().Err(error).Msg("")
		return nil, err.NewCreateUserError("Error occurred", nil)
	}

	if user != nil {
		return nil, err.NewCreateUserError("Username was already used", nil)
	}

	newUser, e := entity.CreateUser(username, password, email, "", phoneNumber)
	if e != nil {
		return nil, e
	}
	createdUser, error := s.userRepo.AddUser(*newUser)

	if error != nil {
		s.logger.Err(error).Msg("")
		return nil, err.NewCreateUserError("Error occurred", nil)
	}

	return &createdUser, nil
}
