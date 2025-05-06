package token

import (
	"strings"
	"time"

	"github.com/TechwizsonORG/auth-service/config/model"
	"github.com/TechwizsonORG/auth-service/entity"
	"github.com/TechwizsonORG/auth-service/err"
	"github.com/TechwizsonORG/auth-service/usecase"
	"github.com/TechwizsonORG/auth-service/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	jwtConfig model.JwtConfig
	logger    zerolog.Logger
	userRepo  usecase.UserRepository
	roleRepo  usecase.RoleRepository
	scopeRepo usecase.ScopeRepository
}

func NewTokenService(logger zerolog.Logger, userRepo usecase.UserRepository, jwtConfig model.JwtConfig, roleRepo usecase.RoleRepository, scopeRepo usecase.ScopeRepository) *Service {
	logger = logger.
		With().
		Str("service", "token").
		Logger()
	return &Service{
		logger:    logger,
		userRepo:  userRepo,
		jwtConfig: jwtConfig,
		roleRepo:  roleRepo,
		scopeRepo: scopeRepo,
	}
}

func (s *Service) getRoles(userId uuid.UUID, roles chan string) {
	roleEntities := s.roleRepo.GetRolesByUserId(userId)
	roleStr := strings.Builder{}
	for i, role := range roleEntities {
		if i > 0 {
			roleStr.WriteRune(',')
		}
		roleStr.WriteString(role.Name)
	}
	roles <- roleStr.String()
}
func (s *Service) getScopes(userId uuid.UUID, scopes chan string) {

	scopeEntities := s.scopeRepo.GetScopesByUserId(userId)
	scopeStr := strings.Builder{}
	for i, scope := range scopeEntities {
		if i > 0 {
			scopeStr.WriteRune(',')
		}
		scopeStr.WriteString(scope.Name)
	}
	scopes <- scopeStr.String()
}

func (s *Service) GenerateToken(userId uuid.UUID, tokenType TokenType) (string, *err.AppError) {
	user, error := s.userRepo.GetUserByID(userId)
	if error != nil {
		s.logger.Err(error).Msgf("Failed to get user by id %s", userId)
		return "", err.NewTokenGenerationError("Failed to get user by id", nil)
	}
	roles := make(chan string)
	scopes := make(chan string)
	go s.getRoles(user.Id, roles)
	go s.getScopes(user.Id, scopes)

	var expired time.Duration
	if tokenType == AccessToken {
		expired = time.Minute * time.Duration(s.jwtConfig.DefaultAccessExpireTime)
	} else {
		expired = time.Hour * time.Duration(s.jwtConfig.DefaultRefreshExpireTime)
	}

	current := util.GetCurrentUtcTime(7)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":        user.Id.String(),
		"email":      user.Email,
		"avatar_url": user.AvatarUrl,
		"iss":        s.jwtConfig.Issuer,
		"iat":        current.Unix(),
		"nbf":        current.Unix(),
		"exp":        current.Add(expired).Unix(),
		"userId":     user.Id.String(),
		"role":       <-roles,
		"scope":      <-scopes,
	})

	token, error := claims.SignedString([]byte(s.jwtConfig.SecretKey))
	if error != nil {
		s.logger.Err(error).Msg("Failed to sign token")
		return "", err.NewTokenGenerationError("Failed to sign token", nil)
	}
	return token, nil
}

func (s *Service) ValidateTokenWithResponse(token string, tokenType TokenType) (user *entity.User, role, scope string, appErr *err.AppError) {
	if strings.Contains(token, "Bearer") {
		token = strings.Split(token, " ")[1]
	}
	jwtToken, parseErr := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.jwtConfig.SecretKey), nil
	})

	if parseErr != nil {
		s.logger.Error().Err(parseErr).Msg("Failed to parse token")
		return nil, "", "", err.NewTokenValidationError("Failed to parse token", nil)
	}

	if jwtToken.Method != jwt.SigningMethodHS256 {
		return nil, "", "", err.NewTokenValidationError("Invalid signing method", nil)
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, "", "", err.NewTokenValidationError("Failed when getting claims", nil)
	}

	current := util.GetCurrentUtcTime(7)

	expire, expiredError := claims.GetExpirationTime()
	if expiredError != nil {
		return nil, "", "", err.NewTokenValidationError("Failed when getting expiration time", nil)
	}
	if current.After(expire.Time) {
		return nil, "", "", err.NewTokenValidationError("Token was expired", nil)
	}
	notBefore, notBeforeError := claims.GetNotBefore()

	if notBeforeError != nil {
		return nil, "", "", err.NewTokenValidationError("Failed when getting not before time", nil)
	}

	if current.Before(notBefore.Time) {
		return nil, "", "", err.NewTokenValidationError("Not before time was failed validation", nil)
	}

	issuer, issError := claims.GetIssuer()
	if issError != nil {
		return nil, "", "", err.NewTokenValidationError("Failed when getting not before time", nil)
	}

	if issuer != s.jwtConfig.Issuer {
		return nil, "", "", err.NewTokenValidationError("Issuer was failed validation", nil)
	}

	userIdStr, ok := claims["userId"]
	if !ok {
		return nil, "", "", err.NewAppError(401, "Couldn't get user id", "Couldn't get user id", nil)
	}

	userId, parseUuidErr := uuid.Parse(userIdStr.(string))
	if parseUuidErr != nil {
		s.logger.Error().Err(parseUuidErr).Msg("")
		return nil, "", "", err.NewAppError(401, "Couldn't parse user id", "Couldn't parse user id", nil)
	}

	user, getUserErr := s.userRepo.GetUserByID(userId)
	if getUserErr != nil {
		s.logger.Error().Err(getUserErr).Msg("")
		return nil, "", "", err.NewAppError(401, "Couldn't find user", "Couldn't find user", nil)
	}
	roles, ok := claims["role"].(string)
	if !ok {
		return nil, "", "", err.NewAppError(401, "Couldn't get user role", "Couldn't get user role", nil)
	}
	scopes, ok := claims["scope"].(string)
	if !ok {
		return nil, "", "", err.NewAppError(401, "Couldn't get user scope", "Couldn't get user scope", nil)
	}
	return user, roles, scopes, nil
}
