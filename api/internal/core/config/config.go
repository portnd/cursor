package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string
	AppPort string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresDSN      string

	MongoHost     string
	MongoPort     string
	MongoUser     string
	MongoPassword string
	MongoDB       string
	MongoURI      string

	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisAddr     string

	JWTSecret string

	GoogleAPIKey   string
	GeminiAPIKey   string // GEMINI_API_KEY for AI estimate/code review
	GroqAPIKey     string // GROQ_API_KEY — when set, use Groq instead of Gemini

	// Optional: AI quota limits for usage display (defaults in code: 15 RPM, 250 RPD if unset)
	AILimitRPM int // e.g. 15 free, 1000 paid
	AILimitRPD int // e.g. 250 free, 10000 paid
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8080"),

		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     getEnv("POSTGRES_USER", "komgrip"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "komgrip_secret"),
		PostgresDB:       getEnv("POSTGRES_DB", "komgrip_db"),

		MongoHost:     getEnv("MONGO_HOST", "localhost"),
		MongoPort:     getEnv("MONGO_PORT", "27017"),
		MongoUser:     getEnv("MONGO_USER", "komgrip"),
		MongoPassword: getEnv("MONGO_PASSWORD", "komgrip_secret"),
		MongoDB:       getEnv("MONGO_DB", "komgrip_logs"),

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", "komgrip_secret"),

		JWTSecret: getEnv("JWT_SECRET", "default_jwt_secret_change_in_production"),

		GoogleAPIKey: getEnv("GOOGLE_API_KEY", ""),
		GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),
		GroqAPIKey:   getEnv("GROQ_API_KEY", ""),

		AILimitRPM: getEnvInt("AI_LIMIT_RPM", 0), // 0 = use default in usage tracker
		AILimitRPD: getEnvInt("AI_LIMIT_RPD", 0),
	}

	cfg.PostgresDSN = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
		cfg.PostgresPort,
	)

	cfg.MongoURI = fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?authSource=admin",
		cfg.MongoUser,
		cfg.MongoPassword,
		cfg.MongoHost,
		cfg.MongoPort,
		cfg.MongoDB,
	)

	cfg.RedisAddr = fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if n, err := strconv.Atoi(value); err == nil {
			return n
		}
	}
	return fallback
}
