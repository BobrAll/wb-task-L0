package main

import (
	"wb-task-L0/internal/db"
	"wb-task-L0/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	db.RunMigrations()

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/orders", handlers.GetOrders)
	}

	r.Run(":8080")
}
