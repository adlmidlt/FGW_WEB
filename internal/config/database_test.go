package config

import (
	"os"
	"testing"
)

func TestConfig_LoadConfigDatabase(t *testing.T) {
	type fields struct {
		PSQL PostgresqlConfig
	}
	type args struct {
		pathToYamlFileTest string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Валидный yaml файл",
			args:    args{pathToYamlFileTest: "testdata/valid_config.yml"},
			wantErr: false,
		},
		{
			name:    "Файла не существует",
			args:    args{pathToYamlFileTest: "testdata/not_found.yml"},
			wantErr: true,
		},
		{
			name:    "Не валидное содержание файла",
			args:    args{pathToYamlFileTest: "testdata/invalid_config.yml"},
			wantErr: true,
		},
		{
			name:    "Пустой файл",
			args:    args{pathToYamlFileTest: "testdata/empty_config.yml"},
			wantErr: false,
		},
	}

	createTestData(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				PSQL: tt.fields.PSQL,
			}
			if err := c.readFromYaml(tt.args.pathToYamlFileTest); (err != nil) != tt.wantErr {
				t.Errorf("readFromYaml() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := c.LoadConfigDatabase(tt.args.pathToYamlFileTest); (err != nil) != tt.wantErr {
				t.Errorf("LoadConfigDatabase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	cleanupTestData(t)
}

// createTestData - создаем тестовые данные.
func createTestData(t *testing.T) {
	err := os.Mkdir("testdata", 0755)
	if err != nil {
		t.Fatalf("Ошибка создания тестовой директории: %v", err)
	}

	validConfig := `PSQL:
  host: "postgresql"
  port: 5432
  name: "test_fgw_web_db"
  user: "test_user1"
  password: "test_pass123word"
  SSLMode: "disable"
  maxOpenConns: 10
  maxIdleConns: 5
  connMaxLifetime: "5m"
  connMaxIdleTime: "1m"
 `
	err = os.WriteFile("testdata/valid_config.yml", []byte(validConfig), 0644)
	if err != nil {
		t.Fatalf("Ошибка создания валидного конфигурационного файла: %v", err)
	}

	invalidConfig := `PSQL:
  host: "postgresql"
  port: not_a_number
  name: "test_fgw_web_db"
 `
	err = os.WriteFile("testdata/invalid_config.yml", []byte(invalidConfig), 0644)
	if err != nil {
		t.Fatalf("Ошибка создания невалидного конфигурационного файла: %v", err)
	}

	err = os.WriteFile("testdata/empty_config.yml", []byte(""), 0644)
	if err != nil {
		t.Fatalf("Ошибка создания пустого конфигурационного файла: %v", err)
	}
}

// cleanupTestData - чистим тестовые данные.
func cleanupTestData(t *testing.T) {
	err := os.RemoveAll("testdata")
	if err != nil {
		t.Fatalf("Ошибка удаления папки testdata: %v", err)
	}
}
