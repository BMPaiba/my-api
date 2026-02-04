package main

import (
	"github/mbpaiba/my-api/internal/config"
	"github/mbpaiba/my-api/internal/db"
	"github/mbpaiba/my-api/internal/db/sqlc"
	"github/mbpaiba/my-api/internal/router"
	"github/mbpaiba/my-api/internal/tcp"

	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	database, err := db.New(cfg.DB)

	if err != nil {
		logger.Fatal("error conectando DB: ", err)
	}

	defer database.Close()

	queries := sqlc.New(database)

	// Iniciar servidor TCP para dispositivo Aanviz EP300Pro
	tcpServer := tcp.NewServer(cfg.TCPPort, logger)
	go func() {
		logger.Info("Servidor TCP iniciando en " + cfg.TCPPort + " para dispositivo Aanviz EP300Pro")
		if err := tcpServer.Start(); err != nil {
			logger.Fatal("Error en servidor TCP: ", err)
		}
	}()

	// Iniciar servidor HTTP
	r := router.Setup(queries)

	logger.Info("Servidor HTTP iniciando en " + cfg.Addr)
	logger.Fatal(r.Run(cfg.Addr))
}
