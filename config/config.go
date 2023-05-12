package config

import (
	"os"
	"strconv"
)

type Config struct {
	StorageType         string
	MinIOEndpoint       string
	MinIOAccessKeyID    string
	MinIOSecretKey      string
	MinIOUseSSL         bool
	MinIOBucketName     string
	MaxDownloadAttempts int

	SurrealDBEndpoint string
}

func LoadConfig() *Config {
	return &Config{
		// StorageType:         "minio",
		// MinIOAccessKeyID:    "ROOTUSER",
		// MinIOSecretKey:      "CHANGEME123",
		// MinIOUseSSL:         false, // Change to "true" if you are using https
		// MinIOBucketName:     "blog",
		// MaxDownloadAttempts: 5,

		// MinIOEndpoint:     "localhost:9009",
		// SurrealDBEndpoint: "ws://localhost:8000/rpc",

		// MinIOEndpoint:       "minio:9000",
		// SurrealDBEndpoint: "ws://surrealdb:8000/rpc",

		// -------------------------------------------------------

		StorageType:         os.Getenv("STORAGE_TYPE"),
		MinIOEndpoint:       os.Getenv("MINIO_ENDPOINT"),
		MinIOAccessKeyID:    os.Getenv("MINIO_ROOT_USER"),
		MinIOSecretKey:      os.Getenv("MINIO_ROOT_PASSWORD"),
		MinIOUseSSL:         os.Getenv("MINIO_USE_SSL") == "true",
		MinIOBucketName:     os.Getenv("MINIO_BUCKET_NAME"),
		MaxDownloadAttempts: getIntEnv("MAX_DOWNLOAD_ATTEMPTS", 5),
		SurrealDBEndpoint:   os.Getenv("SURREALDB_ENDPOINT"),
	}
}

func getIntEnv(key string, defaultValue int) int {
	if val, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultValue
}
