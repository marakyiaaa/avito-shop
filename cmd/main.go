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

	// Подключение базы данных
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Инициализация базы данных и применение миграций
	migrations.InitDB(db, cfg, "./migrationsD")

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)

	// Инициализация сервиса аутентификации
	authService := service.NewAuthService(userRepo, cfg.JWTSecretKey)

	authHandlers := handlers.NewAuthHandlers(authService)

	r := gin.Default()

	// Регистрируем обработчики
	r.POST("/register", authHandlers.RegisterHandler)
	r.POST("/login", authHandlers.LoginHandler)
	r.GET("/balance", middleware.NewCheckAuth(cfg.JWTSecretKey), authHandlers.GetUserBalanceHandler)

	// Определяем порт сервера
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	// Запуск сервера
	log.Printf("Сервер запущен на порту %s...", serverPort)
	err = r.Run(":" + serverPort)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
