package cli

import (
	"context"
	"fmt"
	"github.com/fmo/player-images/internal/ports"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var teamId int
var season int

var logger = logrus.New()

type Adapter struct {
	api ports.APIPorts
}

func NewAdapter(api ports.APIPorts) *Adapter {
	return &Adapter{
		api: api,
	}
}

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
	},
}

func (a Adapter) Run(ctx context.Context) {
	rootCmd := &cobra.Command{
		Use:   "football-data-app",
		Short: "Football Data CLI Application",
	}

	Cmd.Run = func(cmd *cobra.Command, args []string) {
		logger.Info("Starting player command with Adapter")

		err := a.api.SavePlayerImage(ctx, season, teamId)
		if err != nil {
			fmt.Println("error: ", err)
		}

	}

	rootCmd.AddCommand(Cmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
