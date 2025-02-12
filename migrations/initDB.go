package migrations

import (
	"avito_shop/internal/config"
	"database/sql"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Импортируем источник миграций через файл
	"github.com/pressly/goose"
	"log"
)

func InitDB(db *sql.DB, cfg *config.Config, migrationsPath string) {
	//// Подключение к базе данных
	//db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
	//	cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
	//if err != nil {
	//	log.Fatalf("Ошибка подключения к базе данных: %v", err)
	//}
	//defer db.Close()

	// Логирование информации о базе данных из конфигурации
	log.Printf("Используется база данных: %s://%s:%s@%s:%s/%s",
		"postgres", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	err := goose.Up(db, migrationsPath)
	if err != nil {
		log.Fatalf("Ошибка применения миграций: %v", err)
	}

	log.Println("Миграции успешно применены!")
}

//func InitDB(cfg *config.Config, migrationsPath string) {
//	// Строка подключения к базе данных
//	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
//		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
//
//	// Подключаемся к базе данных
//	db, err := sql.Open("postgres", dbURL)
//	if err != nil {
//		log.Fatalf("Ошибка подключения к базе данных: %v", err)
//	}
//	defer db.Close()
//
//	// Выполняем миграции
//	err = goose.Up(db, migrationsPath)
//	if err != nil {
//		log.Fatalf("Ошибка применения миграций: %v", err)
//	}
//
//	log.Println("Миграции успешно применены!")
//}
