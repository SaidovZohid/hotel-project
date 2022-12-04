package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	HttpPort  string
	RedisAddr string
	Auth      Authorization
	Postgres  PostgresConfig
	Smtp      Smtp
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type Authorization struct {
	AuthHeaderKey  string
	AuthpayloadKey string
	SecretKey      string
}

type Smtp struct {
	Sender   string
	Password string
}

func Load(path string) Config {
	godotenv.Load(path + "/.env")
	conf := viper.New()
	conf.AutomaticEnv()
	cfg := Config{
		HttpPort:  conf.GetString("HTTP_PORT"),
		RedisAddr: conf.GetString("REDIS_ADDR"),
		Postgres: PostgresConfig{
			Host:     conf.GetString("POSTGRES_HOST"),
			Port:     conf.GetString("POSTGRES_PORT"),
			User:     conf.GetString("POSTGRES_USER"),
			Password: conf.GetString("POSTGRES_PASSWORD"),
			Database: conf.GetString("POSTGRES_DATABASE"),
		},
		Auth: Authorization{
			AuthHeaderKey:  conf.GetString("AUTHORIZATION_HEADER_KEY"),
			AuthpayloadKey: conf.GetString("AUTHORIZATION_PAYLOAD_KEY"),
			SecretKey:      conf.GetString("SECRET_KEY"),
		},
		Smtp: Smtp{
			Sender:   conf.GetString("SMTP_SENDER"),
			Password: conf.GetString("SMTP_PASSWORD"),
		},
	}

	return cfg
}
