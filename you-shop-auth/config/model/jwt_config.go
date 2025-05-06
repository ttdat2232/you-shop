package model

type JwtConfig struct {
	SecretKey                string
	DefaultAccessExpireTime  int
	DefaultRefreshExpireTime int
	Issuer                   string
}
