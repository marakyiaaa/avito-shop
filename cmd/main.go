package main

import (
	"avito_shop/internal/config"
	"avito_shop/internal/handlers"
	"avito_shop/internal/middlware"
	"avito_shop/internal/repository"
	"avito_shop/internal/service"
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()

	// Подключение к базе данных
	db, err := sql.Open("postgres", "user="+cfg.DBUser+" password="+cfg.DBPassword+" dbname="+cfg.DBName+" host="+cfg.DBHost+" port="+cfg.DBPort+" sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Инициализация репозиториев и сервисов
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, "yourSecretKey")

	// Инициализация обработчиков
	authHandlers := handlers.NewAuthHandlers(authService)

	// Настройка маршрутов
	router := mux.NewRouter()
	router.HandleFunc("/register", authHandlers.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", authHandlers.LoginHandler).Methods("POST")

	// Пример защищенного маршрута
	protected := router.PathPrefix("/protected").Subrouter()
	protected.Use(middlware.NewCheckAuth("yourSecretKey"))
	protected.HandleFunc("/user-balance", authHandlers.GetUserBalanceHandler).Methods("GET")

	// Запуск сервера
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, router))
}
