package storage_test

import (
	"context"
	"emobile/internal/storage"
	"emobile/pkg/logger"
	"testing"
)

func TestPostgresPing(t *testing.T) {
	log := logger.NewLogger(nil)
	conn := "postgres://postgres:postgres@127.0.0.1:5432/library?sslmode=disable"
	pg := storage.NewPostgres(context.Background(), conn, log)
	if err := pg.Ping(context.Background()); err != nil {
		t.Fatal(err)
	}
}
