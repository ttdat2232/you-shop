//go:build prod

package config

import (
	"strconv"

	"github.com/TechwizsonORG/auth-service/config/model"
	"github.com/joho/godotenv"
)

func Init() (databaseConfig *model.DatabaseConfig, serverConfig *model.ServerConfig, logConfig *model.LogConfig, jwtConfig *model.JwtConfig, mode string) {
	mode = "release"
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
		Level:    int8(logLevel),
		FilePath: envMap["LOG_FILE_PATH"],
	}

	defaultAccessExpireTime, err := strconv.Atoi(envMap["JWT_DEFAULT_ACCESS_EXPIRE_TIME"])
	if err != nil {
		panic("Invalid JWT_DEFAULT_ACCESS_EXPIRE_TIME value")
	}

	defaultRefreshExpireTime, err := strconv.Atoi(envMap["JWT_DEFAULT_REFRESH_EXPIRE_TIME"])
	if err != nil {
		panic("Invalid JWT_DEFAULT_REFRESH_EXPIRE_TIME value")
	}

	jwtConfig = &model.JwtConfig{
		SecretKey:                envMap["JWT_SECRET_KEY"],
		DefaultAccessExpireTime:  defaultAccessExpireTime,
		DefaultRefreshExpireTime: defaultRefreshExpireTime,
		Issuer:                   envMap["JWT_ISSUER"],
	}
	return databaseConfig, serverConfig, logConfig, jwtConfig, mode
}
