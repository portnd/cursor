package databases

import (
	"fmt"
	"os"
	"strconv"

	"gorm.io/gorm"
)

// DB is a global var for connect DB
var DB *gorm.DB

// DBConfig represents db configuration
type DBConfig struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
}

// BuildDBConfig use for building DB config
func BuildDBConfig() *DBConfig {
	post, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	dbConfig := DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     post,
		User:     os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_DATABASE"),
	}
	return &dbConfig
}

// DbURL use for create DB connection URL
func DbURL(dbConfig *DBConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)
}
