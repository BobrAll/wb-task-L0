package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"strconv"
	"time"
	"wb-task-L0/internal/models"
)

// redisCache implements Cache using Redis
type redisCache struct {
	client    *redis.Client
	ctx       context.Context
	size      int
	keysList  string
	orderPref string
	ttl       time.Duration
}

var (
	DefaultCacheSize       = 1000
	DefaultCacheTTLMinutes = 5
)

// NewRedis creates a new redisCache from env vars
func NewRedis() Cache {
	cacheSize, err := strconv.Atoi(os.Getenv("CACHE_SIZE"))
	if err != nil {
		log.Printf("WARN: Couldn't parse CACHE_SIZE (%v), using default = %d\n", err, DefaultCacheSize)
		cacheSize = DefaultCacheSize
	}
	ttl, err := strconv.Atoi(os.Getenv("CACHE_TTL_MINUTES"))
	if err != nil {
		log.Printf("WARN: Couldn't parse CACHE_TTL_MINUTES (%v), using default = %d", err, DefaultCacheTTLMinutes)
		ttl = DefaultCacheTTLMinutes
	}

	return &redisCache{
		client:    redis.NewClient(loadRedisOptions()),
		ctx:       context.Background(),
		size:      cacheSize,
		keysList:  "orders:keys",
		orderPref: "orders:data:",
		ttl:       time.Duration(ttl) * time.Minute,
	}
}

// loadRedisOptions loads Redis connection options from env
func loadRedisOptions() *redis.Options {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatal("Couldn't parse REDIS_DB", err)
	}
	return &redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	}
}

// Get retrieves an order by ID
func (c *redisCache) Get(id string) (models.Order, bool) {
	val, err := c.client.Get(c.ctx, c.orderPref+id).Result()
	if errors.Is(err, redis.Nil) {
		return models.Order{}, false
	}
	if err != nil {
		fmt.Println("Redis GET error:", err)
		return models.Order{}, false
	}

	var order models.Order
	if err := json.Unmarshal([]byte(val), &order); err != nil {
		fmt.Println("Unmarshal error", err)
		return models.Order{}, false
	}

	return order, true
}

// Add stores one order
func (c *redisCache) Add(order models.Order) {
	data, err := json.Marshal(order)
	if err != nil {
		fmt.Println("Marshal error", err)
	}

	key := c.orderPref + order.OrderID
	pipe := c.client.TxPipeline()

	if c.ttl > 0 {
		pipe.Set(c.ctx, key, data, c.ttl)
	} else {
		pipe.Set(c.ctx, key, data, 0)
	}

	pipe.LPush(c.ctx, c.keysList, order.OrderID)
	pipe.LTrim(c.ctx, c.keysList, 0, int64(c.size-1))

	if _, err := pipe.Exec(c.ctx); err != nil {
		fmt.Println("Redis Add error", err)
	}
}

// AddAll stores multiple orders
func (c *redisCache) AddAll(orders []models.Order) {
	pipe := c.client.TxPipeline()

	for _, order := range orders {
		data, err := json.Marshal(order)
		if err != nil {
			fmt.Println("Marshal error", err)
			continue
		}

		key := c.orderPref + order.OrderID

		if c.ttl > 0 {
			pipe.Set(c.ctx, key, data, c.ttl)
		} else {
			pipe.Set(c.ctx, key, data, 0)
		}
		pipe.LPush(c.ctx, c.keysList, order.OrderID)
	}

	pipe.LTrim(c.ctx, c.keysList, 0, int64(c.size-1))

	if _, err := pipe.Exec(c.ctx); err != nil {
		fmt.Println("Redis AddAll error", err)
	}
}
