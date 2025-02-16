package e2e_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"avito_shop/internal/repository"
	"avito_shop/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var testDB *pgxpool.Pool

func setupTransactionTestDB() *pgxpool.Pool {
	err := godotenv.Load("../../.env.test")
	if err != nil {
		log.Fatalf("Ошибка загрузки .env.test: %v", err)
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)
	log.Printf("DB_USER=%s, DB_PASSWORD=%s, DB_HOST=%s, DB_PORT=%s, DB_NAME=%s",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_NAME"),
	)

	var db *pgxpool.Pool
	for i := 0; i < 10; i++ {
		db, err = pgxpool.New(context.Background(), dbURL)
		if err == nil {
			break
		}
		log.Println("Ожидание запуска БД...")
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("Не удалось подключиться к тестовой БД: %v", err)
	}
	return db
}
func TestBuyItemE2E(t *testing.T) {
	testDB = setupTransactionTestDB()
	defer testDB.Close()

	router := gin.Default()
	userRepo := repository.NewUserRepository(testDB)
	itemRepo := repository.NewItemRepository(testDB)
	inventoryRepo := repository.NewInventoryRepository(testDB)
	storeService := service.NewStoreService(userRepo, itemRepo, inventoryRepo)

	router.POST("/buy/:item", func(c *gin.Context) {
		userID := 1
		itemName := c.Param("item")
		err := storeService.BuyItem(c, userID, itemName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Покупка успешна"})
	})

	server := httptest.NewServer(router)
	defer server.Close()

	_, err := testDB.Exec(context.Background(), "INSERT INTO users (id, username, balance) VALUES (1, 'test_user', 1000)")
	if err != nil {
		t.Fatalf("Ошибка при вставке данных пользователя: %v", err)
	}
	_, err = testDB.Exec(context.Background(), "INSERT INTO items (id, name, price) VALUES (1, 'sword', 50)")
	if err != nil {
		t.Fatalf("Ошибка при вставке данных предмета: %v", err)
	}

	url := server.URL + "/buy/сup"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(nil))
	if err != nil {
		t.Fatalf("Ошибка покупки: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var newBalance int
	err = testDB.QueryRow(context.Background(), "SELECT balance FROM users WHERE id = 1").Scan(&newBalance)
	if err != nil {
		t.Fatalf("Ошибка при проверке баланса: %v", err)
	}
	assert.Equal(t, 50, newBalance)

	var itemCount int
	err = testDB.QueryRow(context.Background(), "SELECT COUNT(*) FROM inventory WHERE user_id = 1 AND item_type = 'sword'").Scan(&itemCount)
	if err != nil {
		t.Fatalf("Ошибка при проверке предмета в инвентаре: %v", err)
	}
	assert.Equal(t, 1, itemCount)
}
