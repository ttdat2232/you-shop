package size

import (
	"github.com/TechwizsonORG/product-service/entity"
	"github.com/TechwizsonORG/product-service/err"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type SizeService struct {
	sizeRepo SizeRepository
	logger   zerolog.Logger
}

func NewSizeService(sizeRepo SizeRepository, logger zerolog.Logger) *SizeService {
	logger = logger.With().Str("usecase", "size").Logger()
	return &SizeService{
		logger:   logger,
		sizeRepo: sizeRepo,
	}
}

func (s *SizeService) AddSize(name string) (*entity.Size, err.ApplicationError) {
	size := &entity.Size{
		Id:   uuid.New(),
		Name: name,
	}
	addErr := s.sizeRepo.AddSize(size)
	if addErr != nil {
		s.logger.Error().Err(addErr).Msg("")
		return nil, err.CommonError()
	}
	return size, nil

}
func (s *SizeService) GetSizes() []entity.Size {
	sizes, getErr := s.sizeRepo.GetSizes()
	if getErr != nil {
		s.logger.Error().Err(getErr).Msg("")
		return make([]entity.Size, 0)
	}
	return sizes
}
