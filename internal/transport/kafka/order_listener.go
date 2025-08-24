package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"wb-task-L0/internal/db"
	"wb-task-L0/internal/models"
)

func StartOrderListener() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "orders",
		GroupID: "orders-group",
	})
	defer reader.Close()

	dbConn := db.InitConn()
	repo := db.NewOrderRepository(dbConn)

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Reading kafka message error (orders topic): %v", err)
			continue
		}

		var order models.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			log.Println("Error parsing kafka message (new order): %v, data: %s", err, string(msg.Value))
			continue
		}

		err = repo.SaveOrder(order)
		if err != nil {
			log.Println("SaveOrder error: %v, data: %s", err, string(msg.Value))
		}
	}
}
