package storage

import (
	"context"
	"emobile/pkg/logger"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Log logger.Logger
	DB  *pgxpool.Pool
}

func NewPostgres(ctx context.Context, dsn string, log logger.Logger) *Postgres {
	var (
		pgInstance *Postgres
		pgOnce     sync.Once
	)

	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, dsn)
		if err != nil {
			log.Fatal("Unable to connect to database: " + err.Error())
		}

		pgInstance = &Postgres{
			Log: log,
			DB:  db,
		}
	})
	return pgInstance
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.DB.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.DB.Close()
}
