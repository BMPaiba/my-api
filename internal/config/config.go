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
	Addr    string
	TCPPort string
	APIRUL  string
	Env     string
	DB      DBConfig
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		panic("error cargando .env: " + err.Error())
	}

	port := fmt.Sprintf(":%v", env.GetString("PORT", ":3000"))
	
	// Configurar puerto TCP, asegurando que escuche en todas las interfaces
	tcpPortRaw := env.GetString("TCP_PORT", "8888")
	tcpPort := tcpPortRaw
	// Si solo se especifica el puerto (ej: "8888"), agregar 0.0.0.0:
	// Si ya contiene ":", asumimos que ya tiene la IP configurada
	if len(tcpPortRaw) > 0 && tcpPortRaw[0] >= '0' && tcpPortRaw[0] <= '9' {
		// Verificar si no contiene ":" (solo es un número de puerto)
		hasColon := false
		for _, c := range tcpPortRaw {
			if c == ':' {
				hasColon = true
				break
			}
		}
		if !hasColon {
			tcpPort = "0.0.0.0:" + tcpPortRaw
		}
	}

	return Config{
		Addr:    port,
		TCPPort: tcpPort,
		APIRUL:  env.GetString("EXTERNAL_URL", "localhost:3000"),
		Env:     env.GetString("ENV", "development"),
		DB: DBConfig{
			Addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
}
