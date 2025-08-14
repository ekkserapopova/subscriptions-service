package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/fx"
	"log/slog"
)

type PostgresParams struct {
	fx.In

	Cfg    Config
	Logger *slog.Logger
}

func getConnStr(cfg *Config) string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
	)
}

func NewPostgresPool(p PostgresParams) (*pgxpool.Pool, error) {
	connStr := getConnStr(&p.Cfg)

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		p.Logger.Error("pgxpool parse config: " + err.Error())
		return nil, fmt.Errorf("failed to parse connStr: %w", err)
	}

	poolConfig.MaxConns = 10

	ctx, cancel := context.WithTimeout(context.Background(), p.Cfg.ConnectTimeout)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		p.Logger.Error("pgxpool create pool: " + err.Error())
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		p.Logger.Error("pgxpool ping: " + err.Error())
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	p.Logger.Info("created pgx pool")
	return pool, nil
}

func NewPostgresConnect(p PostgresParams) (*sql.DB, error) {
	connStr := getConnStr(&p.Cfg)

	pgxConfig, err := pgx.ParseConfig(connStr)
	if err != nil {
		p.Logger.Error("pgxConfig parse config: " + err.Error())
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	db := stdlib.OpenDB(*pgxConfig)

	ctx, cancel := context.WithTimeout(context.Background(), p.Cfg.ConnectTimeout)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		p.Logger.Error("db ping failed: " + err.Error())
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}
	p.Logger.Info("connected to PostgreSQL")
	return db, nil
}
