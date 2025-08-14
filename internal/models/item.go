package models

type Item struct {
	ChrtID      int64   `json:"chrt_id" db:"chrt_id"`
	OrderID     string  `json:"-" db:"order_id"`
	TrackNumber string  `json:"track_number" db:"track_number"`
	Price       float64 `json:"price" db:"price"`
	RID         string  `json:"rid" db:"rid"`
	Name        string  `json:"name" db:"name"`
	Sale        float32 `json:"sale" db:"sale"`
	Size        int32   `json:"size" db:"size"`
	TotalPrice  float64 `json:"total_price" db:"total_price"`
	NmId        int64   `json:"nm_id" db:"nm_id"`
	Brand       string  `json:"brand" db:"brand"`
	Status      int32   `json:"status" db:"status"`
}
