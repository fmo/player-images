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
