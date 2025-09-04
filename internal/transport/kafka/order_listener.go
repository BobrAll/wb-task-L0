package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"strconv"
	"wb-task-L0/internal/db"
	"wb-task-L0/internal/models"
)

// StartOrderListener consumes order messages from Kafka and saves them to DB
func StartOrderListener() {
	reader := createKafkaReader()
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

// createKafkaReader configures Kafka reader from env vars
func createKafkaReader() *kafka.Reader {
	minBytes, err := strconv.Atoi(os.Getenv("KAFKA_MIN_BYTES"))
	if err != nil {
		log.Fatal("Error parsing KAFKA_MIN_BYTES", err)
	}
	maxBytes, err := strconv.Atoi(os.Getenv("KAFKA_MAX_BYTES"))
	if err != nil {
		log.Fatal("Error parsing KAFKA_MAX_BYTES", err)
	}

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{os.Getenv("KAFKA_ADDR")},
		Topic:    os.Getenv("KAFKA_ORDERS_TOPIC"),
		GroupID:  os.Getenv("KAFKA_ORDERS_GROUP_ID"),
		MinBytes: minBytes,
		MaxBytes: maxBytes,
	})
}
