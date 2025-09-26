package database

import (
	"context"
	"errors"
	"fgw_web/internal/config"
	"fgw_web/logs"
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

// NewPgxPool создаёт пул соединений pgxpool на основе заданной конфигурации PSQL.
func NewPgxPool(ctx context.Context, cfg *config.PostgresqlConfig) (*pgxpool.Pool, error) {
	if ctx == nil {
		return nil, errors.New(logs.E3003)
	}

	if cfg == nil {
		return nil, errors.New(logs.E3004)
	}

	pool, err := newPoolConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", logs.E3005, err)
	}

	return pool, nil
}

// newPoolConfig создаёт конфигурацию пула pgxpool на основе конфигурации PSQL.
func newPoolConfig(ctx context.Context, cfg *config.PostgresqlConfig) (*pgxpool.Pool, error) {
	connStr := formatPSQLConnStr(cfg)
	poolCfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", logs.E3000, err)
	}

	poolCfg.MaxConns = maxConns
	poolCfg.MinConns = minConns
	poolCfg.MaxConnLifetime = maxConnLifetime
	poolCfg.MaxConnIdleTime = maxConnIdleTime

	return pgxpool.NewWithConfig(ctx, poolCfg)
}

// formatPSQLConnStr возвращает строку подключения к PSQL в формате URI.
func formatPSQLConnStr(cfg *config.PostgresqlConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
}
