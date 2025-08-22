package rest

import (
	"database/sql"
	"errors"
	"net/http"
	"wb-task-L0/internal/db"

	"github.com/gin-gonic/gin"
)

func GetAllOrdersIDs(repo *db.OrderRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := repo.GetAllOrdersIDs()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}

func GetOrder(repo *db.OrderRepository) gin.HandlerFunc {
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
		c.JSON(http.StatusOK, result.ToDto())
	}
}
