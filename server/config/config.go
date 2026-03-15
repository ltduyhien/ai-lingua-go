package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	GRPCPort        string
	RedisAddr       string
	RedisTTLSeconds int
	OllamaModel     string
	OllamaBaseURL   string
	HTTPPort        string
}

func Load() (*Config, error) {
	cfg := &Config{}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}
	cfg.GRPCPort = grpcPort

	redisAddr := os.Getenv("REDIS_ADDR")
	cfg.RedisAddr = redisAddr

	redisTTLSecondsRaw := os.Getenv("REDIS_TTL_SECONDS")
	defaultRedisTTLSeconds := 300
	redisTTLSeconds := defaultRedisTTLSeconds
	if redisTTLSecondsRaw != "" {
		if parsed, err := strconv.ParseInt(redisTTLSecondsRaw, 10, 64); err == nil && parsed > 0 {
			redisTTLSeconds = int(parsed)
		}
	}
	cfg.RedisTTLSeconds = redisTTLSeconds

	ollamaModel := os.Getenv("OLLAMA_MODEL")
	cfg.OllamaModel = ollamaModel

	ollamaBaseURL := os.Getenv("OLLAMA_BASE_URL")
	if ollamaBaseURL == "" {
		ollamaBaseURL = "http://localhost:11434"
	}
	cfg.OllamaBaseURL = ollamaBaseURL

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}
	cfg.HTTPPort = httpPort

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.GRPCPort == "" {
		return fmt.Errorf("missing GRPC_PORT configuration")
	}

	if c.OllamaModel == "" {
		return fmt.Errorf("missing OLLAMA_MODEL configuration")
	}

	if c.RedisTTLSeconds <= 0 {
		return fmt.Errorf("REDIS_TTL_SECONDS must be positive (got %d)", c.RedisTTLSeconds)
	}

	if c.RedisAddr != "" {
		ttlDuration := time.Duration(c.RedisTTLSeconds) * time.Second
		if ttlDuration < time.Second {
			return fmt.Errorf("REDIS_TTL_SECONDS too small: %s", ttlDuration.String())
		}
	}

	return nil
}
