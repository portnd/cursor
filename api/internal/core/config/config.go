package config

import (
	"fmt"
	"os"

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
