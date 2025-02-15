package migrations_test

//import (
//	"avito_shop/internal/config"
//	"database/sql"
//	"errors"
//	"github.com/DATA-DOG/go-sqlmock"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"log"
//	"testing"
//)
//
//// MockGoose мок для goose.
//type MockGoose struct {
//	mock.Mock
//}
//
//func (m *MockGoose) Up(db *sql.DB, dir string) error {
//	args := m.Called(db, dir)
//	return args.Error(0)
//}
//
//// MockWriter мок для io.Writer.
//type MockWriter struct {
//	mock.Mock
//}
//
//func (m *MockWriter) Write(p []byte) (n int, err error) {
//	args := m.Called(p)
//	return args.Int(0), args.Error(1)
//}
//
//func TestInitDB(t *testing.T) {
//	db, _, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("Ошибка создания мока для sql.DB: %v", err)
//	}
//	defer func(db *sql.DB) {
//		err := db.Close()
//		if err != nil {
//		}
//	}(db)
//
//	mockGoose := new(MockGoose)
//	mockWriter := new(MockWriter)
//
//	originalLogger := log.Default()
//	log.SetOutput(mockWriter)
//	defer func() {
//		log.SetOutput(originalLogger.Writer())
//	}()
//
//	cfg := &config.Config{
//		DBUser:     "user",
//		DBPassword: "password",
//		DBHost:     "localhost",
//		DBPort:     "5432",
//		DBName:     "testdb",
//	}
//	migrationsPath := "migrations"
//
//	InitDB := func(db *sql.DB, cfg *config.Config, migrationsPath string) {
//		log.Printf("Используется база данных: %s://%s:%s@%s:%s/%s",
//			"postgres", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
//
//		err := mockGoose.Up(db, migrationsPath)
//		if err != nil {
//			log.Panicf("Ошибка применения миграций: %v", err)
//		}
//
//		log.Println("Миграции успешно применены!")
//	}
//
//	t.Run("успешное применение миграций", func(t *testing.T) {
//		mockGoose.On("Up", db, migrationsPath).Return(nil)
//		mockWriter.On("Write", mock.Anything).Return(0, nil).Twice()
//		InitDB(db, cfg, migrationsPath)
//
//		mockGoose.AssertExpectations(t)
//		mockWriter.AssertExpectations(t)
//	})
//
//	t.Run("ошибка применения миграций", func(t *testing.T) {
//		mockGoose.On("Up", db, migrationsPath).Return(errors.New("ошибка миграции"))
//		mockWriter.On("Write", mock.Anything).Return(0, nil).Once()
//
//		assert.Panics(t, func() {
//			InitDB(db, cfg, migrationsPath)
//		})
//		mockGoose.AssertExpectations(t)
//		mockWriter.AssertExpectations(t)
//	})
//}
