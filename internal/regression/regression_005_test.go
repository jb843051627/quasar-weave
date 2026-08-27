package service_test

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"github.com/jb843051627/quasar-weave/internal/store"
)



func TestBug05_FailedTransactionLeavesNoRecord(t *testing.T) {
    dir := t.TempDir()
    repo, err := store.Open(filepath.Join(dir, "q.db"))
    if err != nil { t.Fatal(err) }
    defer repo.Close()
    sentinel := errors.New("abort transaction")
    err = repo.Transaction(context.Background(), func(tx *sql.Tx) error {
        if _, execErr := tx.Exec(`INSERT INTO records(kind,id,payload,updated_at) VALUES('probe','one','{}','now')`); execErr != nil { return execErr }
        return sentinel
    })
    if !errors.Is(err, sentinel) { t.Fatalf("error=%v", err) }
    count, err := repo.Count(context.Background(), "probe")
    if err != nil { t.Fatal(err) }
    if count != 0 { t.Fatalf("rolled-back records=%d", count) }
}
