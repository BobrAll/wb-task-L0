package server

import (
	"github.com/gin-gonic/gin"
	"wb-task-L0/internal/cache"
	"wb-task-L0/internal/db"
	"wb-task-L0/internal/transport/rest"
)

// StartServer initializes and runs the HTTP server
func StartServer(repo *db.OrderRepository, cache cache.Cache) {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/orders", rest.GetAllOrdersIDs(repo))
		api.GET("/orders/:order_id", rest.GetOrder(repo, cache))
	}

	r.Static("/static", "./web")
	r.NoRoute(func(c *gin.Context) {
		c.File("web/index.html")
	})

	r.Run(":8080")
}
