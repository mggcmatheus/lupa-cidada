package config

import "os"

type Config struct {
	Port      string
	Env       string
	Debug     bool // Se true, usa dados mockados; se false, usa banco
	MongoURI  string
	RedisURI  string
	MeiliHost string
	MeiliKey  string
}

func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		Env:       getEnv("ENV", "development"),
		Debug:     getEnv("DEBUG", "false") == "true",
		MongoURI:  getEnv("MONGO_URI", "mongodb://lupa:lupa_secret_2024@localhost:27018/lupa_cidada?authSource=admin"),
		RedisURI:  getEnv("REDIS_URI", "redis://localhost:6380"),
		MeiliHost: getEnv("MEILI_HOST", "http://localhost:7701"),
		MeiliKey:  getEnv("MEILI_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
