package database

import (
	"context"
	"database/sql"
	"fgw_web/internal/config"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"strings"
	"testing"
	"time"
)

func TestNewPostgresDB_WithTestContainer(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := settingTestContainerForPSQL(ctx)
	require.NoError(t, err)
	defer func() {
		if err = pgContainer.Terminate(ctx); err != nil {
			t.Logf("Не удалось остановить контейнер Postgres: %v", err)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx)
	require.NoError(t, err)

	if !strings.Contains(connStr, "sslmode=") {
		connStr += " sslmode=disable"
	}

	host, err := pgContainer.Host(ctx)
	require.NoError(t, err)

	port, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	cfg := &config.Config{PSQL: config.PostgresqlConfig{
		Host:     host,
		Port:     port.Int(),
		Name:     "test_db_name",
		User:     "test_name",
		Password: "test_password",
	}}

	t.Run("Успешное подключение к БД", func(t *testing.T) {
		db, err := sql.Open("postgres", connStr)
		require.NoError(t, err)
		require.NotNil(t, db)

		err = db.PingContext(ctx)
		require.NoError(t, err, "Ping должен работать с sslmode=disable")

		defer func() {
			if err = db.Close(); err != nil {
				t.Logf("failed to close postgres connection: %v", err)
			}
		}()
	})

	t.Run("Не удалось подключиться к БД", func(t *testing.T) {
		fakeDriverName := "postgres1"
		_, err := sql.Open(fakeDriverName, connStr)

		assert.Equal(t, err.Error(), fmt.Sprintf("sql: unknown driver \"%s\" (forgotten import?)", fakeDriverName))
	})

	t.Run("Успешное подключение через NewPostgresDB", func(t *testing.T) {
		db, err := NewPostgresDB(cfg)
		require.NoError(t, err)
		require.NotNil(t, db)

		err = db.PingContext(ctx)
		assert.NoError(t, err)
		defer func() {
			if err = db.Close(); err != nil {
				t.Logf("Не удалось закрыть соединение с БД: %v", err)
			}
		}()
	})
}

func settingTestContainerForPSQL(ctx context.Context) (*postgres.PostgresContainer, error) {
	pgContainer, err := postgres.Run(ctx, "postgres:15-alpine",
		postgres.WithDatabase("test_db_name"),
		postgres.WithUsername("test_name"),
		postgres.WithPassword("test_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	return pgContainer, err
}

func Test_connectionStrToDB(t *testing.T) {
	t.Run("Валидная конфигурация", func(t *testing.T) {
		got := connectionStrToDB(validConfigDB())

		assert.Equal(t, "host=postgresql port=5432 user=test_user password=test_password dbname=test_db_name sslmode=disable", got)
	})

	t.Run("Невалидная конфигурация", func(t *testing.T) {
		got := connectionStrToDB(invalidConfigDB())

		assert.Equal(t, "host=not_localhost port=0 user= password=test_password dbname= sslmode=disable", got)
		assert.NotEqual(t, validConfigDB(), got)
	})

	t.Run("Пустая конфигурация", func(t *testing.T) {
		cfg := &config.Config{PSQL: config.PostgresqlConfig{}}
		got := connectionStrToDB(cfg)

		assert.Equal(t, "host= port=0 user= password= dbname= sslmode=disable", got)
	})
}

func Test_connectionStrPGXPOOLToDB(t *testing.T) {
	t.Run("Валидная конфигурация", func(t *testing.T) {
		got := connectionStrPGXPOOLToDB(validConfigDB())

		assert.Equal(t, "postgres://test_user:test_password@postgresql:5432/test_db_name?sslmode=disable", got)
	})

	t.Run("Невалидная конфигурация", func(t *testing.T) {
		got := connectionStrPGXPOOLToDB(invalidConfigDB())

		assert.NotEqual(t, validConfigDB(), got)
	})

	t.Run("Пустая конфигурация", func(t *testing.T) {
		cfg := &config.Config{}

		assert.Equal(t, defaultConfigDB(), cfg)
	})
}

func Test_createCfgPool(t *testing.T) {
	t.Run("Успешное создание конфигурационного пула", func(t *testing.T) {
		got, err := createCfgPool(validConfigDB())
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})

	t.Run("Конфигурационный файл не найден", func(t *testing.T) {
		got, err := createCfgPool(nil)
		assert.Error(t, err, "конфигурация не найдена")
		assert.Nil(t, got)
	})

	t.Run("Не удалось распарсить строку подключения к БД", func(t *testing.T) {
		got, err := createCfgPool(invalidConfigDB())
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func Test_createPool(t *testing.T) {
	t.Run("Успешное создание пула", func(t *testing.T) {
		got, err := createPool(context.Background(), validConfigDB())
		assert.NoError(t, err)
		assert.NotNil(t, got)
	})

	t.Run("Невалидный конфигурационный файл", func(t *testing.T) {
		got, err := createPool(context.Background(), invalidConfigDB())
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestNewPgxPool(t *testing.T) {
	ctx := context.Background()
	pgContainer, err := settingTestContainerForPSQL(ctx)
	require.NoError(t, err)
	defer func() {
		if err = pgContainer.Terminate(ctx); err != nil {
			t.Logf("Не удалось остановить контейнер Postgres: %v", err)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx)
	require.NoError(t, err)

	if !strings.Contains(connStr, "sslmode=") {
		connStr += " sslmode=disable"
	}

	host, err := pgContainer.Host(ctx)
	require.NoError(t, err)

	port, err := pgContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	cfg := &config.Config{PSQL: config.PostgresqlConfig{
		Host:     host,
		Port:     port.Int(),
		Name:     "test_db_name",
		User:     "test_name",
		Password: "test_password",
	}}

	invCfg := &config.Config{PSQL: config.PostgresqlConfig{
		Host:     "not_localhost",
		Port:     port.Int(),
		Name:     "test_db_name",
		User:     "test_name",
		Password: "test_password",
	}}

	t.Run("Ошибка при создание пула", func(t *testing.T) {
		got, err := NewPgxPool(context.Background(), invalidConfigDB())

		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("Успешное создание пула", func(t *testing.T) {
		got, err := NewPgxPool(context.Background(), cfg)
		assert.NoError(t, err)
		assert.NotNil(t, got)

		err = got.Ping(ctx)
		assert.NoError(t, err)
		got.Close()
	})

	t.Run("Не удалось подключиться к БД", func(t *testing.T) {
		got, err := NewPgxPool(ctx, invCfg)
		require.Error(t, err)
		require.Nil(t, got)
	})
}

func invalidConfigDB() *config.Config {
	return &config.Config{
		PSQL: config.PostgresqlConfig{
			Host:     "not_localhost",
			Port:     0,
			Name:     "",
			User:     "",
			Password: "test_password",
		},
	}
}

func validConfigDB() *config.Config {
	return &config.Config{
		PSQL: config.PostgresqlConfig{
			Host:     "postgresql",
			Port:     5432,
			Name:     "test_db_name",
			User:     "test_user",
			Password: "test_password",
		},
	}
}

func defaultConfigDB() *config.Config {
	return &config.Config{
		PSQL: config.PostgresqlConfig{
			Host:            "",
			Port:            0,
			Name:            "",
			User:            "",
			Password:        "",
			SSLMode:         "",
			MaxOpenConns:    0,
			MaxIdleConns:    0,
			ConnMaxLifetime: 0,
			ConnMaxIdleTime: 0,
		},
	}
}
