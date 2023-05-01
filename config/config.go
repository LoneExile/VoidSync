package config

type Config struct {
	MinIOEndpoint    string
	MinIOAccessKeyID string
	MinIOSecretKey   string
	MinIOUseSSL      bool
	MinIOBucketName  string
}

func LoadConfig() *Config {
	return &Config{
		MinIOEndpoint:    "192.168.1.102:9009",
		MinIOAccessKeyID: "",
		MinIOSecretKey:   "",
		MinIOUseSSL:      false, // Change to "true" if you are using https
		MinIOBucketName:  "blog",
	}
}
