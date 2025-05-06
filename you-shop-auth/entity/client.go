package entity

type Client struct {
	AuditEntity
	Enabled      bool
	ClientId     string
	ClientSecret string
	ClientName   string
	Description  string
	WebsiteUrl   string
	RedirectUrl  string
}
