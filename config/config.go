package config

type Config struct {
	StorageType         string
	MinIOEndpoint       string
	MinIOAccessKeyID    string
	MinIOSecretKey      string
	MinIOUseSSL         bool
	MinIOBucketName     string
	MaxDownloadAttempts int
}

func LoadConfig() *Config {
	return &Config{
		MinIOAccessKeyID: "",
		MinIOSecretKey:   "",
		StorageType:         "minio",
		MinIOEndpoint:       "192.168.1.102:9009",
		MinIOUseSSL:         false, // Change to "true" if you are using https
		MinIOBucketName:     "blog",
		MaxDownloadAttempts: 5,
	}
}
