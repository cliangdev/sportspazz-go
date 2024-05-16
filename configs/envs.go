package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host                       string
	Port                       string
	DBUser                     string
	DBPassword                 string
	DBAddress                  string
	DBName                     string
	DBMigrationDir             string
	FirebaseServiceAccountJson string
	FirebaseApiKey             string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		Host:                       getEnv("HOST", "localhost"),
		Port:                       getEnv("PORT", "4001"),
		DBUser:                     getEnv("DB_USER", "user"),
		DBPassword:                 getEnv("DB_PASSWORD", "password"),
		DBAddress:                  getEnv("DB_ADDRESS", "localhost:3309"),
		DBName:                     getEnv("DB_NAME", "sports_app"),
		DBMigrationDir:             getEnv("DB_MIGRATION_DIR", "/db/migrations"),
		FirebaseServiceAccountJson: getEnv("FIREBASE_SERVICE_ACCOUNT_JSON", "firebase-service-account.json"),
		FirebaseApiKey:             getEnv("FIREBASE_API_KEY", "apiKey"),
	}
}

func getEnv(key, _default string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return _default
}
