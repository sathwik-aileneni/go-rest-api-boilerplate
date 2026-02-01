package riverenqueuer

import (
	"context"
	"database/sql"

	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivermigrate"
)

// MigrateUp runs River schema migrations.
// This is idempotent and safe to call on every application startup.
// It creates the river_job, river_leader, and river_migration tables if they don't exist.
func MigrateUp(ctx context.Context, db *sql.DB) error {
	migrator, err := rivermigrate.New(riverdatabasesql.New(db), nil)
	if err != nil {
		return err
	}

	_, err = migrator.Migrate(ctx, rivermigrate.DirectionUp, nil)
	return err
}
