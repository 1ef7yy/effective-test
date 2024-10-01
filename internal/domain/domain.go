package domain

import (
	"context"
	"emobile/internal/storage"
	"emobile/pkg/logger"
	"os"
)

type Domain interface {
}

type domain struct {
	pg storage.Postgres
}

func NewDomain(log logger.Logger) Domain {
	return &domain{
		pg: *storage.NewPostgres(context.Background(), os.Getenv("POSTGRES_CONN"), log),
	}
}
