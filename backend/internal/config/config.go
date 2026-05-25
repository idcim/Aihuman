package config

import "os"

type Config struct {
	Port           string
	AIServiceURL   string
	VideoEngineURL string
	DatabaseURL    string
	RedisAddr      string
	RedisStream    string
	WorkerEnabled  bool
}

func Load() Config {
	return Config{
		Port:           env("BACKEND_PORT", "8080"),
		AIServiceURL:   env("AI_SERVICE_URL", "http://ai-service:8001"),
		VideoEngineURL: env("VIDEO_ENGINE_URL", "http://video-engine:8002"),
		DatabaseURL:    env("DATABASE_URL", ""),
		RedisAddr:      env("REDIS_ADDR", ""),
		RedisStream:    env("REDIS_STREAM", "video_tasks"),
		WorkerEnabled:  env("WORKER_ENABLED", "true") != "false",
	}
}

func env(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
