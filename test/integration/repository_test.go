package integration

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"wb-task-L0/internal/db"
	"wb-task-L0/internal/models"
)

var pgContainer *postgres.PostgresContainer

func RunMigrations(dsn string) {
	_, filename, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(filename)
	projectRoot := filepath.Join(testDir, "..", "..")
	migrationsPath := filepath.Join(projectRoot, "migrations")

	absolutePath, err := filepath.Abs(migrationsPath)
	if err != nil {
		log.Fatal("Failed to get absolute path for migrations:", err)
	}

	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		log.Fatalf("Migrations directory does not exist: %s", absolutePath)
	}

	migrationURL := "file://" + absolutePath

	m, err := migrate.New(migrationURL, dsn)
	if err != nil {
		log.Fatal("Failed to create migration instance:", err)
	}

	if err := m.Up(); err != nil {
		if err.Error() != "no change" {
			log.Fatal("Failed to run migrations:", err)
		}
	}
}

func setupPostgres(t *testing.T) *db.OrderRepository {
	ctx := context.Background()

	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	require.NoError(t, err)

	pgContainer = container
	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err)

	dsn := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())

	var dbConn *sqlx.DB
	for i := 0; i < 5; i++ {
		dbConn, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	require.NoError(t, err)

	RunMigrations(dsn)

	repo := db.NewOrderRepository(dbConn)
	return repo
}

func teardownPostgres(t *testing.T) {
	require.NoError(t, pgContainer.Terminate(context.Background()))
}

func TestOrderRepository_SaveAndGetOrder(t *testing.T) {
	repo := setupPostgres(t)
	defer teardownPostgres(t)

	order := models.Order{
		OrderID:     "order123",
		TrackNumber: "track123",
		Entry:       "web",
		Delivery: models.Delivery{
			Name:    "John Doe",
			Phone:   "123456789",
			Zip:     "12345",
			City:    "TestCity",
			Address: "Main Street",
			Region:  "TestRegion",
			Email:   "test@example.com",
		},
		Payment: models.Payment{
			Transaction:  "txn123",
			Currency:     "USD",
			Provider:     "TestPay",
			Amount:       100,
			PaymentDt:    time.Now().Unix(),
			Bank:         "TestBank",
			DeliveryCost: 10,
			GoodsTotal:   90,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      1,
				TrackNumber: "track123",
				Price:       50,
				RID:         "rid1",
				Name:        "Item1",
				Sale:        0,
				Size:        "M",
				TotalPrice:  50,
				NmId:        111,
				Brand:       "Brand1",
				Status:      1,
			},
		},
		DateCreated: time.Now(),
	}

	require.NoError(t, repo.SaveOrder(order))
	got, err := repo.GetOrder("order123")
	require.NoError(t, err)
	require.Equal(t, order.OrderID, got.OrderID)
	require.Equal(t, order.Delivery.Name, got.Delivery.Name)
	require.Len(t, got.Items, 1)
}

func TestOrderRepository_GetLatestOrders(t *testing.T) {
	repo := setupPostgres(t)
	defer teardownPostgres(t)

	order := models.Order{
		OrderID:     "order456",
		TrackNumber: "track456",
		Entry:       "mobile",
		Delivery: models.Delivery{
			Name:    "Alice",
			Phone:   "999888777",
			Zip:     "54321",
			City:    "Town",
			Address: "Side Street",
			Region:  "Region",
			Email:   "alice@example.com",
		},
		Payment: models.Payment{
			Transaction:  "txn456",
			Currency:     "EUR",
			Provider:     "PayProvider",
			Amount:       200,
			PaymentDt:    time.Now().Unix(),
			Bank:         "BankX",
			DeliveryCost: 20,
			GoodsTotal:   180,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      2,
				TrackNumber: "track456",
				Price:       200,
				RID:         "rid2",
				Name:        "Item2",
				Sale:        0,
				Size:        "L",
				TotalPrice:  200,
				NmId:        222,
				Brand:       "Brand2",
				Status:      1,
			},
		},
		DateCreated: time.Now(),
	}
	require.NoError(t, repo.SaveOrder(order))

	orders, err := repo.GetLatestOrders(5)
	require.NoError(t, err)
	require.NotEmpty(t, orders)
}

func TestOrderRepository_GetOrdersIDs(t *testing.T) {
	repo := setupPostgres(t)
	defer teardownPostgres(t)

	order := models.Order{
		OrderID:     "order789",
		TrackNumber: "track789",
		Entry:       "api",
		Delivery:    models.Delivery{Name: "Bob"},
		Payment:     models.Payment{Transaction: "txn789", Currency: "USD", Amount: 50},
		DateCreated: time.Now(),
	}
	require.NoError(t, repo.SaveOrder(order))

	ids, total, err := repo.GetOrdersIDs("order", 0, 10)
	require.NoError(t, err)
	require.Contains(t, ids, "order789")
	require.GreaterOrEqual(t, total, int32(1))
}
