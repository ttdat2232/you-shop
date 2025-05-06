package job

import (
	"context"
	"encoding/json"

	"github.com/TechwizsonORG/image-service/background"
	"github.com/TechwizsonORG/image-service/infrastructure/rpc"
	"github.com/TechwizsonORG/image-service/usecase"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Job struct {
	logger zerolog.Logger
}

func NewJob(logger zerolog.Logger) *Job {
	return &Job{logger: logger}
}

func (j *Job) GetOwnersImages(ctx context.Context, rpcService rpc.RpcInterface, imageService usecase.Service) background.JobFunc {
	return func() {
		rpcService.NewRpcQueue("get_owners_images", func(data string) string {

			var ids []string
			err := json.Unmarshal([]byte(data), &ids)
			if err != nil {
				j.logger.Error().Err(err).Msg("Error unmarshal data:")
				return "{}"
			}
			uuids := []uuid.UUID{}
			for _, id := range ids {
				uuidParsed, parseErr := uuid.Parse(id)
				if parseErr != nil {
					continue
				}
				uuids = append(uuids, uuidParsed)
			}
			result := imageService.GetImageByOwnerIds(uuids, "https", "api.youshop.fun", "api/v1/images")

			jsonResult, err := json.Marshal(result)
			if err != nil {
				j.logger.Error().Err(err).Msg("Error marshal data:")
				return ""
			}
			return string(jsonResult)
		})
	}
}
