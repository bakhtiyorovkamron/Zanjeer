package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

// Config ...
type Config struct {
	OtpTimeout                int // seconds
	ContextTimeout            int
	Environment               string // develop, staging, production
	LogLevel                  string // DEBUG, INFO ...
	HTTPPort                  string
	PostgresHost              string
	PostgresPort              string
	PostgresDatabase          string
	PostgresUser              string
	PostgresPassword          string
	PostgresConnectionTimeOut int // seconds
	PostgresConnectionTry     int
	BaseUrl                   string
	SMTPEmail                 string
	SMTPEmailPass             string
	SMTPHost                  string
	SMTPPort                  string
	SignInKey                 string
	AuthConfigPath            string
	CSVFilePath               string
	RedisHost                 string
	RedisPort                 string
	AccessTokenTimout         int // MINUTES
	MaxImageSize              int // Mb
	DefaultOtp                string
}

// Load loads environment vars and inflates Config
func Load() Config {
	dotenvFilePath := cast.ToString(getOrReturnDefault("DOT_ENV_PATH", "config/.env"))
	err := godotenv.Load(dotenvFilePath)

	if err != nil {
		fmt.Println(".env file not found. Default configuration is being used.")
	}
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))
	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "DEBUG"))
	c.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", "8000"))
	c.BaseUrl = cast.ToString(getOrReturnDefault("BASE_URL", "http://localhost:8000/v1/"))

	// Postgres
	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "204.12.199.7"))
	c.PostgresPort = cast.ToString(getOrReturnDefault("POSTGRES_PORT", 5432))
	c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DATABASE", "zanjeer"))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "zanjeer"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "zanjer"))
	c.PostgresConnectionTimeOut = cast.ToInt(getOrReturnDefault("POSTGRES_CONNECTION_TIMEOUT", 5))
	c.PostgresConnectionTry = cast.ToInt(getOrReturnDefault("POSTGRES_CONNECTION_TRY", 1))

	c.SignInKey = cast.ToString(getOrReturnDefault("SIGN_IN_KEY", "ASJDKLFJASasdFASE2SD2dafa"))
	c.AuthConfigPath = cast.ToString(getOrReturnDefault("AUTH_CONFIG_PATH", "./config/auth.conf"))
	c.CSVFilePath = cast.ToString(getOrReturnDefault("CSV_FILE_PATH", "./config/auth.csv"))
	// Email sending
	c.SMTPEmail = cast.ToString(getOrReturnDefault("SMTP_EMAIL", "youremail@gmail.com"))
	c.SMTPEmailPass = cast.ToString(getOrReturnDefault("SMTP_EMAIL_PASS", "YOUR_EMAIL_PASSWORD"))
	c.SMTPHost = cast.ToString(getOrReturnDefault("SMTP_HOST", "smtp.gmail.com"))
	c.SMTPPort = cast.ToString(getOrReturnDefault("SMTP_PORT", "587"))

	// in mermory storage
	c.RedisHost = cast.ToString(getOrReturnDefault("REDIS_HOST", "localhost"))
	c.RedisPort = cast.ToString(getOrReturnDefault("REDIS_PORT", "6379"))
	c.OtpTimeout = cast.ToInt(getOrReturnDefault("OTP_TIMEOUT", 300))
	c.ContextTimeout = cast.ToInt(getOrReturnDefault("CONTEXT_TIMOUT", 7))
	c.AccessTokenTimout = cast.ToInt(getOrReturnDefault("ACCESS_TOKEN_TIMEOUT", 300))

	// Media
	c.MaxImageSize = cast.ToInt(getOrReturnDefault("MAX_IMAGE_SIZE", 5))

	c.DefaultOtp = cast.ToString(getOrReturnDefault("DEFAULT_OTP", "111111"))
	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}
	return defaultValue
}
