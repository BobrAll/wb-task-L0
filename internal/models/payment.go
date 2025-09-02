package models

import "wb-task-L0/internal/transport/dto"

// Payment represents payment details for an order
type Payment struct {
	Transaction  string  `json:"transaction" db:"transaction"`
	RequestID    string  `json:"request_id" db:"request_id"`
	Currency     string  `json:"currency" db:"currency"`
	Provider     string  `json:"provider" db:"provider"`
	Amount       int32   `json:"amount" db:"amount"`
	PaymentDt    int64   `json:"payment_dt" db:"payment_dt"`
	Bank         string  `json:"bank" db:"bank"`
	DeliveryCost float64 `json:"delivery_cost" db:"delivery_cost"`
	GoodsTotal   int32   `json:"goods_total" db:"goods_total"`
	CustomFee    float64 `json:"custom_fee" db:"custom_fee"`
}

// ToDto converts Payment model to DTO
func (payment Payment) ToDto() dto.PaymentDto {
	return dto.PaymentDto{
		Currency:     payment.Currency,
		Provider:     payment.Provider,
		Amount:       payment.Amount,
		PaymentDt:    payment.PaymentDt,
		Bank:         payment.Bank,
		DeliveryCost: payment.DeliveryCost,
		GoodsTotal:   payment.GoodsTotal,
		CustomFee:    payment.CustomFee,
	}
}
