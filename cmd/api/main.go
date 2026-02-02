package main

import (
	"github/mbpaiba/my-api/internal/config"
	"github/mbpaiba/my-api/internal/db"
	"github/mbpaiba/my-api/internal/db/sqlc"
	"github/mbpaiba/my-api/internal/router"

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

	r := router.Setup(queries)

	logger.Info("server iniciando en " + cfg.Addr)
	logger.Fatal(r.Run(cfg.Addr))
}
