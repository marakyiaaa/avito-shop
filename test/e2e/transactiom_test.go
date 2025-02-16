package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"avito_shop/internal/handler"
	"avito_shop/internal/repository"
	"avito_shop/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// setupTestDB загружает переменные окружения из .env.test и устанавливает соединение с тестовой БД.
func setupTestDB() *pgxpool.Pool {
	err := godotenv.Load("../../.env.test")
	if err != nil {
		log.Fatalf("Ошибка загрузки .env.test: %v", err)
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
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

func TestSendCoinsE2E(t *testing.T) {
	db := setupTestDB()
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	txService := service.NewTransactionService(userRepo, transactionRepo)
	txHandler := handler.NewTransactionHandler(txService)

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("user_id", 1)
		c.Next()
	})
	router.POST("/sendCoin", txHandler.SendCoinHandler)

	server := httptest.NewServer(router)
	defer server.Close()

	ctx := context.Background()
	_, err := db.Exec(ctx, "INSERT INTO users (id, username, balance) VALUES (1, 'sender', 1000)")
	if err != nil {
		t.Fatalf("Ошибка при вставке отправителя: %v", err)
	}
	_, err = db.Exec(ctx, "INSERT INTO users (id, username, balance) VALUES (2, 'recipient', 500)")
	if err != nil {
		t.Fatalf("Ошибка при вставке получателя: %v", err)
	}

	payload := map[string]interface{}{
		"to_user": "2",
		"amount":  200,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Ошибка при маршалинге JSON: %v", err)
	}

	url := server.URL + "/sendCoin"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatalf("Ошибка создания запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var senderBalance, recipientBalance int
	err = db.QueryRow(ctx, "SELECT balance FROM users WHERE id = 1").Scan(&senderBalance)
	if err != nil {
		t.Fatalf("Ошибка проверки баланса отправителя: %v", err)
	}
	err = db.QueryRow(ctx, "SELECT balance FROM users WHERE id = 2").Scan(&recipientBalance)
	if err != nil {
		t.Fatalf("Ошибка проверки баланса получателя: %v", err)
	}
	assert.Equal(t, 800, senderBalance)
	assert.Equal(t, 700, recipientBalance)
}
