package config

import "os"

type Config struct {
	Port           string
	AIServiceURL   string
	VideoEngineURL string
}

func Load() Config {
	return Config{
		Port:           env("BACKEND_PORT", "8080"),
		AIServiceURL:   env("AI_SERVICE_URL", "http://ai-service:8001"),
		VideoEngineURL: env("VIDEO_ENGINE_URL", "http://video-engine:8002"),
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
