package database

import (
	"context"
	"database/sql"
	"errors"
	"fgw_web/internal/config"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"time"
)

const (
	maxConns        = 20
	minConns        = 5
	maxConnLifetime = time.Hour
	maxConnIdleTime = time.Minute * 30
)

func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionStrToDB(cfg))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewPgxPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := createPool(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}

func createPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	poolCfg, err := createCfgPool(cfg)
	if err != nil || poolCfg == nil {
		return nil, fmt.Errorf("ошибка создания подключения пула: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании пула: %w", err)
	}

	return pool, nil
}

func createCfgPool(cfg *config.Config) (*pgxpool.Config, error) {
	if cfg == nil {
		return nil, errors.New("конфигурация не найдена")
	}

	poolCfg, err := pgxpool.ParseConfig(connectionStrPGXPOOLToDB(cfg))
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга строки подключения: %w", err)
	}

	poolCfg.MaxConns = maxConns
	poolCfg.MinConns = minConns
	poolCfg.MaxConnLifetime = maxConnLifetime
	poolCfg.MaxConnIdleTime = maxConnIdleTime

	return poolCfg, nil
}

// connectionStrSQLToDB - строка подключения к БД.
func connectionStrToDB(cfg *config.Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PSQL.Host, cfg.PSQL.Port, cfg.PSQL.User, cfg.PSQL.Password, cfg.PSQL.Name,
	)
}

// connectionStrPGXPOOLToDB - строка подключения к БД.
func connectionStrPGXPOOLToDB(cfg *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PSQL.User, cfg.PSQL.Password, cfg.PSQL.Host, cfg.PSQL.Port, cfg.PSQL.Name)
}
