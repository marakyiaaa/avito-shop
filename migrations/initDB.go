package migrations

import (
	"avito_shop/internal/config"
	"database/sql"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Импортируем источник миграций через файл
	"github.com/pressly/goose"
	"log"
)

func InitDB(cfg *config.Config, migrationsPath string) {
	// Подключение к базе данных
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	err = goose.Up(db, migrationsPath)
	if err != nil {
		log.Fatalf("Ошибка применения миграций: %v", err)
	}

	log.Println("Миграции успешно применены!")
}
