package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/TechwizsonORG/image-service/config"
	infraFile "github.com/TechwizsonORG/image-service/infrastructure/file"
	"github.com/rs/zerolog"
)

func main() {
	args := os.Args
	_, _, _, _, s3ProxyConfig, _, _ := config.Init()
	multi := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%s]", i))
		},
	})
	log := zerolog.New(multi).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Caller().
		Logger()
	log.Info().
		Str("delete", "img1,img2,img3").
		Msg("Action")
	if len(args) < 2 {
		panic("Please provide arguments")
	}

	commands := strings.Split(args[1], "=")
	if len(commands) < 2 {
		panic("Please provide action and image name")
	}
	action := commands[0]
	imageNames := strings.Split(commands[1], ",")
	fileService := infraFile.NewTypeService(*s3ProxyConfig)
	switch action {
	case "delete":
		var wg sync.WaitGroup
		wg.Add(len(imageNames))
		for i, imageName := range imageNames {
			go func(imageName string, count int) {
				defer wg.Done()
				if err := fileService.DeleteFile(fmt.Sprintf("%s:%d/%s/%s", s3ProxyConfig.Host, s3ProxyConfig.Port, s3ProxyConfig.Folder, imageName)); err != nil {
					log.Error().Err(err).Msgf("Failed to delete image - number %d", count)
				} else {
					log.Info().Msgf("Image deleted successfully - number %d", count)
				}
			}(imageName, i)
		}
		wg.Wait()
	}

}
