package main

import (
	"log"
	"os"
	"strconv"
	"wb-task-L0/internal/cache"
	"wb-task-L0/internal/config"
	"wb-task-L0/internal/db"
	"wb-task-L0/internal/server"
	"wb-task-L0/internal/transport/kafka"
)

func main() {
	config.LoadEnv()
	db.RunMigrations()
	dbConn := db.InitConn()
	repo := db.NewOrderRepository(dbConn)
	defer repo.Db.Close()

	go kafka.StartOrderListener()

	ordersCache := cache.NewRedis()
	cacheSize, _ := strconv.Atoi(os.Getenv("CACHE_SIZE"))
	orders, err := repo.GetLatestOrders(cacheSize)
	if err != nil {
		log.Fatal(err)
	}

	ordersCache.AddAll(orders)
	server.StartServer(repo, ordersCache)
}
