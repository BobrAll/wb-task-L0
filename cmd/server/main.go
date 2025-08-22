package main

import (
	"wb-task-L0/internal/db"
	"wb-task-L0/internal/server"
)

func main() {
	db.RunMigrations()
	dbConn := db.InitConn()

	repo := db.NewOrderRepository(dbConn)
	defer repo.Db.Close()
	server.StartServer(repo)
}
