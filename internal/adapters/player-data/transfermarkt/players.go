package transfermarkt

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fmo/player-images/internal/application/core/domain"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type Adapter struct {
	rapidApiKey string
}

func NewAdapter(rapidApiKey string) (*Adapter, error) {
	return &Adapter{
		rapidApiKey: rapidApiKey,
	}, nil
}

type Data struct {
	Players []domain.Player `json:"data"`
}

func (a Adapter) GetPlayers(season, teamId int) []domain.Player {
	response := a.request(
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

	a.logTheResponse(playerResponse)

	return playerResponse.Players
}

func (a Adapter) logTheResponse(playerResponse Data) {
	playerNames := make([]string, 0, 3)
	for i, p := range playerResponse.Players {
		if i >= 3 {
			break
		}
		playerNames = append(playerNames, p.Name)
	}

	log.WithFields(log.Fields{
		"firstThreeNames": playerNames,
	}).Info("Rapid api response summary with player names")
}

func (a Adapter) request(url string) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Add("X-RapidAPI-Key", a.rapidApiKey)

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
