package cache

import (
	"sync"
	"wb-task-L0/internal/models"
)

type memoryCache struct {
	mu       sync.RWMutex
	size     int
	orderMap map[string]models.Order
	keys     []string
	index    int
}

func New(size int) Cache {
	return &memoryCache{
		size:     size,
		keys:     make([]string, 0, size),
		orderMap: make(map[string]models.Order),
	}
}

func (c *memoryCache) Get(id string) (models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, ok := c.orderMap[id]
	return order, ok
}

func (c *memoryCache) Add(order models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.orderMap[order.OrderID]; ok {
		c.orderMap[order.OrderID] = order
		return
	}

	if len(c.keys) < c.size {
		c.keys = append(c.keys, order.OrderID)
		c.orderMap[order.OrderID] = order
		return
	}

	oldKey := c.keys[c.index]
	delete(c.orderMap, oldKey)

	c.keys[c.index] = order.OrderID
	c.orderMap[order.OrderID] = order

	c.index = (c.index + 1) % c.size
}

func (c *memoryCache) AddAll(orders []models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, order := range orders {
		if _, ok := c.orderMap[order.OrderID]; ok {
			c.orderMap[order.OrderID] = order
			continue
		}

		if len(c.keys) < c.size {
			c.keys = append(c.keys, order.OrderID)
			c.orderMap[order.OrderID] = order
			continue
		}

		oldKey := c.keys[c.index]
		delete(c.orderMap, oldKey)

		c.keys[c.index] = order.OrderID
		c.orderMap[order.OrderID] = order

		c.index = (c.index + 1) % c.size
	}
}
