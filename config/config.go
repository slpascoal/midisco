package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Nenhum arquivo .env encontrado. Continuando com as variáveis de ambiente.")
	}

	return &Config{
		DBHost:     getEnv("MYSQL_HOST", "localhost"),
		DBPort:     getEnv("MYSQL_PORT", "3306"),
		DBUser:     getEnv("MYSQL_USER", "root"),
		DBPassword: getEnv("MYSQL_PASSWORD", ""),
		DBName:     getEnv("MYSQL_DB", "midisco"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
