package main

import (
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
	server.StartServer(repo)
}
