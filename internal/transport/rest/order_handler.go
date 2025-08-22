package rest

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"wb-task-L0/internal/db"

	"github.com/gin-gonic/gin"
)

func GetAllOrdersIDs(repo *db.OrderRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		search := c.Query("search")
		page, err := parseInt32(c.DefaultQuery("page", "0"))
		if err != nil || page < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Param page must be a num greater or equal to 0"})
			return
		}

		size, err := parseInt32(c.DefaultQuery("size", "10"))
		if err != nil || size < 1 || size > 50 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Param size must be a num between 1 and 50"})
			return
		}

		ordersIds, totalOrders, err := repo.GetOrdersIDs(search, page, size)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			log.Println(err)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"orders_ids":   ordersIds,
				"total_orders": totalOrders,
			})
		}
	}
}

func parseInt32(numStr string) (int32, error) {
	num, err1 := strconv.ParseInt(numStr, 10, 32)
	if err1 != nil {
		return 0, err1
	}
	return int32(num), nil
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
		c.JSON(http.StatusOK, gin.H{"order": result.ToDto()})
	}
}
