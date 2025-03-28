package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAuditEvents, downAuditEvents)
}

func upAuditEvents(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
	CREATE TABLE events (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		payload jsonb NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`)
	return err
}

func downAuditEvents(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
	DROP TABLE IF EXISTS events;
	`)
	return err
}
