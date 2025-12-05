package models

type Vendor struct {
	VendorId string `gorm:"primaryKey" json:"vendorId"`
	Name string `json:"name"`
	Contact string `json:"contact"`
	Phone string `json:"phone"`
	Email string `json:"email"`
	Address string `json:"address"`
}