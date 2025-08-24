package models

import "wb-task-L0/internal/transport/dto"

type Item struct {
	ChrtID      int64   `json:"chrt_id" db:"chrt_id"`
	TrackNumber string  `json:"track_number" db:"track_number"`
	Price       float64 `json:"price" db:"price"`
	RID         string  `json:"rid" db:"rid"`
	Name        string  `json:"name" db:"name"`
	Sale        float32 `json:"sale" db:"sale"`
	Size        string  `json:"size" db:"size"`
	TotalPrice  float64 `json:"total_price" db:"total_price"`
	NmId        int64   `json:"nm_id" db:"nm_id"`
	Brand       string  `json:"brand" db:"brand"`
	Status      int32   `json:"status" db:"status"`
}

func (item Item) ToDto() dto.ItemDto {
	return dto.ItemDto{
		TrackNumber: item.TrackNumber,
		Price:       item.Price,
		Name:        item.Name,
		Sale:        item.Sale,
		Size:        item.Size,
		TotalPrice:  item.TotalPrice,
		Brand:       item.Brand,
	}
}
