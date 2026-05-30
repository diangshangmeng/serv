package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort         string
	ServerHost         string
	ServerURL          string
	ServerImageURL     string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	RedisHost          string
	RedisPort          string
	RedisPassword      string
	RedisDB            int
	JWTSecret          string
	JWTExpireHours     int
	SMSAccessKeyID     string
	SMSAccessKeySecret string
	ImageUploadPath    string
	MaxImageSize       int64
	LogLevel           string
	LogFilename        string
	LogMaxSize         int
	LogMaxBackups      int
	LogMaxAge          int
	LogCompress        bool
}

var AppConfig *Config

func LoadConfig() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	AppConfig = &Config{
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		ServerHost:         getEnv("SERVER_HOST", "localhost"),
		ServerURL:          getEnv("SERVER_URL", "https://api.dianshangmeng.tech"),
		ServerImageURL:     getEnv("SERVER_IMAGE_URL", "https://api.dianshangmeng.tech/"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "3306"),
		DBUser:             getEnv("DB_USER", "dsm_user"),
		DBPassword:         getEnv("DB_PASSWORD", "asdf3asRDSfEre4DAS79"),
		DBName:             getEnv("DB_NAME", "dsm"),
		RedisHost:          getEnv("REDIS_HOST", "localhost"),
		RedisPort:          getEnv("REDIS_PORT", "6379"),
		RedisPassword:      getEnv("REDIS_PASSWORD", "gx0ifqGcAdl3EMgiPlYBm5orYGZTOq"),
		RedisDB:            getEnvAsInt("REDIS_DB", 6),
		JWTSecret:          getEnv("JWT_SECRET", "xK9#pL2$mN7!qR5*tY8&uI3_oP6^aS4+dF1-gH0=jK2:lZ9;vX4,cV6.bM3?nB5"),
		JWTExpireHours:     getEnvAsInt("JWT_EXPIRE_HOURS", 720),
		SMSAccessKeyID:     getEnv("SMS_ACCESS_KEY_ID", ""),
		SMSAccessKeySecret: getEnv("SMS_ACCESS_KEY_SECRET", ""),
		ImageUploadPath:    getEnv("IMAGE_UPLOAD_PATH", "/uploads"),
		MaxImageSize:       getEnvAsInt64("MAX_IMAGE_SIZE", 2097152),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		LogFilename:        getEnv("LOG_FILENAME", "./logs/app.log"),
		LogMaxSize:         getEnvAsInt("LOG_MAX_SIZE", 100),
		LogMaxBackups:      getEnvAsInt("LOG_MAX_BACKUPS", 30),
		LogMaxAge:          getEnvAsInt("LOG_MAX_AGE", 7),
		LogCompress:        getEnvAsBool("LOG_COMPRESS", false),
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		var result int
		_, err := fmt.Sscanf(value, "%d", &result)
		if err == nil {
			return result
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		var result int64
		_, err := fmt.Sscanf(value, "%d", &result)
		if err == nil {
			return result
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if value == "true" || value == "1" {
			return true
		} else if value == "false" || value == "0" {
			return false
		}
	}
	return defaultValue
}
