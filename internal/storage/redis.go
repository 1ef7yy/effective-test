package storage

import (
	"context"
	"emobile/internal/models"
	"emobile/pkg/logger"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Log logger.Logger
	*redis.Client
}

func NewRedis(log logger.Logger) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_CONN"),
		Password: "",
		DB:       0,
	})

	return &Redis{
		Log:    log,
		Client: rdb,
	}, nil
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

func (r *Redis) LoadDataFromPG(ctx context.Context) error {
	return nil
}

func (r *Redis) GetHash(ctx context.Context, key string) (map[string]string, error) {
	return r.Client.HGetAll(ctx, key).Result()
}

func (r *Redis) SetHash(ctx context.Context, key string, value map[string]string) error {
	return r.Client.HSet(ctx, key, "song_id", value).Err()
}

func (r *Redis) GetSong(group, song string) (models.RedisSongDTO, error) {

	val, err := r.GetHash(context.Background(), song)

	// redis empty

	if err == redis.Nil {
		return models.RedisSongDTO{}, nil
	}

	if err != nil {
		r.Log.Error(err.Error())
		return models.RedisSongDTO{}, err
	}

	fmt.Printf("Value redis: %v\n", val)
	ReleaseDateTime, err := time.Parse("2006-01-02", val["release_date"])

	if err != nil {
		r.Log.Error(err.Error())
		return models.RedisSongDTO{}, err
	}

	return models.RedisSongDTO{
		SongID:      val["song_id"],
		GroupID:     val["group_id"],
		SongName:    val["song_name"],
		ReleaseDate: ReleaseDateTime,
		SongText:    val["song_text"],
		Link:        val["link"],
	}, nil

}

func (r *Redis) NewSong(data models.NewSongFormattedReq) error {
	var Song models.RedisSongDTO
	Song.SongName = data.SongName
	MapReleaseDate := data.ReleaseDate.UTC().String()
	Song.SongText = data.SongText
	Song.Link = data.Link

	err := r.SetHash(context.Background(), data.SongName, map[string]string{
		"song_id":      Song.SongID,
		"group_id":     Song.GroupID,
		"song_name":    Song.SongName,
		"release_date": MapReleaseDate,
		"song_text":    Song.SongText,
		"link":         Song.Link,
	})

	if err != nil {
		r.Log.Error(err.Error())
		return err
	}

	return nil
}
