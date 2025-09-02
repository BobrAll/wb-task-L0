package models

import (
	"time"
	"wb-task-L0/internal/transport/dto"
)

// Order represents a customer's order with delivery, payment, and items
type Order struct {
	OrderID           string    `json:"order_uid" db:"order_uid"`
	TrackNumber       string    `json:"track_number" db:"track_number"`
	Entry             string    `json:"entry" db:"entry"`
	Delivery          Delivery  `json:"delivery" db:"delivery"`
	Payment           Payment   `json:"payment" db:"payment"`
	Items             []Item    `json:"items" db:"-"`
	Locale            string    `json:"locale" db:"locale"`
	InternalSignature string    `json:"internal_signature" db:"internal_signature"`
	CustomerID        string    `json:"customer_id" db:"customer_id"`
	DeliveryService   string    `json:"delivery_service" db:"delivery_service"`
	ShardKey          int32     `json:"shard_key" db:"shard_key"`
	SmID              int32     `json:"sm_id" db:"sm_id"`
	DateCreated       time.Time `json:"date_created" db:"date_created"`
	OofShard          string    `json:"oof_shard" db:"oof_shard"`
}

// ToDto converts Order model to DTO
func (order *Order) ToDto() dto.OrderDto {
	itemsDto := make([]dto.ItemDto, len(order.Items))
	for i, item := range order.Items {
		itemsDto[i] = item.ToDto()
	}

	return dto.OrderDto{
		OrderID:         order.OrderID,
		TrackNumber:     order.TrackNumber,
		Entry:           order.Entry,
		Delivery:        order.Delivery.ToDto(),
		Payment:         order.Payment.ToDto(),
		Items:           itemsDto,
		DeliveryService: order.DeliveryService,
		DateCreated:     order.DateCreated,
	}
}
