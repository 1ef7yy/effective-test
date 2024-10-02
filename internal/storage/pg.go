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

func (pg *Postgres) Info(ctx context.Context, song, group string) (models.InfoDTOResponse, error) {

	val, err := pg.DB.Query(ctx, "GET song_details FROM songs WHERE song_name = $1 AND group_name = $2", song, group)

	if err != nil {
		return models.InfoDTOResponse{}, err
	}

	defer val.Close()

	var info models.InfoDTOResponse
	for val.Next() {
		if err := val.Scan(&info); err != nil {
			return models.InfoDTOResponse{}, err
		}
	}

	return info, nil

}

func (pg *Postgres) GetGroup(ctx context.Context, group string) (models.Group, error) {

	val, err := pg.DB.Query(ctx, "SELECT group_id, group_name FROM groups WHERE group_name = $1 ORDER BY group_name", group)

	if err != nil {
		return models.Group{}, err
	}

	defer val.Close()

	var info models.Group

	for val.Next() {
		if err := val.Scan(&info.GroupID, &info.GroupName); err != nil {
			return models.Group{}, err
		}
	}

	return info, nil
}

func (pg *Postgres) GetGroupSongs(ctx context.Context, group string, limit, offset int) ([]models.Song, error) {

	val, err := pg.DB.Query(ctx, "SELECT song_name, release_date, song_text, link FROM songs WHERE group_name = $1 LIMIT $2 OFFSET $3 ORDER BY song_name, release_date", group, limit, offset)

	if err != nil {
		return []models.Song{}, err
	}

	defer val.Close()

	var info []models.Song
	for val.Next() {
		if err := val.Scan(&info); err != nil {
			return []models.Song{}, err
		}
	}

	return info, nil
}
