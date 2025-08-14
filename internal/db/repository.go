package db

import (
	"errors"
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

type OrderRepository struct {
	Db *sqlx.DB
}

func InitConn() *sqlx.DB {
	db, err := sqlx.Connect("postgres", getDBConnStr())
	if err != nil {
		log.Fatal("cannot connect to db: %w", err)
	}
	return db
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func loadEnv() {
	_, currFile, _, _ := runtime.Caller(0)
	currDir := filepath.Dir(currFile)
	envPath := filepath.Join(currDir, "..", "..", "configs", ".env")
	if err := godotenv.Load(envPath); err != nil {
		panic(fmt.Sprintf("Error loading .env file: %v", err))
	}
}

func getDBConnStr() string {
	loadEnv()
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
	m, err := migrate.New("file://migrations", getDBConnStr())
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}
}

func (r *OrderRepository) GetAllOrdersIDs() ([]string, error) {
	query := `SELECT o.order_uid FROM orders o`
	var ids []string
	err := r.Db.Select(&ids, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders ids")
	}
	return ids, nil
}

func (r *OrderRepository) GetOrder(orderUID string) (models.Order, error) {
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
	err := r.Db.Get(&order, orderQuery, orderUID)
	if err != nil {
		return models.Order{}, err
	}

	var items []models.Item

	err = r.Db.Select(&items, itemsQuery, orderUID)
	if err != nil {
		return models.Order{}, err
	}
	order.Items = items

	return order, nil
}
