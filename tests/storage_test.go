package storage_test

import (
	"context"
	"emobile/internal/storage"
	"emobile/pkg/logger"
	"os"
	"testing"
)

func TestPostgresPing(t *testing.T) {
	log := logger.NewLogger(nil)
	pg := storage.NewPostgres(context.Background(), os.Getenv("POSTGRES_DSN"), log)

	if err := pg.Ping(context.Background()); err != nil {
		t.Fatal(err)
	}
}
