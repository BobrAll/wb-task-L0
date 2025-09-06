package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"

	"wb-task-L0/internal/cache"
	"wb-task-L0/internal/models"
)

var redisContainer *redis.RedisContainer

// setupRedis starts a test Redis container and returns the cache
func setupRedis(t *testing.T) cache.Cache {
	ctx := context.Background()

	container, err := redis.RunContainer(ctx,
		testcontainers.WithImage("redis:7-alpine"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithOccurrence(1).
				WithStartupTimeout(20*time.Second),
		),
	)
	require.NoError(t, err)

	redisContainer = container

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "6379/tcp")
	require.NoError(t, err)

	os.Setenv("REDIS_ADDR", host+":"+port.Port())
	os.Setenv("REDIS_PASSWORD", "")
	os.Setenv("REDIS_DB", "0")
	os.Setenv("CACHE_SIZE", "100")
	os.Setenv("CACHE_TTL_MINUTES", "1")

	c := cache.NewRedis()
	return c
}

// teardownRedis stops the Redis container
func teardownRedis(t *testing.T) {
	require.NoError(t, redisContainer.Terminate(context.Background()))
}

// makeTestOrder creates a test order with a given ID
func makeTestOrder(id string) models.Order {
	return models.Order{
		OrderID:     id,
		TrackNumber: "track-" + id,
		Entry:       "web",
		Delivery: models.Delivery{
			Name:  "John",
			Phone: "123456",
			City:  "City",
		},
		Payment: models.Payment{
			Transaction: "txn-" + id,
			Currency:    "USD",
			Amount:      42,
		},
		Items: []models.Item{
			{
				ChrtID:      1,
				TrackNumber: "track-" + id,
				Price:       42,
				Name:        "Item-" + id,
				TotalPrice:  42,
			},
		},
		DateCreated: time.Now(),
	}
}

// TestRedisCache_AddAndGet tests adding a single order to Redis and retrieving it
func TestRedisCache_AddAndGet(t *testing.T) {
	c := setupRedis(t)
	defer teardownRedis(t)

	order := makeTestOrder("order123")
	c.Add(order)

	got, ok := c.Get("order123")
	require.True(t, ok)
	require.Equal(t, order.OrderID, got.OrderID)
	require.Equal(t, order.Delivery.Name, got.Delivery.Name)
	require.Len(t, got.Items, 1)
}

// TestRedisCache_AddAllAndGet tests adding multiple orders to Redis and retrieving them
func TestRedisCache_AddAllAndGet(t *testing.T) {
	c := setupRedis(t)
	defer teardownRedis(t)

	orders := []models.Order{
		makeTestOrder("order1"),
		makeTestOrder("order2"),
		makeTestOrder("order3"),
	}
	c.AddAll(orders)

	for _, o := range orders {
		got, ok := c.Get(o.OrderID)
		require.True(t, ok)
		require.Equal(t, o.OrderID, got.OrderID)
	}
}
