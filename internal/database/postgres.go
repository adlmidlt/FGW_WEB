package database

import (
	"database/sql"
	"fgw_web/internal/config"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

const (
	maxConns          = 20
	minConns          = 5
	maxConnLifetime   = time.Hour
	maxConnIdleTime   = time.Minute * 30
	healthCheckPeriod = time.Minute
)

func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionStrToDB(cfg))
	if err != nil {
		return nil, err
	}

	return db, nil
}

//func NewPostgresDB(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
//	configDB, err := pgxpool.ParseConfig(connectionStrToDB(cfg))
//	if err != nil {
//		return nil, fmt.Errorf("не удалось распарсить: %v", err)
//	}
//
//	pool, err := cfgNewPoolWithCtx(ctx, configDB)
//	if err != nil {
//		return nil, err
//	}
//
//	pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
//	defer pingCancel()
//
//	if err = pool.Ping(pingCtx); err != nil {
//		return nil, fmt.Errorf("ping database: %w", err)
//	}
//
//	return pool, nil
//}

// cfgNewPoolWithCtx конфигурация нового пула с контекстом.
//func cfgNewPoolWithCtx(ctx context.Context, configDB *pgxpool.Config) (*pgxpool.Pool, error) {
//	if configDB == nil {
//		return nil, errors.New("конфигурация обязательна и не может быть nil")
//	}
//
//	configDB.MaxConns = maxConns
//	configDB.MinConns = minConns
//	configDB.MaxConnLifetime = maxConnLifetime
//	configDB.MaxConnIdleTime = maxConnIdleTime
//	configDB.HealthCheckPeriod = healthCheckPeriod
//
//	pool, err := pgxpool.NewWithConfig(ctx, configDB)
//	if err != nil {
//		return nil, fmt.Errorf("ошибка создание пула: %w", err)
//	}
//
//	return pool, nil
//}

// connectionStrToDB - строка подключения к БД.
func connectionStrToDB(cfg *config.Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.PSQL.Host, cfg.PSQL.Port, cfg.PSQL.User,
		cfg.PSQL.Password, cfg.PSQL.Name,
	)
}
