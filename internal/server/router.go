package server

import (
	"github.com/gin-gonic/gin"
	"wb-task-L0/internal/db"
	"wb-task-L0/internal/transport/rest"
)

func StartServer(repo *db.OrderRepository) {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/orders", rest.GetAllOrdersIDs(repo))
		api.GET("/orders/:order_id", rest.GetOrder(repo))
	}

	r.NoRoute(func(c *gin.Context) {
		c.File("web/index.html")
	})

	r.Run(":8080")
}
