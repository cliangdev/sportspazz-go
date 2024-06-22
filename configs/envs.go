package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host              string
	Port              string
	DBUser            string
	DBPassword        string
	DBHost            string
	DBPort            string
	DBName            string
	DBMigrationDir    string
	FirebaseProjectID string
	GoogleMapApiKey   string
	GCPApiKey         string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		Host:            getEnv("HOST", "localhost"),
		Port:            getEnv("PORT", "4001"),
		DBUser:          getEnv("DB_USER", "user"),
		DBPassword:      getEnv("DB_PASSWORD", "password"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBName:          getEnv("DB_NAME", "sports_app"),
		DBMigrationDir:  getEnv("DB_MIGRATION_DIR", "/db/migrations"),
		GoogleMapApiKey: getEnv("GOOGLE_MAP_API_KEY", "apiKey"),
		GCPApiKey:       getEnv("GCP_SERVICE_ACCOUNT_API_KEY", "apiKey"),
	}
}

func getEnv(key, _default string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return _default
}
