package dto

// PaymentDto represents payment details for API responses
type PaymentDto struct {
	Currency     string  `json:"currency" db:"currency"`
	Provider     string  `json:"provider" db:"provider"`
	Amount       int32   `json:"amount" db:"amount"`
	PaymentDt    int64   `json:"payment_dt" db:"payment_dt"`
	Bank         string  `json:"bank" db:"bank"`
	DeliveryCost float64 `json:"delivery_cost" db:"delivery_cost"`
	GoodsTotal   int32   `json:"goods_total" db:"goods_total"`
	CustomFee    float64 `json:"custom_fee" db:"custom_fee"`
}
