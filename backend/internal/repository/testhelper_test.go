package repository_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// newTestDB spins up a disposable Postgres container and returns a ready *sqlx.DB.
func newTestDB(t *testing.T) *sqlx.DB {
	t.Helper()
	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithDatabase(uuid.NewString()),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(wait.ForListeningPort("5432/tcp")),
	)
	require.NoError(t, err)

	t.Cleanup(func() { _ = container.Terminate(ctx) })

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := sqlx.Connect("postgres", dsn)
	require.NoError(t, err)

	migrate(t, db)
	return db
}

// migrate applies the minimal schema required by the repository tests.
func migrate(t *testing.T, db *sqlx.DB) {
	t.Helper()
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id            SERIAL       PRIMARY KEY,
			username      TEXT         UNIQUE NOT NULL,
			password_hash TEXT         NOT NULL,
			created_at    TIMESTAMPTZ  DEFAULT NOW()
		);
		CREATE TABLE IF NOT EXISTS categories (
			id      SERIAL       PRIMARY KEY,
			user_id INTEGER      REFERENCES users(id) ON DELETE CASCADE,
			name    VARCHAR(50)  NOT NULL
		);
		CREATE TABLE IF NOT EXISTS tasks (
			id           SERIAL       PRIMARY KEY,
			user_id      INT          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			task_text    TEXT         NOT NULL,
			priority     SMALLINT     NOT NULL DEFAULT 2 CHECK (priority BETWEEN 1 AND 3),
			deadline     DATE,
			is_completed BOOLEAN      NOT NULL DEFAULT FALSE,
			position     INT          NOT NULL DEFAULT 0,
			category_id  INTEGER      REFERENCES categories(id) ON DELETE SET NULL,
			created_at   TIMESTAMPTZ  DEFAULT NOW()
		);
		CREATE TABLE IF NOT EXISTS task_collaborators (
			task_id INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			PRIMARY KEY (task_id, user_id)
		);
	`)
	require.NoError(t, err)
}
