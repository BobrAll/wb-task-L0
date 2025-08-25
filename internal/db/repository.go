package db

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

func (r *OrderRepository) GetOrdersIDs(search string, page int32, size int32) ([]string, int32, error) {
	query := `
        SELECT 
			ARRAY_AGG(order_uid) AS ids,
			MAX(total_count) AS total_count
		FROM (
			SELECT 
				order_uid,
				COUNT(*) OVER() AS total_count
			FROM orders
			WHERE order_uid LIKE $1
			ORDER BY date_created
			OFFSET $2 LIMIT $3
		) AS subquery;
    `

	searchPattern := "%" + search + "%"
	offset := page * size

	var idsBytes []byte
	var count int32

	err := r.Db.QueryRow(query, searchPattern, offset, size).Scan(&idsBytes, &count)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get orders: %w", err)
	}

	ids := parsePostgresArr(idsBytes)
	return ids, count, nil
}

func parsePostgresArr(arrayBytes []byte) []string {
	idsStr := string(arrayBytes)
	idsStr = strings.Trim(idsStr, "{}")
	ids := strings.Split(idsStr, ",")

	if len(ids) == 1 && ids[0] == "" {
		ids = []string{}
	}
	return ids
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
		WHERE track_number = $1
	`

	var order models.Order
	err := r.Db.Get(&order, orderQuery, orderUID)
	if err != nil {
		return models.Order{}, err
	}

	var items []models.Item

	err = r.Db.Select(&items, itemsQuery, order.TrackNumber)
	if err != nil {
		return models.Order{}, err
	}
	order.Items = items

	return order, nil
}

func (r *OrderRepository) SaveOrder(order models.Order) error {
	tx, err := r.Db.Beginx()
	if err != nil {
		return err
	}

	var deliveryID int
	err = tx.QueryRowx(`
		INSERT INTO deliveries (name, phone, zip, city, address, region, email)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id
	`,
		order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email,
	).Scan(&deliveryID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.NamedExec(`
		INSERT INTO payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES (:transaction, :request_id, :currency, :provider, :amount, :payment_dt, :bank, :delivery_cost, :goods_total, :custom_fee)
		ON CONFLICT (transaction) DO NOTHING
	`, order.Payment)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		INSERT INTO orders (order_uid, track_number, entry, delivery_id, payment_transaction,
		                    locale, internal_signature, customer_id, delivery_service,
		                    shard_key, sm_id, date_created, oof_shard)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
	`,
		order.OrderID, order.TrackNumber, order.Entry, deliveryID, order.Payment.Transaction,
		order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.ShardKey, order.SmID, order.DateCreated, order.OofShard,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(order.Items) > 0 {
		_, err = tx.NamedExec(`
			INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
			VALUES (:chrt_id, :track_number, :price, :rid, :name, :sale, :size, :total_price, :nm_id, :brand, :status)
		`, order.Items)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *OrderRepository) GetLatestOrders(n int) ([]models.Order, error) {
	query := `
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
		ORDER BY o.date_created DESC
		LIMIT $1
	`

	var orders []models.Order
	if err := r.Db.Select(&orders, query, n); err != nil {
		return nil, err
	}

	trackNumbers := make([]string, 0, len(orders))
	for _, o := range orders {
		trackNumbers = append(trackNumbers, o.TrackNumber)
	}

	itemsQuery := `
		SELECT *
		FROM items
		WHERE track_number = ANY($1)
	`
	var items []models.Item
	if err := r.Db.Select(&items, itemsQuery, pq.StringArray(trackNumbers)); err != nil {
		return nil, err
	}

	itemsMap := make(map[string][]models.Item)
	for _, item := range items {
		itemsMap[item.TrackNumber] = append(itemsMap[item.TrackNumber], item)
	}

	for i, o := range orders {
		orders[i].Items = itemsMap[o.TrackNumber]
	}

	return orders, nil
}
