package main

import (
	"avito_shop/internal/config"
	"avito_shop/internal/handlers"
	"avito_shop/internal/middleware"
	"avito_shop/internal/repository"
	"avito_shop/internal/service"
	"avito_shop/migrations"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Подключение БД
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Инициализация БД и применение миграций
	migrations.InitDB(db, cfg, "./migrations")

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)
	itemRepo := repository.NewItemRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	// Инициализация сервисов
	authService := service.NewAuthService(userRepo, cfg.JWTSecretKey)
	storeService := service.NewStoreService(userRepo, itemRepo)
	transactionService := service.NewTransactionService(userRepo, transactionRepo)
	infoService := service.NewInfoService(userRepo, itemRepo, transactionRepo)

	// Инициализация обработчиков
	authHandlers := handlers.NewAuthHandlers(authService)
	storeHandlers := handlers.NewStoreHandler(storeService)
	transactionHandlers := handlers.NewTransactionHandler(transactionService)
	infoHandler := handlers.NewInfoHandler(infoService)

	r := gin.Default()
	r.Use(middleware.NewCORS())

	// Регистрируем обработчики
	r.POST("/api/auth", authHandlers.AuthHandler)
	r.GET("/api/info", middleware.NewCheckAuth(cfg.JWTSecretKey), infoHandler.GetUserInfoHandler)
	r.POST("/api/sendCoin", middleware.NewCheckAuth(cfg.JWTSecretKey), transactionHandlers.SendCoinHandler)
	r.GET("/api/buy/:item", middleware.NewCheckAuth(cfg.JWTSecretKey), storeHandlers.BuyItemHandler) //пост или гет?

	// Определяем порт сервера
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	log.Printf("Сервер запущен на порту %s...", serverPort)
	err = r.Run(":" + serverPort)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
