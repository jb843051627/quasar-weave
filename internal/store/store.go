package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/jb843051627/quasar-weave/internal/model"
	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
	mu sync.RWMutex
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf("create database directory: %w", err)
		}
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	s := &Store{db: db}
	if err := s.schema(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) schema(ctx context.Context) error {
	const ddl = `
CREATE TABLE IF NOT EXISTS records (
  kind TEXT NOT NULL,
  id TEXT NOT NULL,
  payload BLOB NOT NULL,
  updated_at TEXT NOT NULL,
  PRIMARY KEY(kind, id)
);
CREATE INDEX IF NOT EXISTS records_kind_updated ON records(kind, updated_at);
CREATE TABLE IF NOT EXISTS events (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  subject TEXT NOT NULL,
  action TEXT NOT NULL,
  payload TEXT NOT NULL,
  created_at TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS events_subject_created ON events(subject, created_at);
`
	if _, err := s.db.ExecContext(ctx, ddl); err != nil {
		return fmt.Errorf("create schema: %w", err)
	}
	return nil
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) DB() *sql.DB {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db
}

func (s *Store) Save(ctx context.Context, kind, id string, value any) error {
	if err := model.ValidateID(kind, id); err != nil {
		return err
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal %s/%s: %w", kind, id, err)
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	_, err = s.db.ExecContext(ctx, `
INSERT INTO records(kind, id, payload, updated_at) VALUES(?, ?, ?, ?)
ON CONFLICT(kind, id) DO UPDATE SET payload=excluded.payload, updated_at=excluded.updated_at`, kind, id, raw, now)
	if err != nil {
		return fmt.Errorf("save %s/%s: %w", kind, id, err)
	}
	return nil
}

func (s *Store) Load(ctx context.Context, kind, id string, into any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	var raw []byte
	err := s.db.QueryRowContext(ctx, `SELECT payload FROM records WHERE kind=? AND id=?`, kind, id).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return model.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("load %s/%s: %w", kind, id, err)
	}
	if err := json.Unmarshal(raw, into); err != nil {
		return fmt.Errorf("decode %s/%s: %w", kind, id, err)
	}
	return nil
}

func (s *Store) Delete(ctx context.Context, kind, id string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	result, err := s.db.ExecContext(ctx, `DELETE FROM records WHERE kind=? AND id=?`, kind, id)
	if err != nil {
		return fmt.Errorf("delete %s/%s: %w", kind, id, err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (s *Store) List(ctx context.Context, kind string, into func([]byte) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	rows, err := s.db.QueryContext(ctx, `SELECT payload FROM records WHERE kind=? ORDER BY updated_at, id`, kind)
	if err != nil {
		return fmt.Errorf("list %s: %w", kind, err)
	}
	defer rows.Close()
	for rows.Next() {
		var raw []byte
		if err := rows.Scan(&raw); err != nil {
			return fmt.Errorf("scan %s: %w", kind, err)
		}
		if err := into(raw); err != nil {
			return err
		}
	}
	return rows.Err()
}

func (s *Store) Count(ctx context.Context, kind string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return 0, errors.New("store is closed")
	}
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM records WHERE kind=?`, kind).Scan(&count); err != nil {
		return 0, fmt.Errorf("count %s: %w", kind, err)
	}
	return count, nil
}

func (s *Store) Event(ctx context.Context, subject, action, payload string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO events(subject, action, payload, created_at) VALUES(?, ?, ?, ?)`, subject, action, payload, time.Now().UTC().Format(time.RFC3339Nano))
	if err != nil {
		return fmt.Errorf("append event: %w", err)
	}
	return nil
}

func (s *Store) Events(ctx context.Context, subject string, limit int) ([]model.Event, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, errors.New("store is closed")
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id, subject, action, payload, created_at FROM events WHERE (?='' OR subject=?) ORDER BY id DESC LIMIT ?`, subject, subject, limit)
	if err != nil {
		return nil, fmt.Errorf("list events: %w", err)
	}
	defer rows.Close()
	result := make([]model.Event, 0, limit)
	for rows.Next() {
		var sequence int64
		var created string
		var event model.Event
		if err := rows.Scan(&sequence, &event.Subject, &event.Action, &event.Payload, &created); err != nil {
			return nil, err
		}
		event.ID = fmt.Sprint(sequence)
		event.CreatedAt, err = time.Parse(time.RFC3339Nano, created)
		if err != nil {
			return nil, fmt.Errorf("parse event time: %w", err)
		}
		result = append(result, event)
	}
	return result, rows.Err()
}

func (s *Store) Transaction(ctx context.Context, fn func(*sql.Tx) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("transaction: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func SaveJSON[T any](ctx context.Context, s *Store, kind, id string, value T) error {
	return s.Save(ctx, kind, id, value)
}

func LoadJSON[T any](ctx context.Context, s *Store, kind, id string) (T, error) {
	var value T
	if err := s.Load(ctx, kind, id, &value); err != nil {
		return value, err
	}
	return value, nil
}

func ListJSON[T any](ctx context.Context, s *Store, kind string) ([]T, error) {
	items := make([]T, 0)
	err := s.List(ctx, kind, func(raw []byte) error {
		var item T
		if err := json.Unmarshal(raw, &item); err != nil {
			return fmt.Errorf("decode %s: %w", kind, err)
		}
		items = append(items, item)
		return nil
	})
	return items, err
}
