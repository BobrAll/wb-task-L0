package rest

import (
	"database/sql"
	"errors"
	"net/http"
	"wb-task-L0/internal/db"

	"github.com/gin-gonic/gin"
)

func RegisterHandlers(repo *db.OrderRepository) {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/orders", getAllOrdersIDs(repo))
		api.GET("/orders/:order_id", getOrder(repo))
	}

	r.Run(":8080")
}

func getAllOrdersIDs(repo *db.OrderRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := repo.GetAllOrdersIDs()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}

func getOrder(repo *db.OrderRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderID := c.Param("order_id")
		result, err := repo.GetOrder(orderID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Internal server error",
					"details": err.Error(),
				})
			}
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
