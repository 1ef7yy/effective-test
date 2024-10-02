package storage

import (
	"context"
	"emobile/internal/models"
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

func (pg *Postgres) GetSong(group, song string) (models.Song, error) {

	val, err := pg.DB.Query(context.Background(), "SELECT * FROM songs WHERE group_id = $1 AND song_name = $2", group, song)

	if err != nil {
		return models.Song{}, err
	}

	var Song models.Song

	for val.Next() {

		if err := val.Scan(&Song.SongID, &Song.GroupID, &Song.ReleaseDate, &Song.SongName, &Song.SongText, &Song.Link); err != nil {
			return models.Song{}, err
		}
	}

	return Song, nil
}
