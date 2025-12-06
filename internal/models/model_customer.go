package models

type Customer struct {
	CustomerId string `gorm:"primaryKey" json:"customerId"`
	FirstName string `json:"firstName" validate:"required,min=1,max=60"`
	LastName string `json:"lastName"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone" validate:"e164,min=11,max=12"`
	Address string `json:"address"`
}