package e2e_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"avito_shop/internal/handler"
	"avito_shop/internal/models/entities"
	"avito_shop/internal/repository"
	"avito_shop/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type StoreE2ETestSuite struct {
	suite.Suite
	server *httptest.Server
	db     *pgxpool.Pool
}

func (suite *StoreE2ETestSuite) SetupSuite() {
	// Запускаем контейнер с PostgreSQL
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
			"POSTGRES_DB":       "shop_test",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(30 * time.Second),
	}
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	suite.Require().NoError(err)

	// Получаем хост и порт контейнера
	host, err := postgresContainer.Host(ctx)
	suite.Require().NoError(err)
	port, err := postgresContainer.MappedPort(ctx, "5432")
	suite.Require().NoError(err)

	// Подключаемся к базе данных с использованием pgxpool
	dbURL := fmt.Sprintf("postgresql://testuser:testpassword@%s:%s/shop_test?sslmode=disable", host, port.Port())
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer pool.Close()
	fmt.Printf("Connecting to database at %s:%s\n", host, port.Port())

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("failed to ping database:", err)
	} else {
		fmt.Println("Connected to the database!")
	}

	// Применяем миграции
	err = applyMigrations(dbURL)
	suite.Require().NoError(err)

	// Инициализируем репозитории
	userRepo := repository.NewUserRepository(pool)
	itemRepo := repository.NewItemRepository(pool)
	inventoryRepo := repository.NewInventoryRepository(pool)

	// Инициализируем сервис
	storeService := service.NewStoreService(userRepo, itemRepo, inventoryRepo)

	// Инициализируем роутер и хендлер
	router := gin.Default()
	storeHandler := handler.NewStoreHandler(storeService)
	router.POST("/buy/:item", storeHandler.BuyItemHandler)

	// Запускаем тестовый сервер
	suite.server = httptest.NewServer(router)
}

func (suite *StoreE2ETestSuite) TearDownSuite() {
	// Останавливаем сервер
	suite.server.Close()

	// Закрываем соединение с базой данных
	if suite.db != nil {
		suite.db.Close()
	}
}

func (suite *StoreE2ETestSuite) TestBuyItem() {
	// Создаем пользователя
	user := &entities.User{
		Username: "testuser",
		Password: "testpass",
		Balance:  1000,
	}
	err := suite.db.QueryRow(context.Background(),
		"INSERT INTO users (username, password, balance) VALUES ($1, $2, $3) RETURNING id",
		user.Username, user.Password, user.Balance).Scan(&user.ID)
	suite.Require().NoError(err)

	// Создаем товар
	item := &entities.Item{
		Name:  "testitem",
		Price: 500,
	}
	err = suite.db.QueryRow(context.Background(),
		"INSERT INTO items (name, price) VALUES ($1, $2) RETURNING id",
		item.Name, item.Price).Scan(&item.ID)
	suite.Require().NoError(err)

	// Отправляем запрос на покупку товара
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/buy/%s", suite.server.URL, item.Name), nil)
	suite.Require().NoError(err)
	req.Header.Set("user_id", fmt.Sprintf("%d", user.ID))

	resp, err := http.DefaultClient.Do(req)
	suite.Require().NoError(err)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	// Проверяем ответ
	suite.Equal(http.StatusOK, resp.StatusCode)

	// Проверяем, что баланс пользователя уменьшился
	var balance int
	err = suite.db.QueryRow(context.Background(), "SELECT balance FROM users WHERE id = $1", user.ID).Scan(&balance)
	suite.Require().NoError(err)
	suite.Equal(500, balance)

	// Проверяем, что товар добавлен в инвентарь
	var count int
	err = suite.db.QueryRow(context.Background(), "SELECT COUNT(*) FROM inventory WHERE user_id = $1 AND item_name = $2", user.ID, item.Name).Scan(&count)
	suite.Require().NoError(err)
	suite.Equal(1, count)
}

func TestStoreE2ETestSuite(t *testing.T) {
	suite.Run(t, new(StoreE2ETestSuite))
}

func applyMigrations(dbURL string) error {
	m, err := migrate.New(
		"file://../../migrations",
		dbURL,
	)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
