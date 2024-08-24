package transfermarkt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Player struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Data struct {
	Players []Player `json:"data"`
}

type PlayersApi struct {
	logger      *logrus.Logger
	redisClient *redis.Client
}

func NewPlayersApi(l *logrus.Logger, rdb *redis.Client) PlayersApi {
	return PlayersApi{
		logger:      l,
		redisClient: rdb,
	}
}

func (p PlayersApi) GetPlayers(season, teamId int) []Player {
	cacheKey := fmt.Sprintf("team:%d:season:%d", teamId, season)
	cachedPlayers, err := p.redisClient.Get(context.Background(), cacheKey).Result()
	if errors.Is(err, redis.Nil) {
		p.logger.Info("Cache miss, fetching from API")
		response := request(
			fmt.Sprintf(
				"https://transfermarkt-db.p.rapidapi.com/v1/clubs/squad?season_id=%d&locale=UK&club_id=%d",
				season,
				teamId,
			),
		)
		var playerResponse Data
		if err := json.Unmarshal(response, &playerResponse); err != nil {
			log.Fatalf("Error unmarshalling json: %v\n", err)
		}

		p.logTheResponse(playerResponse)

		playersJson, err := json.Marshal(playerResponse.Players)
		if err != nil {
			log.Fatalf("Error marshalling json: %v\n", err)
		}
		p.redisClient.Set(context.Background(), cacheKey, playersJson, 24*30*time.Hour)

		return playerResponse.Players
	} else if err != nil {
		log.Fatalf("Error getting from redis: %v\n", err)
	}

	p.logger.Info("Cache hit, returning cached data")
	var players []Player
	if err := json.Unmarshal([]byte(cachedPlayers), &players); err != nil {
		log.Fatalf("Error unmarshalling cached json: %v\n", err)
	}

	return players
}

func (p PlayersApi) logTheResponse(playerResponse Data) {
	playerNames := make([]string, 0, 3)
	for i, p := range playerResponse.Players {
		if i >= 3 {
			break
		}
		playerNames = append(playerNames, p.Name)
	}

	p.logger.WithFields(logrus.Fields{
		"firstThreeNames": playerNames,
	}).Info("Rapid api response summary with player names")
}

func request(url string) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Add("X-RapidAPI-Key", os.Getenv("RAPID_API_KEY"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	return body
}
