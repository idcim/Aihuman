package video

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(ctx context.Context, databaseURL string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := waitForDB(ctx, db, 30*time.Second); err != nil {
		_ = db.Close()
		return nil, err
	}

	store := &PostgresStore{db: db}
	if err := store.runMigrations(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

func (s *PostgresStore) Save(ctx context.Context, task Task) error {
	renderCommand, err := json.Marshal(task.RenderCommand)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO video_tasks (
			id, product_name, script, avatar_image_url, duration_seconds, status, render_command, output_url, object_key, preview_seconds, error_message, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (id) DO UPDATE SET
			product_name = EXCLUDED.product_name,
			script = EXCLUDED.script,
			avatar_image_url = EXCLUDED.avatar_image_url,
			duration_seconds = EXCLUDED.duration_seconds,
			status = EXCLUDED.status,
			render_command = EXCLUDED.render_command,
			output_url = EXCLUDED.output_url,
			object_key = EXCLUDED.object_key,
			preview_seconds = EXCLUDED.preview_seconds,
			error_message = EXCLUDED.error_message,
			updated_at = EXCLUDED.updated_at
	`, task.ID, task.ProductName, task.Script, task.AvatarImageURL, task.DurationSeconds, string(task.Status), string(renderCommand), task.OutputURL, task.ObjectKey, task.PreviewSeconds, task.ErrorMessage, task.CreatedAt, task.UpdatedAt)
	return err
}

func (s *PostgresStore) Get(ctx context.Context, id string) (Task, error) {
	task, err := scanTask(s.db.QueryRowContext(ctx, `
		SELECT id, product_name, script, avatar_image_url, duration_seconds, status, render_command, output_url, object_key, preview_seconds, error_message, created_at, updated_at
		FROM video_tasks
		WHERE id = $1
	`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return Task{}, ErrTaskNotFound
	}
	return task, err
}

func (s *PostgresStore) List(ctx context.Context) ([]Task, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, product_name, script, avatar_image_url, duration_seconds, status, render_command, output_url, object_key, preview_seconds, error_message, created_at, updated_at
		FROM video_tasks
		ORDER BY created_at DESC
		LIMIT 100
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *PostgresStore) UpdateStatus(ctx context.Context, id string, status TaskStatus, result RenderResult, errorMessage string) (Task, error) {
	renderCommandJSON, err := json.Marshal(result.Command)
	if err != nil {
		return Task{}, err
	}

	task, err := scanTask(s.db.QueryRowContext(ctx, `
		UPDATE video_tasks
		SET status = $2,
			render_command = $3::jsonb,
			output_url = $4,
			object_key = $5,
			preview_seconds = $6,
			error_message = $7,
			updated_at = $8
		WHERE id = $1
		RETURNING id, product_name, script, avatar_image_url, duration_seconds, status, render_command, output_url, object_key, preview_seconds, error_message, created_at, updated_at
	`, id, string(status), string(renderCommandJSON), result.OutputURL, result.ObjectKey, result.PreviewSeconds, errorMessage, time.Now().UTC()))
	if errors.Is(err, sql.ErrNoRows) {
		return Task{}, ErrTaskNotFound
	}
	return task, err
}

func (s *PostgresStore) Close() error {
	return s.db.Close()
}

func (s *PostgresStore) runMigrations(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`); err != nil {
		return err
	}

	entries, err := fs.ReadDir(migrationFS, "migrations")
	if err != nil {
		return err
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		if err := s.applyMigration(ctx, entry.Name()); err != nil {
			return err
		}
	}
	return nil
}

func (s *PostgresStore) applyMigration(ctx context.Context, version string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var exists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", version).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return tx.Commit()
	}

	sqlBytes, err := migrationFS.ReadFile("migrations/" + version)
	if err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, string(sqlBytes)); err != nil {
		return fmt.Errorf("apply migration %s: %w", version, err)
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
		return err
	}
	return tx.Commit()
}

type taskScanner interface {
	Scan(dest ...interface{}) error
}

func scanTask(scanner taskScanner) (Task, error) {
	var task Task
	var status string
	var renderCommand []byte
	err := scanner.Scan(
		&task.ID,
		&task.ProductName,
		&task.Script,
		&task.AvatarImageURL,
		&task.DurationSeconds,
		&status,
		&renderCommand,
		&task.OutputURL,
		&task.ObjectKey,
		&task.PreviewSeconds,
		&task.ErrorMessage,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	task.Status = TaskStatus(status)
	if len(renderCommand) > 0 {
		if jsonErr := json.Unmarshal(renderCommand, &task.RenderCommand); jsonErr != nil && err == nil {
			err = jsonErr
		}
	}
	return task, err
}

func waitForDB(ctx context.Context, db *sql.DB, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		if err := db.PingContext(ctx); err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
