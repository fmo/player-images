package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fmo/player-images/internal/application/core/domain"
	"github.com/fmo/player-images/internal/ports"
	log "github.com/sirupsen/logrus"
	"time"
)

type Application struct {
	image      ports.ImagePort
	playerData ports.PlayerData
	cache      ports.CachePort
}

func NewApplication(
	image ports.ImagePort,
	playerData ports.PlayerData,
	cache ports.CachePort,
) *Application {
	return &Application{
		image:      image,
		playerData: playerData,
		cache:      cache,
	}
}

func (a Application) SavePlayerImage(ctx context.Context, season, teamId int) error {
	cacheKey := fmt.Sprintf("team:%d:season:%d", season, teamId)
	cachedPlayers, err := a.cache.Get(ctx, cacheKey)
	if err != nil {
		return err
	}

	var players []domain.Player

	if cachedPlayers != "" {
		if err := json.Unmarshal([]byte(cachedPlayers), &players); err != nil {
			log.Fatalf("Error unmarshalling cached json: %v\n", err)
		}
	} else {
		players = a.playerData.GetPlayers(season, teamId)
		playersJson, err := json.Marshal(players)
		if err != nil {
			log.Fatalf("Error marshalling json: %v\n", err)
		}
		a.cache.Set(ctx, cacheKey, playersJson, 24*30*time.Hour)
	}

	for _, player := range players {
		imageName := fmt.Sprintf("%s.png", player.ID)

		if a.image.CheckImageAlreadyUploaded(imageName) {
			log.Infof("Image %s is already uploaded", imageName)
		} else {
			log.Infof("Image %s is not yet uploaded, doing so", imageName)
			err := a.image.Upload(imageName, player.Image)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
