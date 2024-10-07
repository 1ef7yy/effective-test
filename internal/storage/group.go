package storage

import (
	"context"
	"emobile/internal/models"
	"fmt"
	"time"

	"github.com/jackc/pgx"
)

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

func (pg *Postgres) GetGroupName(groupID string) (string, error) {

	var groupName string

	err := pg.DB.QueryRow(context.Background(), "SELECT group_name FROM groups WHERE group_id = $1", groupID).Scan(&groupName)

	if err != nil {
		pg.Log.Error(err.Error())
		return "", err
	}

	return groupName, nil
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

func (pg *Postgres) DeleteGroup(groupID string) error {
	return pg.DB.QueryRow(context.Background(), "DELETE FROM groups WHERE group_id = $1", groupID).Scan()
}
