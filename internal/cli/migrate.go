package cli

import (
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/zlobste/fake-wallet/internal/assets"
	"github.com/zlobste/fake-wallet/internal/config"
)

var migrations = &migrate.PackrMigrationSource{
	Box: assets.Migrations,
}

func MigrateUp(cfg config.Config) error {
	applied, err := migrate.Exec(cfg.DB(), "postgres", migrations, migrate.Up)
	if err != nil {
		return errors.Wrap(err, "failed to apply migrations")
	}
	cfg.Logging().WithField("applied", applied).Info("migrations applied")
	return nil
}

func MigrateDown(cfg config.Config) error {
	applied, err := migrate.Exec(cfg.DB(), "postgres", migrations, migrate.Down)
	if err != nil {
		return errors.Wrap(err, "failed to apply migrations")
	}
	cfg.Logging().WithField("applied", applied).Info("migrations applied")
	return nil
}