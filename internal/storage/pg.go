package storage

import (
	"context"
	"emobile/internal/errors"
	"emobile/internal/models"
	"emobile/pkg/logger"
	"sync"
	"time"

	"github.com/jackc/pgx"
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

func (pg *Postgres) GetAllSongs() ([]models.Song, error) {

	val, err := pg.DB.Query(context.Background(), "SELECT * FROM songs")

	if err != nil {
		pg.Log.Error(err.Error())
		return nil, err
	}

	if val.RawValues() == nil {
		return nil, errors.NewHTTPError(404, "Not found")
	}

	defer val.Close()

	var songs []models.Song

	for val.Next() {
		var song models.Song

		if err := val.Scan(&song.SongID, &song.GroupID, &song.ReleaseDate, &song.SongName, &song.SongText, &song.Link); err != nil {
			pg.Log.Error(err.Error())
			return nil, err
		}

		songs = append(songs, song)
	}

	return songs, nil
}

func (pg *Postgres) GetSong(group_name, song string) (models.SongDB, error) {

	group_id, err := pg.GetGroupID(group_name)

	if err != nil {
		pg.Log.Error(err.Error())
		return models.SongDB{}, err
	}

	var Song models.SongDB
	var release_date string

	err = pg.DB.QueryRow(context.Background(), "SELECT song_id, group_id, release_date::text, song_name, song_text, link FROM songs WHERE group_id = $1 AND song_name = $2", group_id, song).Scan(&Song.SongID, &Song.GroupID, &release_date, &Song.SongName, &Song.SongText, &Song.Link)

	if err != nil {
		pg.Log.Error(err.Error())
		return models.SongDB{}, err
	}

	if Song.SongID == "" {
		return models.SongDB{}, errors.NewHTTPError(404, "Not found")
	}

	Song.ReleaseDate, err = time.Parse("2006-01-02", release_date)

	if err != nil {
		pg.Log.Error(err.Error())
		return models.SongDB{}, err
	}

	return Song, nil
}
func (pg *Postgres) GetGroupID(groupName string) (string, error) {

	if groupName == "" {
		return "", errors.NewHTTPError(400, "Bad request, missing group name")
	}

	var groupID string
	err := pg.DB.QueryRow(context.Background(), "SELECT group_id FROM groups WHERE group_name = $1", groupName).Scan(&groupID)

	if err != nil {
		pg.Log.Error(err.Error())
		if err == pgx.ErrNoRows {
			return "", errors.NewHTTPError(404, "Not found")
		}
		return "", err
	}

	if groupID == "" {
		return "", errors.NewHTTPError(404, "Not found")
	}

	return groupID, nil
}

func (pg *Postgres) NewSong(data models.NewSongReq) (string, error) {

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

	if songID == "" {
		return "", errors.NewHTTPError(500, "Internal server error")
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
