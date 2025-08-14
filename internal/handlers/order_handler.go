package handlers

import (
	"log"
	"net/http"
	"wb-task-L0/internal/db"

	"github.com/gin-gonic/gin"
)

func GetOrders(c *gin.Context) {
	result, err := db.GetOrders("b563feb7b2b84b6test")
	if err == nil {
		c.JSON(http.StatusOK, result)
	} else {
		log.Fatal(err)
	}
}
