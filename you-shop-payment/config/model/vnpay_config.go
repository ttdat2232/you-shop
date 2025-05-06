package model

type VnpayConfig struct {
	OrderType    string
	Command      string
	PageUrl      string
	SecretKey    string
	MerchantCode string
	ReturnUrl    string
	Version      string
	Currency     string
	Locale       string
}
