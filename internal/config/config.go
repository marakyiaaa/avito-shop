package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// Config представляет конфигурацию приложения.
// Содержит настройки для сервера, базы данных и JWT.
type Config struct {
	ServerPort   string
	DBPort       string
	DBUser       string
	DBName       string
	DBPassword   string
	DBHost       string
	JWTSecretKey string
}

// Load загружает конфигурацию из переменных окружения.
// Возвращает указатель на Config, заполненный значениями из окружения.
func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Файл .env не найден, используются переменные окружения из системы")
	}

	return &Config{
		ServerPort:   os.Getenv("SERVER_PORT"),
		DBPort:       os.Getenv("DATABASE_PORT"),
		DBUser:       os.Getenv("DATABASE_USER"),
		DBName:       os.Getenv("DATABASE_NAME"),
		DBPassword:   os.Getenv("DATABASE_PASSWORD"),
		DBHost:       os.Getenv("DATABASE_HOST"),
		JWTSecretKey: os.Getenv("JWT_SECRET_KEY"),
	}
}
