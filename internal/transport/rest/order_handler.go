package rest

import (
	"database/sql"
	"errors"
	"net/http"
	"wb-task-L0/internal/db"

	"github.com/gin-gonic/gin"
)

func GetOrder(c *gin.Context) {
	orderID := c.Param("order_id")
	result, err := db.GetOrders(orderID)
	if err == nil {
		c.JSON(http.StatusOK, result)
	} else {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}
}
