package dto

// DeliveryDto represents delivery information for API responses
type DeliveryDto struct {
	Name    string `json:"name" db:"name"`
	Zip     string `json:"zip" db:"zip"`
	City    string `json:"city" db:"city"`
	Address string `json:"address" db:"address"`
	Region  string `json:"region" db:"region"`
}
