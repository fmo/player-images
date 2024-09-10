package config

import (
	"os"
)

func GetS3Bucket() string {
	return os.Getenv("S3_BUCKET")
}

func GetAwsRegion() string {
	return os.Getenv("AWS_REGION")
}

func GetRedisAddr() string {
	return os.Getenv("REDIS_ADDR")
}

func GetRedisPassword() string {
	return os.Getenv("REDIS_PASSWORD")
}

func GetRapidApiKey() string {
	return os.Getenv("RAPID_API_KEY")
}
