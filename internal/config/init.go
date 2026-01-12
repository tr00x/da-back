package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func loadEnvVariable(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}

type Config struct {
	ClientID              string
	MIGRATE               string
	MIGRATE_PATH          string
	DB_HOST               string
	DB_PORT               string
	DB_USER               string
	DB_PASSWORD           string
	DB_NAME               string
	DB_URL                string
	ACCESS_KEY            string
	ACCESS_TIME           time.Duration
	REFRESH_KEY           string
	REFRESH_TIME          time.Duration
	APP_VERSION           string
	PORT                  string
	APP_MODE              string
	LOGGER_FOLDER_PATH    string
	LOGGER_FILENAME       string
	STATIC_PATH           string
	DEFAULT_IMAGE_WIDTHS  []uint
	IMAGE_BASE_URL        string
	FILE_BASE_URL         string
	FIREBASE_ACCOUNT_FILE string
	FCM_CHANNEL_ID        string
	TWILIO_ACCOUNT_SID    string
	TWILIO_AUTH_TOKEN     string
	TWILIO_PHONE_NUMBER   string
}

var ENV Config

func Init() *Config {
	godotenv.Load(".env")

	ENV.PORT = loadEnvVariable("PORT")

	ENV.ClientID = loadEnvVariable("GOOGLE_CLIENT_ID")

	ENV.DB_HOST = loadEnvVariable("DB_HOST")
	ENV.DB_PORT = loadEnvVariable("DB_PORT")
	ENV.DB_USER = loadEnvVariable("DB_USER")
	ENV.DB_PASSWORD = loadEnvVariable("DB_PASSWORD")
	ENV.DB_NAME = loadEnvVariable("DB_NAME")

	ENV.APP_MODE = loadEnvVariable("APP_MODE")

	ENV.LOGGER_FOLDER_PATH = loadEnvVariable("LOGGER_FOLDER_PATH")
	ENV.LOGGER_FILENAME = loadEnvVariable("LOGGER_FILENAME")

	ENV.ACCESS_KEY = loadEnvVariable("ACCESS_KEY")
	AT, _ := time.ParseDuration(loadEnvVariable(("ACCESS_TIME")))
	ENV.ACCESS_TIME = AT
	ENV.REFRESH_KEY = loadEnvVariable("REFRESH_KEY")
	RT, _ := time.ParseDuration(loadEnvVariable(("REFRESH_TIME")))
	ENV.REFRESH_TIME = RT

	ENV.APP_VERSION = loadEnvVariable("APP_VERSION")
	ENV.STATIC_PATH = loadEnvVariable("STATIC_PATH")
	ENV.MIGRATE = loadEnvVariable("MIGRATE")
	ENV.DEFAULT_IMAGE_WIDTHS = []uint{320, 640} // if change these sizes, u must change in pkg/files.go too (line: 147)
	ENV.IMAGE_BASE_URL = loadEnvVariable("IMAGE_BASE_URL")
	ENV.FILE_BASE_URL = loadEnvVariable("FILE_BASE_URL")
	ENV.FIREBASE_ACCOUNT_FILE = loadEnvVariable("FIREBASE_ACCOUNT_FILE")
	ENV.FCM_CHANNEL_ID = loadEnvVariable("FCM_CHANNEL_ID")
	ENV.TWILIO_ACCOUNT_SID = loadEnvVariable("TWILIO_ACCOUNT_SID")
	ENV.TWILIO_AUTH_TOKEN = loadEnvVariable("TWILIO_AUTH_TOKEN")
	ENV.TWILIO_PHONE_NUMBER = loadEnvVariable("TWILIO_PHONE_NUMBER")
	ENV.MIGRATE_PATH = loadEnvVariable("MIGRATE_PATH")
	return &ENV
}
