package command

import (
	"github.com/fmo/player-images/internal/redis"
	"github.com/fmo/player-images/internal/s3"
	"github.com/fmo/player-images/internal/transfermarkt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var teamId int
var season int

var logger = logrus.New()

func init() {
	Cmd.Flags().IntVarP(&teamId, "teamId", "t", 541, "Team Id")
	Cmd.Flags().IntVarP(&season, "season", "s", 2023, "Season")

	logger.Out = os.Stdout

	logger.Level = logrus.DebugLevel
}

var Cmd = &cobra.Command{
	Use:   "player-images",
	Short: "Get player images from Transfermarkt",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Starting player command")

		s3Service, err := s3.NewS3Service(logger)
		if err != nil {
			logger.Fatalf("Cant connect to s3 %v", err)
		}

		redisClient := redis.NewRedisClient()
		r := transfermarkt.NewPlayersApi(logger, redisClient)
		players := r.GetPlayers(season, teamId)

		for _, player := range players {
			imageAlreadyUploaded := false
			imageInfo := ""
			if player.Image != "" {
				imageAlreadyUploaded, err = s3Service.Save(player.ID, player.Image)
				if err != nil {
					logger.Error(err)
				} else {
					if imageAlreadyUploaded {
						imageInfo = "Image already before uploaded"
					} else {
						imageInfo = "Image uploaded"
					}
				}
			}

			logger.WithFields(logrus.Fields{
				"player image": imageInfo,
				"player id":    player.ID,
				"player name":  player.Name,
			}).Infof("Image update process")
		}

	},
}
