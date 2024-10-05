package storage

import (
	"context"
	"emobile/internal/models"
	"emobile/pkg/logger"
	"fmt"
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

func (pg *Postgres) GetSong(group_name, song string) (models.Song, error) {

	group_id, err := pg.GetGroupID(group_name)

	fmt.Printf("Group id: %s\n", group_id)

	if err != nil {
		pg.Log.Error(err.Error())
		return models.Song{}, err
	}

	val, err := pg.DB.Query(context.Background(), "SELECT * FROM songs WHERE group_id = $1 AND song_name = $2", group_id, song)

	if err != nil {
		pg.Log.Error(err.Error())
		return models.Song{}, err
	}

	var Song models.Song

	for val.Next() {

		if err := val.Scan(&Song.SongID, &Song.GroupID, &Song.ReleaseDate, &Song.SongName, &Song.SongText, &Song.Link); err != nil {
			pg.Log.Error(err.Error())
			return models.Song{}, err
		}
	}

	return Song, nil
}

func (pg *Postgres) GetGroupID(groupName string) (string, error) {

	var groupID string
	err := pg.DB.QueryRow(context.Background(), "SELECT group_id FROM groups WHERE group_name = $1", groupName).Scan(&groupID)

	if err != nil {
		pg.Log.Error(err.Error())
		return "", err
	}

	return groupID, nil

}

func (pg *Postgres) NewSong(data models.NewSongReq) (string, error) {

	fmt.Printf("Data pg: %v\n", data)

	var group_id string

	err := pg.DB.QueryRow(context.Background(), "SELECT group_id FROM groups WHERE group_name = $1", data.GroupName).Scan(&group_id)
	if err != nil {
		pg.Log.Error(err.Error())
		return "", err
	}

	var songID string
	err = pg.DB.QueryRow(context.Background(), "INSERT INTO songs (group_id, song_name, release_date, song_text, link) VALUES ($1, $2, $3, $4, $5) RETURNING song_id",
		group_id, data.SongName, data.ReleaseDate, data.SongText, data.Link).Scan(&songID)

	if err != nil {
		pg.Log.Error(err.Error())
		return "", err
	}

	if err != nil {
		pg.Log.Error(err.Error())
		return "", err
	}

	return songID, nil

}

func (pg *Postgres) NewGroup(ctx context.Context, data models.NewGroupReq) (string, error) {

	val, err := pg.DB.Query(context.Background(), "INSERT INTO groups (group_name) VALUES ($1) RETURNING group_id",
		data.GroupName)

	if err != nil {
		pg.Log.Error(err.Error())
		return "", err
	}

	var groupID string
	defer val.Close()

	for val.Next() {

		if err := val.Scan(&groupID); err != nil {
			pg.Log.Error(err.Error())
			return "", err

		}
	}

	return groupID, nil
}

func (pg *Postgres) GetGroups() ([]models.Group, error) {
	val, err := pg.DB.Query(context.Background(), "SELECT * FROM groups")

	if err != nil {
		pg.Log.Error(err.Error())
		return nil, err
	}

	var groups []models.Group

	for val.Next() {
		var group models.Group

		if err := val.Scan(&group.GroupID, &group.GroupName); err != nil {
			pg.Log.Error(err.Error())
			return nil, err
		}

		groups = append(groups, group)
	}

	return groups, nil
}
