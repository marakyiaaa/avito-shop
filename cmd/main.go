package main

import (
	"avito_shop/internal/repository"
	"database/sql"
)

func main() {
	// Подключение к базе данных
	connStr := "user=postgres dbname=shop password=password sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)

}
