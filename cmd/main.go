package main

import (
	"context"
	"github.com/fmo/player-images/config"
	"github.com/fmo/player-images/internal/adapters/cache/redis"
	"github.com/fmo/player-images/internal/adapters/cli"
	"github.com/fmo/player-images/internal/adapters/image/s3"
	"github.com/fmo/player-images/internal/adapters/player-data/transfermarkt"
	"github.com/fmo/player-images/internal/application/core/api"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	environment := os.Getenv("ENVIRONMENT")
	if environment != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	}

	ctx := context.Background()

	playerAdapter, err := transfermarkt.NewAdapter(config.GetRapidApiKey())
	if err != nil {
		log.Fatalf("Failed to connect to transfermarkt. Error: %v", err)
	}

	imageAdapter, err := s3.NewAdapter(config.GetAwsRegion(), config.GetS3Bucket())
	if err != nil {
		log.Fatalf("Failed to connect to aws. Error: %v", err)
	}

	cacheAdapter, err := redis.NewAdapter(config.GetRedisAddr(), config.GetRedisPassword())
	if err != nil {
		log.Fatalf("Failed to connect to redis. Error: %v", err)
	}

	application := api.NewApplication(imageAdapter, playerAdapter, cacheAdapter)
	cliAdapter := cli.NewAdapter(application)
	cliAdapter.Run(ctx)
}
