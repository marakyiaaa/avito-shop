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
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Инициализация базы данных
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)

	// Инициализация сервиса аутентификации
	authService := service.NewAuthService(userRepo, cfg.JWTSecretKey)

	// Применяем миграции
	migrations.InitDB(cfg)

	// Создаем маршруты
	r := mux.NewRouter()

	// Регистрируем обработчики
	r.HandleFunc("/register", handlers.RegisterHandler(authService)).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler(authService)).Methods("POST")
	r.HandleFunc("/balance", middleware.NewCheckAuth(cfg.JWTSecretKey)(handlers.BalanceHandler(authService))).Methods("GET")
	r.HandleFunc("/send", middleware.NewCheckAuth(cfg.JWTSecretKey)(handlers.SendCoinsHandler(authService))).Methods("POST")

	// Запуск сервера
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	log.Printf("Starting server on port %s...", serverPort)
	err = http.ListenAndServe(":"+serverPort, r)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
