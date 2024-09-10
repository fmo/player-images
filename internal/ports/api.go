package ports

import "context"

type APIPorts interface {
	SavePlayerImage(ctx context.Context, season, teamId int) error
}
