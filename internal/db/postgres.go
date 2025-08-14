package db

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"wb-task-L0/internal/models"
)

func init() {
	_, currFile, _, _ := runtime.Caller(0)
	currDir := filepath.Dir(currFile)
	envPath := filepath.Join(currDir, "..", "..", "configs", ".env")
	if err := godotenv.Load(envPath); err != nil {
		panic(fmt.Sprintf("Error loading .env file: %v", err))
	}
}

func getDBConnStr() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_SSL_MODE"),
	)
}

func RunMigrations() {
	m, err := migrate.New(
		"file://migrations",
		getDBConnStr(),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
}

func GetOrders(orderUID string) (models.Order, error) {
	db, err := sqlx.Connect("postgres", getDBConnStr())
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer db.Close()

	orderQuery := `
		SELECT
			o.order_uid,
			o.track_number,
			o.entry,
			o.locale,
			o.internal_signature,
			o.customer_id,
			o.delivery_service,
			o.shard_key,
			o.sm_id,
			o.date_created,
			o.oof_shard,
		
			d.name    AS "delivery.name",
			d.phone   AS "delivery.phone",
			d.zip     AS "delivery.zip",
			d.city    AS "delivery.city",
			d.address AS "delivery.address",
			d.region  AS "delivery.region",
			d.email   AS "delivery.email",
		
			p.transaction   AS "payment.transaction",
			p.request_id    AS "payment.request_id",
			p.currency      AS "payment.currency",
			p.provider      AS "payment.provider",
			p.amount        AS "payment.amount",
			p.payment_dt    AS "payment.payment_dt",
			p.bank          AS "payment.bank",
			p.delivery_cost AS "payment.delivery_cost",
			p.goods_total   AS "payment.goods_total",
			p.custom_fee    AS "payment.custom_fee"
		
		FROM orders o
		JOIN deliveries d ON o.delivery_id = d.id
		JOIN payments   p ON o.payment_transaction = p.transaction
		WHERE o.order_uid = $1
`

	itemsQuery := `
		SELECT *
		FROM items
		WHERE order_id = $1
`
	var order models.Order
	err = db.Get(&order, orderQuery, orderUID)
	if err != nil {
		log.Fatal("failed to select order:", err)
	}

	var items []models.Item

	err = db.Select(&items, itemsQuery, orderUID)
	if err != nil {
		return models.Order{}, err
	}
	order.Items = items

	return order, nil
}
