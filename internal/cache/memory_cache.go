package cache

import (
	"sync"
	"wb-task-L0/internal/models"
)

// memoryCache implements Cache interface using in-memory storage
type memoryCache struct {
	mu       sync.RWMutex
	size     int
	orderMap map[string]models.Order
	keys     []string
	index    int
}

// New creates a new memoryCache with given size of orders
func New(maxOrders int) Cache {
	return &memoryCache{
		size:     maxOrders,
		keys:     make([]string, 0, maxOrders),
		orderMap: make(map[string]models.Order),
	}
}

// Get retrieves an order from cache by ID
func (c *memoryCache) Get(id string) (models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, ok := c.orderMap[id]
	return order, ok
}

// Add stores or updates a single order in cache
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

// AddAll stores or updates multiple orders in cache
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
