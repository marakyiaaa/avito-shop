package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"avito_shop/internal/config"
	"avito_shop/internal/handler"
	"avito_shop/internal/middleware"
	"avito_shop/internal/repository"
	"avito_shop/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()
	if err != nil {
		log.Printf("Ошибка при запуске сервера: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	itemRepo := repository.NewItemRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecretKey)
	storeService := service.NewStoreService(userRepo, itemRepo, inventoryRepo)
	transactionService := service.NewTransactionService(userRepo, transactionRepo)
	infoService := service.NewInfoService(userRepo, itemRepo, inventoryRepo, transactionRepo)

	authHandlers := handler.NewAuthHandlers(authService)
	storeHandlers := handler.NewStoreHandler(storeService)
	transactionHandlers := handler.NewTransactionHandler(transactionService)
	infoHandler := handler.NewInfoHandler(infoService)

	r := gin.Default()
	r.Use(middleware.NewCORS())

	r.GET("/api/info", middleware.NewCheckAuth(cfg.JWTSecretKey), infoHandler.GetUserInfoHandler)
	r.POST("/api/auth", authHandlers.AuthHandler)
	r.POST("/api/sendCoin", middleware.NewCheckAuth(cfg.JWTSecretKey), transactionHandlers.SendCoinHandler)
	r.POST("/api/buy/:item", middleware.NewCheckAuth(cfg.JWTSecretKey), storeHandlers.BuyItemHandler)

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	log.Printf("Сервер запущен на порту %s...", serverPort)
	err = r.Run(":" + serverPort)
	if err != nil {
		log.Printf("Ошибка при запуске сервера: %v", err)
	}
}
