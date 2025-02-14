package migrations

import (
	"avito_shop/internal/config"
	"database/sql"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Импортируем источник миграций через файл
	"github.com/pressly/goose"
	"log"
)

// InitDB инициализирует соединение с базой данных и применяет миграции.
func InitDB(db *sql.DB, cfg *config.Config, migrationsPath string) {
	log.Printf("Используется база данных: %s://%s:%s@%s:%s/%s",
		"postgres", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	err := goose.Up(db, migrationsPath)
	if err != nil {
		log.Fatalf("Ошибка применения миграций: %v", err)
	}

	log.Println("Миграции успешно применены!")
}
