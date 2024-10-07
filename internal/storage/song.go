package storage

import (
	"context"
	"emobile/internal/models"
	"fmt"
	"time"

	"github.com/jackc/pgx"
)

func (pg *Postgres) GetAllSongs() ([]models.Song, error) {

	val, err := pg.DB.Query(context.Background(), "SELECT song_id, group_id, release_date::text, song_name, song_text, link FROM songs")

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
