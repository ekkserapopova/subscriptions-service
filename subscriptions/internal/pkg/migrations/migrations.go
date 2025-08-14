package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/fx"
	"log/slog"
)

type Params struct {
	fx.In

	DB     *sql.DB
	Logger *slog.Logger
}

//go:embed schema/*.sql
var migrationFiles embed.FS

func RunMigrations(params Params) error {
	sourceDriver, err := iofs.New(migrationFiles, "schema")
	if err != nil {
		params.Logger.Error("failed to load migration files: ", err.Error())
		return fmt.Errorf("failed to initialize migrations source driver: %w", err)
	}

	dbDriver, err := postgres.WithInstance(params.DB, &postgres.Config{})
	if err != nil {
		params.Logger.Error("failed to initialize postgres driver: ", err.Error())
		return fmt.Errorf("failed to initialize postgres driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		params.Logger.Error("failed to initialize migrate instance: ", err.Error())
		return fmt.Errorf("failed to initialize migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		params.Logger.Error("failed to run migrations: ", err.Error())
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	err = sourceDriver.Close()
	if err != nil {
		params.Logger.Error("failed to close source driver: ", err.Error())
		return fmt.Errorf("failed to close source driver: %w", err)
	}

	err = dbDriver.Close()
	if err != nil {
		params.Logger.Error("failed to close db driver: ", err.Error())
		return fmt.Errorf("failed to close db driver: %w", err)
	}

	params.Logger.Info("migrations done")
	return nil
}
