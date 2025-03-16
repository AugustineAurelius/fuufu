package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upDdlForUser, downDdlForUser)
}

func upDdlForUser(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
	CREATE TABLE users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR NOT NULL,
		email VARCHAR NOT NULL,
		hashed_password VARCHAR NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ
	);

	CREATE OR REPLACE FUNCTION update_updated_at_column()
	RETURNS TRIGGER AS $$
	BEGIN
		NEW.updated_at = NOW();
		RETURN NEW;
	END;
	$$ LANGUAGE plpgsql;

	CREATE TRIGGER update_task_updated_at
	BEFORE UPDATE ON users
	FOR EACH ROW
	EXECUTE FUNCTION update_updated_at_column();
`)
	return err
}

func downDdlForUser(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
	DROP TRIGGER IF EXISTS update_task_updated_at ON users;
	DROP TABLE IF EXISTS users;
	`)
	return err
}
