package database

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"github.com/jackc/pgx/v5"
	tern "github.com/jackc/tern/v2/migrate"
	"github.com/mustaphalimar/gopanel/pkg/config"
	"github.com/mustaphalimar/gopanel/pkg/logger"
)

// go:embed migrations/*sql
var migrations embed.FS

func Migrate(ctx context.Context, cfg *config.Config, logger *logger.Logger) error {
	conn, err := pgx.Connect(ctx, cfg.Database.GetDatabaseURI())
	if err != nil {
		return nil
	}
	defer conn.Close(ctx)

	m, err := tern.NewMigrator(ctx, conn, "schema_version")
	if err != nil {
		return fmt.Errorf("error constructing database migrator: %w", err)
	}

	subTree, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("error retrieving database migrations subtree: %w", err)
	}

	if err := m.LoadMigrations(subTree); err != nil {
		return fmt.Errorf("loading database migrations: %w", err)
	}

	from, err := m.GetCurrentVersion(ctx)
	if err != nil {
		return fmt.Errorf("error retreiving current database migration version")
	}

	if err := m.Migrate(ctx); err != nil {
		return err
	}

	if from == int32(len(m.Migrations)) {
		logger.InfoContext(ctx, "database schema up to date", "version", len(m.Migrations))
	} else {
		logger.InfoContext(ctx, "migrated database schema", "from", from, "to", len(m.Migrations))
	}
	return nil
}
