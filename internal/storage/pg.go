package storage

import (
	"context"
	"emobile/internal/models"
	"emobile/pkg/logger"
	"fmt"
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

	val, err := pg.DB.Query(context.Background(), "SELECT song_id, group_id, release_date::text, song_name, song_text, link FROM songs")

	fmt.Println("1")

	if err != nil {
		pg.Log.Error(err.Error())
		return nil, err
	}

	defer val.Close()

	var songs []models.Song

	for val.Next() {
		var song models.Song

		var release_date string

		if err := val.Scan(&song.SongID, &song.GroupID, &release_date, &song.SongName, &song.SongText, &song.Link); err != nil {
			pg.Log.Error(err.Error())
			return nil, err
		}

		song.GroupName, err = pg.GetGroupName(song.GroupID)

		if song.GroupName == "" {
			return []models.Song{}, nil
		}

		if err != nil {
			pg.Log.Error(err.Error())
			return nil, err
		}

		song.ReleaseDate, err = time.Parse("2006-01-02", release_date)

		if err != nil {
			pg.Log.Error(err.Error())
			return nil, err
		}

		songs = append(songs, song)
	}

	return songs, nil
}

func (pg *Postgres) GetSong(group_name, song string) (models.SongDB, error) {

	group_id, err := pg.GetGroupID(group_name)

	if group_id == "" {
		return models.SongDB{
			SongID: "",
		}, nil
	}

	if err != nil {
		pg.Log.Error(err.Error())
		return models.SongDB{}, err
	}

	var Song models.SongDB
	var release_date string

	err = pg.DB.QueryRow(context.Background(), "SELECT song_id, group_id, release_date::text, song_name, song_text, link FROM songs WHERE group_id = $1 AND song_name = $2", group_id, song).Scan(&Song.SongID, &Song.GroupID, &release_date, &Song.SongName, &Song.SongText, &Song.Link)

	if err == pgx.ErrNoRows {
		return models.SongDB{
			SongID: "",
		}, nil
	}

	if err != nil {
		pg.Log.Error(err.Error())
		return models.SongDB{}, err
	}

	if Song.SongID == "" {
		return models.SongDB{
			SongID: "",
		}, nil
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
		return "", nil
	}

	var groupID string
	err := pg.DB.QueryRow(context.Background(), "SELECT group_id FROM groups WHERE group_name = $1", groupName).Scan(&groupID)

	if err != nil {
		pg.Log.Error(err.Error())
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return groupID, nil
}

func (pg *Postgres) GetGroupName(groupID string) (string, error) {

	var groupName string

	err := pg.DB.QueryRow(context.Background(), "SELECT group_name FROM groups WHERE group_id = $1", groupID).Scan(&groupName)

	if err != nil {
		pg.Log.Error(err.Error())
		return "", err
	}

	return groupName, nil
}

func (pg *Postgres) NewSong(data models.NewSongReq) (string, error) {

	var group_id string

	err := pg.DB.QueryRow(context.Background(), "SELECT group_id FROM groups WHERE group_name = $1", data.GroupName).Scan(&group_id)

	if err == pgx.ErrNoRows {
		pg.Log.Error(err.Error())
		return "", fmt.Errorf("group %s not found", data.GroupName)
	}
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
		return "", fmt.Errorf("song_id is empty")
	}
	return songID, nil

}

func (pg *Postgres) EditSong(data models.EditSongReq) (string, error) {
	song_id, err := pg.GetSong(data.GroupName, data.SongName)

	if err != nil {
		pg.Log.Error(err.Error())
		return "", err
	}

	if song_id.SongID == "" {
		return "", nil
	}

	_, err = pg.DB.Exec(context.Background(), "UPDATE songs SET song_name = $1, release_date = $2, song_text = $3, link = $4 WHERE song_id = $5",
		data.SongName, data.ReleaseDate, data.SongText, data.Link, song_id.SongID)

	if err != nil {
		pg.Log.Error(err.Error())
		return "", err
	}

	return song_id.SongID, nil
}

func (pg *Postgres) DeleteSong(data models.DeleteSongReq) (string, error) {
	song_id, err := pg.GetSong(data.GroupName, data.SongName)

	if err != nil {
		pg.Log.Error(err.Error())
		return "", err
	}

	if song_id.SongID == "" {
		return "", nil
	}

	_, err = pg.DB.Exec(context.Background(), "DELETE FROM songs WHERE song_id = $1", song_id.SongID)

	if err != nil {
		pg.Log.Error(err.Error())
		return "", err
	}

	return song_id.SongID, nil

}

func (pg *Postgres) NewGroup(ctx context.Context, data models.NewGroupReq) (string, error) {

	val, err := pg.DB.Query(context.Background(), "SELECT group_id FROM groups WHERE group_name = $1", data.GroupName)

	if err != nil {
		pg.Log.Error(err.Error())
		return "", err
	}

	defer val.Close()

	var groupID string

	if val.Next() {

		if err := val.Scan(&groupID); err != nil {
			pg.Log.Error(err.Error())
			return "", err
		}
	}

	if groupID != "" {
		return "", fmt.Errorf("group %s already exists", data.GroupName)
	}

	val, err = pg.DB.Query(context.Background(), "INSERT INTO groups (group_name) VALUES ($1) RETURNING group_id",
		data.GroupName)

	if err != nil {
		pg.Log.Error(err.Error())
		return "", err
	}

	defer val.Close()

	for val.Next() {

		if err := val.Scan(&groupID); err != nil {
			pg.Log.Error(err.Error())
			return "", err

		}
	}

	if groupID == "" {
		return "", fmt.Errorf("group_id is empty")
	}

	return groupID, nil
}

func (pg *Postgres) GetAllGroups() ([]models.Group, error) {
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

func (pg *Postgres) GetGroupSongs(group_name string) ([]models.Song, error) {
	groupID, err := pg.GetGroupID(group_name)

	val, err := pg.DB.Query(context.Background(), "SELECT song_id, group_id, release_date::text, song_name, song_text, link FROM songs WHERE group_id = $1", groupID)

	if err != nil {
		pg.Log.Error(err.Error())
		return nil, err
	}

	var songs []models.Song

	for val.Next() {

		var song models.Song
		var release_date string

		if err := val.Scan(&song.SongID, &song.GroupID, &release_date, &song.SongName, &song.SongText, &song.Link); err != nil {
			pg.Log.Error(err.Error())
			return nil, err
		}

		song.ReleaseDate, err = time.Parse("2006-01-02", release_date)

		if err != nil {
			pg.Log.Error(err.Error())
			return nil, err
		}

		songs = append(songs, song)
	}

	return songs, nil
}

func (pg *Postgres) EditGroup(data models.Group) (string, error) {

	var groupID string

	err := pg.DB.QueryRow(context.Background(), "UPDATE groups SET group_name = $1 WHERE group_id = $2 RETURNING group_id", data.GroupName, data.GroupID).Scan(&groupID)

	if err != nil {
		pg.Log.Error(err.Error())
		return "", err
	}

	return groupID, nil
}
