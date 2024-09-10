package ports

import "github.com/fmo/player-images/internal/application/core/domain"

type PlayerData interface {
	GetPlayers(season, teamId int) []domain.Player
}
