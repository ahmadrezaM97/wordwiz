package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

// Config ...
type Config struct {
	App      AppCfg
	HTTP     HTTPCfg
	Postgres PostgresCfg
}

type AppCfg struct {
	Name        string
	Version     string
	Environment string // debug, dev, stage, prod
	LogLevel    string
}

func loadAppEnv() AppCfg {
	return AppCfg{
		Name:        cast.ToString(getOrReturnDefaultValue("APP_NAME", "wordwiz")),
		Version:     cast.ToString(getOrReturnDefaultValue("APP_VERSION", "1.0.0")),
		Environment: cast.ToString(getOrReturnDefaultValue("ENVIRONMENT_NAME", "debug")),
		LogLevel:    cast.ToString(getOrReturnDefaultValue("LOG_LEVEL", "debug")),
	}
}

type HTTPCfg struct {
	Host string
	Port int

	DefaultOffset string
	DefaultLimit  string

	DefaultTimeout time.Duration
}

func loadHTTPEnv() HTTPCfg {
	return HTTPCfg{
		Host:           cast.ToString(getOrReturnDefaultValue("HOST", "localhost")),
		Port:           cast.ToInt(getOrReturnDefaultValue("PORT", 8080)),
		DefaultOffset:  cast.ToString(getOrReturnDefaultValue("DEFAULT_OFFSET", "0")),
		DefaultLimit:   cast.ToString(getOrReturnDefaultValue("DEFAULT_LIMIT", "10")),
		DefaultTimeout: time.Duration(cast.ToInt64(getOrReturnDefaultValue("HTTP_DEFAULT_TIMEOUT_IN_S", 5))) * time.Second,
	}
}

type PostgresCfg struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

func loadPostgresEnv() PostgresCfg {
	return PostgresCfg{
		Host:     cast.ToString(getOrReturnDefaultValue("POSTGRES_HOST", "0.0.0.0")),
		Port:     cast.ToInt(getOrReturnDefaultValue("POSTGRES_PORT", 5432)),
		Database: cast.ToString(getOrReturnDefaultValue("POSTGRES_DATABASE", "wordwiz_db")),
		User:     cast.ToString(getOrReturnDefaultValue("POSTGRES_USER", "wordwiz_user")),
		Password: cast.ToString(getOrReturnDefaultValue("POSTGRES_PASSWORD", "wordwiz$123qwerty!")),
	}
}

func (c *PostgresCfg) MakeURL() string {
	// ex: "postgres://username:password@localhost:5432/database_name"
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
}

// Load ...
func Load() Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	config := Config{}

	config.App = loadAppEnv()

	config.HTTP = loadHTTPEnv()

	config.Postgres = loadPostgresEnv()

	return config
}

func getOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)

	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}
