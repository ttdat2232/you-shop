package token

type TokenType int8

const (
	AccessToken TokenType = iota
	RefreshToken
)

type ValidateTokenRequest struct {
	Token string    `json:"token"`
	Type  TokenType `json:"type"`
}
