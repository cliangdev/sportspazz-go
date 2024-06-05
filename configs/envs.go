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
	DBAddress         string
	DBName            string
	DBMigrationDir    string
	FirebaseProjectID string
	FirebaseApiKey    string
	GoogleMapApiKey   string
	GCPApiKey         string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		Host:              getEnv("HOST", "localhost"),
		Port:              getEnv("PORT", "4001"),
		DBUser:            getEnv("DB_USER", "user"),
		DBPassword:        getEnv("DB_PASSWORD", "password"),
		DBAddress:         getEnv("DB_ADDRESS", "localhost:3309"),
		DBName:            getEnv("DB_NAME", "sports_app"),
		DBMigrationDir:    getEnv("DB_MIGRATION_DIR", "/db/migrations"),
		FirebaseProjectID: getEnv("FIREBASE_PROJECT_ID", "firebase"),
		FirebaseApiKey:    getEnv("FIREBASE_API_KEY", "apiKey"),
		GoogleMapApiKey:   getEnv("GOOGLE_MAP_API_KEY", "apiKey"),
		GCPApiKey:         getEnv("GCP_SERVICE_ACCOUNT_API_KEY", "apiKey"),
	}
}

func getEnv(key, _default string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return _default
}
