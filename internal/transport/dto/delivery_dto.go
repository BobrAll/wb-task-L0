package dto

type DeliveryDto struct {
	Name    string `json:"name" db:"name"`
	Zip     int32  `json:"zip" db:"zip"`
	City    string `json:"city" db:"city"`
	Address string `json:"address" db:"address"`
	Region  string `json:"region" db:"region"`
}
