package cache

import "wb-task-L0/internal/models"

type Cache interface {
	Get(key string) (models.Order, bool)
	Add(order models.Order)
	AddAll(orders []models.Order)
}
