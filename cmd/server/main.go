package main

import (
	"log"
	"wb-task-L0/internal/cache"
	"wb-task-L0/internal/db"
	"wb-task-L0/internal/server"
	"wb-task-L0/internal/transport/kafka"
)

func main() {
	db.RunMigrations()
	dbConn := db.InitConn()

	repo := db.NewOrderRepository(dbConn)
	defer repo.Db.Close()
	go kafka.StartOrderListener()

	cacheSize := 1000
	ordersCache := cache.New(cacheSize)

	orders, err := repo.GetLatestOrders(cacheSize)
	if err != nil {
		log.Fatal(err)
	}

	ordersCache.AddAll(orders)
	server.StartServer(repo, ordersCache)
}
