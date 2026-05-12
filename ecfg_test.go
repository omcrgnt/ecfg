package ecfg_test

import (
	"os"
	"testing"
	"time"

	"github.com/omcrgnt/ecfg"
)

func TestParse_Success(t *testing.T) {
	type Extra struct {
		ID int `ecfg:"ID"`
	}

	type Config struct {
		AppID   string        `ecfg:"APP"`
		Timeout time.Duration `ecfg:"TIMEOUT"`
		DB      *Extra        `ecfg:"DATABASE"`
	}

	// Устанавливаем окружение
	os.Setenv("APP", "test-app")
	os.Setenv("TIMEOUT", "5s")
	os.Setenv("DATABASE_ID", "42")
	defer os.Clearenv()

	cfg, err := ecfg.Parse[Config]()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.AppID != "test-app" {
		t.Errorf("Expected APP to be 'test-app', got %s", cfg.AppID)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("Expected TIMEOUT to be 5s, got %v", cfg.Timeout)
	}
	if cfg.DB == nil || cfg.DB.ID != 42 {
		t.Errorf("Expected DATABASE_ID to be 42, got %v", cfg.DB)
	}
}

func TestParse_ValidationErrors(t *testing.T) {
	t.Run("missing ecfg tag at root", func(t *testing.T) {
		type BadConfig struct {
			NoTag string // Нет тега ecfg
		}
		_, err := ecfg.Parse[BadConfig]()
		if err == nil {
			t.Error("Expected error due to missing mandatory ecfg tag at root, but got nil")
		}
	})

	t.Run("invalid integer type", func(t *testing.T) {
		type Config struct {
			Port int `ecfg:"PORT"`
		}
		os.Setenv("PORT", "not-a-number")
		defer os.Unsetenv("PORT")

		_, err := ecfg.Parse[Config]()
		if err == nil {
			t.Error("Expected error for invalid integer format, but got nil")
		}
	})

	t.Run("invalid duration format", func(t *testing.T) {
		type Config struct {
			TTL time.Duration `ecfg:"TTL"`
		}
		os.Setenv("TTL", "100 лет") // Невалидный формат для time.ParseDuration
		defer os.Unsetenv("TTL")

		_, err := ecfg.Parse[Config]()
		if err == nil {
			t.Error("Expected error for invalid duration format, but got nil")
		}
	})
}

func TestParse_NestedLogic(t *testing.T) {
	type Deep struct {
		Key string `ecfg:"KEY"`
	}
	type Mid struct {
		Inner Deep // Без тега, должно взяться имя поля Inner
	}
	type Root struct {
		Module Mid `ecfg:"MODULE"`
	}

	// Путь должен быть MODULE_INNER_KEY
	os.Setenv("MODULE_INNER_KEY", "secret")
	defer os.Unsetenv("MODULE_INNER_KEY")

	cfg, err := ecfg.Parse[Root]()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Module.Inner.Key != "secret" {
		t.Errorf("Expected secret, got %s", cfg.Module.Inner.Key)
	}
}

func TestParse_WithPrefix(t *testing.T) {
	type Config struct {
		Port int `ecfg:"PORT"`
	}

	os.Setenv("MYAPP_PORT", "8080")
	defer os.Unsetenv("MYAPP_PORT")

	// Передаем опцию WithPrefix
	cfg, err := ecfg.Parse[Config](ecfg.WithPrefix("MYAPP"))
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Port != 8080 {
		t.Errorf("Expected 8080, got %d", cfg.Port)
	}
}
