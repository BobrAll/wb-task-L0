package models

import (
	"wb-task-L0/internal/transport/dto"
)

// Delivery represents delivery information for an order
type Delivery struct {
	Name    string `json:"name" db:"name"`
	Phone   string `json:"phone" db:"phone"`
	Zip     string `json:"zip" db:"zip"`
	City    string `json:"city" db:"city"`
	Address string `json:"address" db:"address"`
	Region  string `json:"region" db:"region"`
	Email   string `json:"email" db:"email"`
}

// ToDto converts Delivery model to DTO
func (delivery Delivery) ToDto() dto.DeliveryDto {
	return dto.DeliveryDto{
		Name:    delivery.Name,
		Zip:     delivery.Zip,
		City:    delivery.City,
		Address: delivery.Address,
		Region:  delivery.Region,
	}
}
