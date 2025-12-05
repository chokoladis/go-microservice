package models

type Service struct {
	ServiceId string `gorm:"primaryKey" json:"serviceId"`
	Name string `json:"name"`
	Price float64 `json:"price"`
}