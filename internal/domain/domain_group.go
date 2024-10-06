package domain

import (
	"context"
	"emobile/internal/errors"
	"emobile/internal/models"
)

func (d *domain) NewGroup(data models.NewGroupReq) (string, errors.APIError) {

	groupID, err := d.pg.NewGroup(context.Background(), data)

	if err != nil {
		d.pg.Log.Error(err.Error())
		return "", errors.NewHTTPError(500, err.Error())
	}

	if groupID == "" {
		return "", errors.NewHTTPError(500, "Failed to create group")
	}

	return groupID, nil

}

func (d *domain) GetAllGroups() ([]models.Group, errors.APIError) {
	groups, err := d.pg.GetAllGroups()
	if err != nil {
		return nil, errors.NewHTTPError(500, err.Error())
	}
	if len(groups) == 0 {
		return nil, errors.NewHTTPError(404, "No groups found")
	}
	return groups, nil
}

func (d *domain) GetGroupSongs(group string) ([]models.Song, errors.APIError) {

	songs, err := d.pg.GetGroupSongs(group)
	if err != nil {
		return nil, errors.NewHTTPError(500, err.Error())
	}

	if len(songs) == 0 {
		return nil, errors.NewHTTPError(404, "No songs found")
	}
	return songs, nil
}

func (d *domain) EditGroup(group_name models.Group) errors.APIError {
	groupID, err := d.pg.EditGroup(group_name)

	if err != nil {
		return errors.NewHTTPError(500, err.Error())
	}

	if groupID == "" {
		return errors.NewHTTPError(404, "Not found")
	}

	return nil

}
