package config

import (
	"fmt"
	"github/mbpaiba/my-api/internal/env"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type Config struct {
	Addr   string
	APIRUL string
	Env    string
	DB     DBConfig
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		panic("error cargando .env: " + err.Error())
	}

	port := fmt.Sprintf(":%v", env.GetString("PORT", ":3000"))

	return Config{
		Addr:   port,
		APIRUL: env.GetString("EXTERNAL_URL", "localhost:3000"),
		Env:    env.GetString("ENV", "development"),
		DB: DBConfig{
			Addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
}
