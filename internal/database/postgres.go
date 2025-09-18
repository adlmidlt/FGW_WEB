package database

import (
	"context"
	"fgw_web/internal/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type PostgresDB struct {
	DB *gorm.DB
}

// NewPostgresConnection - новое подключение к БД.
func NewPostgresConnection(ctx context.Context, config *config.Config) (*PostgresDB, error) {
	db := &PostgresDB{}

	if err := db.connectWithRetry(ctx, config, 5, 5*time.Second); err != nil {
		return nil, err
	}

	return db, nil
}

// SetupDatabase - создает и настраивает подключение к PSQL базе данных.
func SetupDatabase(cfg config.Config) *PostgresDB {
	ctxInit, cancelInit := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancelInit()

	psqlDBConn, err := NewPostgresConnection(ctxInit, &cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return psqlDBConn
}

// connectWithRetry - повторные попытки подключения.
func (db *PostgresDB) connectWithRetry(ctx context.Context, config *config.Config, maxAttempts int, delay time.Duration) error {
	var err error
	for i := 0; i < maxAttempts; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err = db.connect(config); err == nil {
			return nil
		}

		log.Printf("Попытка подключения к базе данных %d/%d не удалась: %v", i+1, maxAttempts, err)

		if i < maxAttempts-1 {
			time.Sleep(delay)
		}
	}

	return fmt.Errorf("не удалось подключиться после %d попыток: %w", maxAttempts, err)
}

// connect - основное подключение.
func (db *PostgresDB) connect(config *config.Config) error {
	gormDB, err := gorm.Open(
		postgres.Open(dsn(config)),
		&gorm.Config{
			Logger:      logger.Default.LogMode(logger.Info),
			PrepareStmt: true,
		},
	)
	if err != nil {
		return err
	}

	// Настройка пула соединений GORM
	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(config.PSQL.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.PSQL.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.PSQL.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.PSQL.ConnMaxIdleTime)

	db.DB = gormDB

	return nil
}

// dsn - формирование строки подключения к БД.
func dsn(cfg *config.Config) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.PSQL.Host, cfg.PSQL.User, cfg.PSQL.Password, cfg.PSQL.Name, cfg.PSQL.Port)
}

// CloseDatabase - безопасно закрывает соединение с базой данных.
func CloseDatabase(db *PostgresDB) {
	if err := db.Close(); err != nil {
		log.Printf("Предупреждение: Ну удалось закрыть соединение с БД: %v", err)
	}
}

// Close - закрывает соединение с базой данных и возвращает ошибку операции.
func (db *PostgresDB) Close() error {
	if db.DB == nil {
		return nil
	}

	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
