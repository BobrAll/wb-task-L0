package dto

// ItemDto represents a product in an order for API responses
type ItemDto struct {
	TrackNumber string  `json:"track_number" db:"track_number"`
	Price       float64 `json:"price" db:"price"`
	Name        string  `json:"name" db:"name"`
	Sale        float32 `json:"sale" db:"sale"`
	Size        string  `json:"size" db:"size"`
	TotalPrice  float64 `json:"total_price" db:"total_price"`
	Brand       string  `json:"brand" db:"brand"`
}
