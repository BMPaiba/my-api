package db

import (
	"context"
	"fmt"
	"github/mbpaiba/my-api/internal/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(cfg config.DBConfig) (*pgxpool.Pool, error) {
	maxIdleTime, err := time.ParseDuration(cfg.MaxIdleTime)
	if err != nil {
		return nil, fmt.Errorf("error parseando MaxIdleTime: %w", err)
	}

	poolCfg, err := pgxpool.ParseConfig(cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("error parseando config de DB: %w", err)
	}

	poolCfg.MaxConns = int32(cfg.MaxOpenConns)
	poolCfg.MinConns = int32(cfg.MaxIdleConns)
	poolCfg.MaxConnIdleTime = maxIdleTime

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("error creando pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error conectando DB: %w", err)
	}

	return pool, nil
}
