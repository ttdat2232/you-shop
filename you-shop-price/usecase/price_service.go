package usecase

import (
	"fmt"

	"github.com/TechwizsonORG/price-service/entity"
	"github.com/TechwizsonORG/price-service/err"
	"github.com/TechwizsonORG/price-service/usecase/event"
	"github.com/TechwizsonORG/price-service/usecase/rpc/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type PriceService struct {
	logger    zerolog.Logger
	priceRepo Repository
}

func NewPriceService(logger zerolog.Logger, priceRepo Repository) *PriceService {
	return &PriceService{
		logger:    logger,
		priceRepo: priceRepo,
	}
}

func (ps *PriceService) GetCurrentPrices(productIds []string) (map[string]float64, *err.AppError) {
	result := make(map[string]float64)
	if len(productIds) == 0 {
		return result, nil
	}
	uuids := make([]uuid.UUID, len(productIds))
	for i, id := range productIds {
		parsedId, err := uuid.Parse(id)
		if err != nil {
			ps.logger.Error().Err(err).Msg("Failed to parse product id")
			return result, nil
		}
		uuids[i] = parsedId
	}
	prices, err := ps.priceRepo.GetCurrentPricesByProductIds(uuids)
	if err != nil {
		ps.logger.Error().Err(err).Msg("Failed to get current prices")
		return result, nil
	}
	for _, price := range prices {
		result[price.ProductId.String()] = price.Amount
	}
	return result, nil
}

func (ps *PriceService) CreateNewPriceList(description string, currency entity.Currency) (*entity.PriceList, *err.AppError) {
	nPriceList := &entity.PriceList{
		AuditEntity: entity.AuditEntity{
			Id: uuid.New(),
		},
	}
	createdPriceList, createErr := ps.priceRepo.AddNewPriceList(*nPriceList)

	if createErr != nil {
		ps.logger.Error().Err(createErr).Msg("")
		return nil, err.NewAppError(500, "Failed to create price list", "Failed to create price list", nil)
	}
	return createdPriceList, nil
}
func (p *PriceService) UpdatePrice(productId, colorId, sizeId uuid.UUID, price float64) (bool, *err.AppError) {
	p.logger.Debug().Msg("Updating price")
	priceEntity, getPriceErr := p.priceRepo.GetPrice(productId, colorId, sizeId)
	if getPriceErr != nil {
		return false, err.NewAppError(500, "getting price failed", "getting price failed", nil)
	}
	if priceEntity == nil {
		return false, err.NewAppError(400, "counldn't found price", "counldn't found price", nil)
	}
	if priceEntity.Amount == price {
		p.logger.Debug().Msg("Price wasn't change")
		return true, nil
	}
	ok, updateErr := p.priceRepo.UpdatePrice(productId, sizeId, colorId, price)

	if !ok || updateErr != nil {
		if updateErr != nil {
			p.logger.Error().Err(updateErr).Msg("")
		}
		return false, err.NewAppError(500, "updated fail", "updated fail", nil)
	}

	p.logger.Info().Msgf("Update price of product's %s successfully", productId.String())
	return true, nil
}

func (p *PriceService) GetTotalPrice(req model.TotalPriceRequest) (float64, []entity.Price, *err.AppError) {

	prices, getPriceErr := p.priceRepo.GetCurrentPrices(req.Items)
	if getPriceErr != nil {
		p.logger.Error().Err(getPriceErr).Msg("")
		return 0, nil, nil
	}

	var totalPrice float64 = 0
	priceMap := make(map[string]float64, len(prices))

	for _, price := range prices {
		key := fmt.Sprintf("%s-%s-%s", price.ProductId, price.ColorId, price.SizeId)
		priceMap[key] = price.Amount
	}

	for _, item := range req.Items {
		key := fmt.Sprintf("%s-%s-%s", item.ProductId, item.ColorId, item.SizeId)
		totalPrice += priceMap[key] * float64(item.Quantity)
	}
	return totalPrice, prices, nil
}

func (p *PriceService) CreateNewPrices(event event.CreatedInventoriesEvent) ([]*entity.Price, *err.AppError) {
	prices := make([]*entity.Price, 0, len(event.CreatedInventories))
	for _, inventory := range event.CreatedInventories {
		prices = append(prices, entity.NewPrice(inventory.Price, inventory.ProductId, inventory.ColorId, inventory.SizeId))
	}
	if addErr := p.priceRepo.AddNewPrices(prices); addErr != nil {
		p.logger.Error().Err(addErr).Msg("")
		return nil, err.NewAppError(500, "adding prices failed", "adding prices failed", nil)
	}
	return prices, nil
}
