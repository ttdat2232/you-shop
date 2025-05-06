package color

import (
	"github.com/TechwizsonORG/product-service/entity"
	"github.com/TechwizsonORG/product-service/err"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type ColorService struct {
	logger    zerolog.Logger
	colorRepo Repository
}

func NewColorService(logger zerolog.Logger, colorRepo Repository) *ColorService {
	logger = logger.With().Str("usecase", "color").Logger()
	return &ColorService{
		logger:    logger,
		colorRepo: colorRepo,
	}
}

func (c *ColorService) AddColor(name string) (*entity.Color, err.ApplicationError) {
	colorEntity := &entity.Color{
		Id:   uuid.New(),
		Name: name,
	}
	addErr := c.colorRepo.AddColor(colorEntity)
	if addErr != nil {
		c.logger.Error().Err(addErr).Msg("")
		return nil, err.NewProductError(500, "adding product failed", "adding product failed", nil)
	}
	return colorEntity, nil
}

func (c *ColorService) GetColors() []entity.Color {
	colors, getErr := c.colorRepo.GetColors()
	if getErr != nil {
		c.logger.Error().Err(getErr).Msg("")
		return make([]entity.Color, 0)
	}
	return colors
}
func (c *ColorService) GetById(uuid.UUID) (entity.Color, err.ApplicationError) {
	return entity.Color{}, nil
}
