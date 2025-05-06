package entity

type VendorStatus int8

const (
	Actice VendorStatus = iota + 1
	Inactice
)

type Vendor struct {
	AuditEntity
	Name        string
	Status      VendorStatus
	Key         string
	BaseEnpoint string
	CallBackUrl string
}
