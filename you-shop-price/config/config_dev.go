//go:build dev

package config

import (
	"strconv"

	"github.com/TechwizsonORG/price-service/config/model"
	"github.com/joho/godotenv"
)

func Init() (databaseConfig *model.DatabaseConfig, serverConfig *model.ServerConfig, logConfig *model.LogConfig, authServerConfig *model.AuthServerConfig, rabbitMqConfig *model.RabbitMqConfig, mode string) {

	mode = "debug"
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	var envMap = map[string]string{}
	envMap, err = godotenv.Read(".env")
	if err != nil {
		panic("Error reading .env file")
	}

	port, err := strconv.Atoi(envMap["DB_PORT"])
	if err != nil {
		panic("Invalid DB_PORT value")
	}

	apiPort, err := strconv.Atoi(envMap["API_SERVER_PORT"])
	if err != nil {
		panic("Invalid API_SERVER_PORT value")
	}

	logLevel, err := strconv.Atoi(envMap["LOG_LEVEL"])
	if err != nil {
		panic("Invalid LOG_LEVEL value")
	}

	databaseConfig = &model.DatabaseConfig{
		Host:     envMap["DB_HOST"],
		Port:     port,
		User:     envMap["DB_USER"],
		Password: envMap["DB_PASSWORD"],
		Name:     envMap["DB_NAME"],
		SSLMode:  envMap["DB_SSLMODE"],
	}

	serverConfig = &model.ServerConfig{
		Host: envMap["API_SERVER_HOST"],
		Port: apiPort,
	}

	logConfig = &model.LogConfig{
		Level:    logLevel,
		FilePath: envMap["LOG_FILE_PATH"],
	}

	authServerConfig = &model.AuthServerConfig{
		AuthServerUrl: envMap["AUTH_SERVER_URL"],
	}

	rabbitMqPort, err := strconv.Atoi(envMap["MSG_BROKER_PORT"])
	if err != nil {
		panic("Invalid RABBITMQ_PORT value")
	}
	rabbitMqConfig = &model.RabbitMqConfig{
		Host:     envMap["MSG_BROKER_HOST"],
		Port:     rabbitMqPort,
		Username: envMap["MSG_BROKER_USERNAME"],
		Password: envMap["MSG_BROKER_PASSWORD"],
		Vhost:    envMap["MSG_BROKER_VHOST"],
	}
	return databaseConfig, serverConfig, logConfig, authServerConfig, rabbitMqConfig, mode
}
