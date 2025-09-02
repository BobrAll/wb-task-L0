package dto

import "time"

// OrderDto represents an order for API responses
type OrderDto struct {
	OrderID         string      `json:"order_uid"`
	TrackNumber     string      `json:"track_number"`
	Entry           string      `json:"entry"`
	Delivery        DeliveryDto `json:"delivery"`
	Payment         PaymentDto  `json:"payment"`
	Items           []ItemDto   `json:"items"`
	DeliveryService string      `json:"delivery_service"`
	DateCreated     time.Time   `json:"date_created"`
}
