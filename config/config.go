package config

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
		MinIOAccessKeyID: "",
		MinIOSecretKey:   "",
		StorageType: "minio",
		// MinIOEndpoint: "localhost:9009",
		MinIOEndpoint:       "minio:9000",
		MinIOUseSSL:         false, // Change to "true" if you are using https
		MinIOBucketName:     "blog",
		MaxDownloadAttempts: 5,
		// SurrealDBEndpoint: "ws://localhost:8000/rpc",
		SurrealDBEndpoint: "ws://surrealdb:8000/rpc",
	}
}
