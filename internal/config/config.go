package config

import (
	"os"
)

type Config struct {
	ServerPort   string
	DBPort       string
	DBUser       string
	DBName       string
	DBPassword   string
	DBHost       string
	JWTSecretKey string
}

func Load() *Config {
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
