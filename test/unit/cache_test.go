package unit

import (
	"testing"
	"wb-task-L0/internal/cache"
	"wb-task-L0/internal/models"
)

func makeOrder(id string) models.Order {
	return models.Order{OrderID: id}
}

func TestAddAndGet(t *testing.T) {
	c := cache.New(2)

	order1 := makeOrder("1")
	c.Add(order1)

	if got, ok := c.Get("1"); !ok || got.OrderID != "1" {
		t.Errorf("expected order 1, got %+v, ok=%v", got, ok)
	}

	if _, ok := c.Get("2"); ok {
		t.Errorf("expected no order with id 2")
	}
}

func TestAddOverwriteExisting(t *testing.T) {
	c := cache.New(2)

	order1 := makeOrder("1")
	c.Add(order1)

	updated := makeOrder("1")
	updated.TrackNumber = "TRACK-UPDATED"
	c.Add(updated)

	got, ok := c.Get("1")
	if !ok {
		t.Fatalf("expected order with id 1")
	}
	if got.TrackNumber != "TRACK-UPDATED" {
		t.Errorf("expected updated TrackNumber, got %s", got.TrackNumber)
	}
}

func TestEviction(t *testing.T) {
	c := cache.New(2)

	c.Add(makeOrder("1"))
	c.Add(makeOrder("2"))
	c.Add(makeOrder("3"))

	if _, ok := c.Get("1"); ok {
		t.Errorf("expected order 1 to be evicted")
	}
	if _, ok := c.Get("2"); !ok {
		t.Errorf("expected order 2 to remain")
	}
	if _, ok := c.Get("3"); !ok {
		t.Errorf("expected order 3 to remain")
	}
}

func TestAddAll(t *testing.T) {
	c := cache.New(3)

	orders := []models.Order{
		makeOrder("1"),
		makeOrder("2"),
		makeOrder("3"),
	}
	c.AddAll(orders)

	for _, o := range orders {
		if _, ok := c.Get(o.OrderID); !ok {
			t.Errorf("expected order %s to be in cache", o.OrderID)
		}
	}
}

func TestAddAllEviction(t *testing.T) {
	c := cache.New(2)

	orders := []models.Order{
		makeOrder("1"),
		makeOrder("2"),
		makeOrder("3"),
	}
	c.AddAll(orders)

	if _, ok := c.Get("1"); ok {
		t.Errorf("expected order 1 to be evicted")
	}
	if _, ok := c.Get("2"); !ok {
		t.Errorf("expected order 2 to remain")
	}
	if _, ok := c.Get("3"); !ok {
		t.Errorf("expected order 3 to remain")
	}
}
