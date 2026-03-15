package config

// The import block lists packages this file depends on.
// We import "os" to read environment variables that hold configuration values.
// We import "strconv" to convert string values (from the environment) into integers.
// We import "time" so we can convert seconds (from the environment) into time.Duration values later if needed.
// We import "fmt" to build readable error messages during configuration validation and loading.
import (
	"os"
	"strconv"
	"time"
	"fmt"
)

// Config holds all configuration for the ai-lingua-go server.
// It is populated from environment variables (12-factor style) and validated at startup.
type Config struct {
	// GRPCPort is the port the gRPC server listens on (e.g. "50051").
	// Env: GRPC_PORT
	GRPCPort string

	// RedisAddr is the Redis server address (e.g. "localhost:6379").
	// Env: REDIS_ADDR. Empty means cache is disabled.
	RedisAddr string

	// RedisTTLSeconds is how long translations stay in cache.
	// Env: REDIS_TTL_SECONDS
	RedisTTLSeconds int

	// OllamaModel is the model name for Ollama (e.g. "llama2").
	// Env: OLLAMA_MODEL
	OllamaModel string

	// OllamaBaseURL is the Ollama API base URL (e.g. "http://localhost:11434").
	// Env: OLLAMA_BASE_URL
	OllamaBaseURL string
}

// Load reads configuration values from environment variables and returns a Config instance.
// It applies default values where appropriate and returns an error if required settings are missing or invalid.
func Load() (*Config, error) {
	// Create an empty Config value that we will fill field by field.
	cfg := &Config{}

	// Read the GRPC_PORT environment variable for the gRPC server port.
	grpcPort := os.Getenv("GRPC_PORT")
	// If GRPC_PORT is empty, fall back to a safe default port value (50051).
	if grpcPort == "" {
		grpcPort = "50051"
	}
	// Store the (possibly defaulted) gRPC port into the Config struct.
	cfg.GRPCPort = grpcPort

	// Read the REDIS_ADDR environment variable for the Redis server address.
	redisAddr := os.Getenv("REDIS_ADDR")
	// Store the Redis address (which may be empty to indicate that Redis is disabled) into the Config struct.
	cfg.RedisAddr = redisAddr

	// Read the REDIS_TTL_SECONDS environment variable for the cache time-to-live value.
	redisTTLSecondsRaw := os.Getenv("REDIS_TTL_SECONDS")
	// Define a default TTL (in seconds) that will be used if the environment variable is not set or invalid.
	defaultRedisTTLSeconds := 300
	// Start with the default TTL value.
	redisTTLSeconds := defaultRedisTTLSeconds
	// If the raw TTL string from the environment is not empty we attempt to parse it as an integer.
	if redisTTLSecondsRaw != "" {
		// Parse the TTL string as a base-10 integer with 64-bit precision.
		if parsed, err := strconv.ParseInt(redisTTLSecondsRaw, 10, 64); err == nil && parsed > 0 {
			// If parsing succeeded and the value is positive, use the parsed value instead of the default.
			redisTTLSeconds = int(parsed)
		}
	}
	// Store the final TTL value in seconds into the Config struct.
	cfg.RedisTTLSeconds = redisTTLSeconds

	// Read the OLLAMA_MODEL environment variable for the Ollama model name.
	ollamaModel := os.Getenv("OLLAMA_MODEL")
	// Store the model name string into the Config struct (it will be validated later).
	cfg.OllamaModel = ollamaModel

	// Read the OLLAMA_BASE_URL environment variable for the Ollama API base URL.
	ollamaBaseURL := os.Getenv("OLLAMA_BASE_URL")
	// If the base URL is empty we fall back to the standard local Ollama endpoint.
	if ollamaBaseURL == "" {
		ollamaBaseURL = "http://localhost:11434"
	}
	// Store the (possibly defaulted) Ollama base URL into the Config struct.
	cfg.OllamaBaseURL = ollamaBaseURL

	// Call the Validate method to ensure all required fields are present and consistent.
	if err := cfg.Validate(); err != nil {
		// If validation fails we return a nil config pointer and the validation error.
		return nil, err
	}

	// If we reach this line the configuration is valid and fully populated, so we return it with a nil error.
	return cfg, nil
}

// Validate checks the Config fields for required values and reasonable constraints.
// It returns an error describing the first validation problem that it encounters.
func (c *Config) Validate() error {
	// If the gRPC port field is empty we cannot start a server, so we treat it as a configuration error.
	if c.GRPCPort == "" {
		// Build and return an error that clearly explains which configuration value is missing.
		return fmt.Errorf("missing GRPC_PORT configuration")
	}

	// If the Ollama model name is empty the translation logic would not know which model to call.
	if c.OllamaModel == "" {
		// Return an error that tells the operator to set the OLLAMA_MODEL environment variable.
		return fmt.Errorf("missing OLLAMA_MODEL configuration")
	}

	// If the Redis TTL is zero or negative it would not make sense for caching logic.
	if c.RedisTTLSeconds <= 0 {
		// Return an error that indicates that the TTL must be a positive number of seconds.
		return fmt.Errorf("REDIS_TTL_SECONDS must be positive (got %d)", c.RedisTTLSeconds)
	}

	// If we have a non-empty Redis address we could optionally check that the TTL is not unreasonably small or large.
	// We convert the TTL into a time.Duration to reason about it in terms of time.
	if c.RedisAddr != "" {
		// Construct a time.Duration value from the integer number of seconds.
		ttlDuration := time.Duration(c.RedisTTLSeconds) * time.Second
		// If the duration is less than a second we consider it too small to be useful for caching.
		if ttlDuration < time.Second {
			// Return an error that explains that the TTL is too small.
			return fmt.Errorf("REDIS_TTL_SECONDS too small: %s", ttlDuration.String())
		}
	}

	// If we reach this line all validation checks have passed, so we return nil to indicate success.
	return nil
}
