package database

import (
	"context"
	"fgw_web/internal/config"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"strings"
	"testing"
	"time"
)

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
		Port:     -1,
		Name:     "",
		User:     "test_name",
		Password: "test_password",
	}}

	t.Run("Ошибка парсинга строки подключения при создание пула", func(t *testing.T) {
		got, err := NewPgxPool(ctx, &invCfg.PSQL)

		fmt.Println(got, err)
		assert.Error(t, err)
		assert.Nil(t, got)

	})

	t.Run("Успешное создание пула", func(t *testing.T) {
		got, err := NewPgxPool(ctx, &cfg.PSQL)
		assert.NoError(t, err)
		assert.NotNil(t, got)

		err = got.Ping(ctx)
		assert.NoError(t, err)
		got.Close()
	})
}

// settingTestContainerForPSQL запускает контейнер с настроенными параметрами.
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

func Test_NewPgxPool(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{PSQL: config.PostgresqlConfig{
		Host:     "postgresql",
		Port:     5432,
		Name:     "test_db_name",
		User:     "test_name",
		Password: "test_password",
	}}
	invCfg := &config.Config{PSQL: config.PostgresqlConfig{
		Host:     "not_localhost",
		Port:     -1,
		Name:     "",
		User:     "test_name",
		Password: "test_password",
	}}

	type args struct {
		ctx context.Context
		cfg *config.PostgresqlConfig
	}
	tests := []struct {
		name    string
		args    args
		want    *pgxpool.Pool
		wantErr bool
	}{
		{
			name: "nil - контекст",
			args: args{
				ctx: nil,
				cfg: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "nil - конфиг",
			args: args{
				ctx: ctx,
				cfg: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Не удалось создать пул",
			args: args{
				ctx: ctx,
				cfg: &invCfg.PSQL,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Не удалось подключиться к БД",
			args: args{
				ctx: ctx,
				cfg: &cfg.PSQL,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPgxPool(tt.args.ctx, tt.args.cfg)

			if tt.wantErr {
				assert.Error(t, err)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("NewPgxPool() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equalf(t, tt.want, got, "NewPgxPool(%v, %v)", tt.args.ctx, tt.args.cfg)
		})
	}
}
