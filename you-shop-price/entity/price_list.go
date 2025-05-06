package entity

type Currency int8

const (
	VND Currency = iota + 1
	USD
)

type PriceList struct {
	AuditEntity
	Description string
	Currency    Currency
}
