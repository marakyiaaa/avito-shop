package migrations

import (
	"avito_shop/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"log"
)

// InitDB инициализирует базу данных и выполняет миграции.
func InitDB(config *config.Config) *migrate.Migrate {
	// Формируем строку подключения
	connStr := "postgres://" + config.DBUser + ":" + config.DBPassword + "@" + config.DBHost + ":" + config.DBPort + "/" + config.DBName + "?sslmode=disable"

	// Создаем источник миграций
	source, err := file.New("file://migrations")
	if err != nil {
		log.Fatalf("Ошибка создания источника миграций: %v", err)
	}

	// Создаем подключение к базе данных
	db, err := postgres.WithInstance(config.DBConnection(), &postgres.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	// Создаем миграцию
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", db)
	if err != nil {
		log.Fatalf("Ошибка создания миграций: %v", err)
	}

	// Запускаем миграции
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Ошибка применения миграций: %v", err)
	}

	return m
}
