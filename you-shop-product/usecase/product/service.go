package product

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/TechwizsonORG/product-service/config/model"
	"github.com/TechwizsonORG/product-service/entity"
	"github.com/TechwizsonORG/product-service/err"
	"github.com/TechwizsonORG/product-service/usecase/inventory"
	messagequeue "github.com/TechwizsonORG/product-service/usecase/message_queue"
	"github.com/TechwizsonORG/product-service/usecase/rpc"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	productRepo   ProductRepository
	inventoryRepo inventory.InventoryRepository
	logger        zerolog.Logger
	msgQueue      messagequeue.MessageQueue
	rpcService    rpc.RpcInterface
	httpEndpoint  model.HttpEndpoint
	rpcEndpoint   model.RpcServerEndpoint
}

func NewService(httpEndpoint model.HttpEndpoint, repo ProductRepository, log zerolog.Logger, msgQueue messagequeue.MessageQueue, rpcService rpc.RpcInterface, rpcEndpoint model.RpcServerEndpoint, inventoryRepo inventory.InventoryRepository) *Service {
	logger := log.
		With().
		Str("product", "service").
		Logger()
	return &Service{
		productRepo:   repo,
		logger:        logger,
		msgQueue:      msgQueue,
		httpEndpoint:  httpEndpoint,
		rpcService:    rpcService,
		rpcEndpoint:   rpcEndpoint,
		inventoryRepo: inventoryRepo,
	}
}

func (s *Service) SearchProducts(query string) []entity.Product {
	products, err := s.productRepo.Search(query)
	if err != nil {
		s.logger.Err(err).Msg("")
		return []entity.Product{}
	}
	return products
}

func (s *Service) GetProducts(page int, pageSize int) (count int, products []entity.Product) {
	count, err := s.productRepo.Count()
	if err != nil {
		return 0, []entity.Product{}
	}
	products, err = s.productRepo.List(page, pageSize)
	if err != nil {
		return count, []entity.Product{}
	}
	return count, products
}

func (s *Service) GetProduct(id string) (product *entity.Product, appErr err.ApplicationError) {
	entityID, parseErr := uuid.Parse(id)
	if parseErr != nil {
		return nil, err.NewProductError(400, "Cannot parse UUID", "", nil)
	}
	result, appErr := s.productRepo.Get(entityID)
	return result, appErr
}

func (s *Service) CreateProduct(name, description, sku, userManual string, productImages map[string]*multipart.File, thumbnailImage *multipart.FileHeader) (*entity.Product, err.ApplicationError) {

	thumbnailFile, openErr := thumbnailImage.Open()
	if openErr != nil {
		return nil, err.NewProductError(500, "couldn't open thumbnail image", "", nil)
	}
	defer thumbnailFile.Close()
	nProduct := entity.NewProduct(name, description, sku, userManual)
	thumbnailImageUrls := s.uploadImages(*nProduct, map[string]*multipart.File{thumbnailImage.Filename: &thumbnailFile})
	if len(thumbnailImageUrls) > 0 {
		nProduct.Thumbnail = thumbnailImageUrls[0]
	}
	if len(productImages) > 0 {
		go s.uploadImages(*nProduct, productImages)
	}
	createdProduct, createErr := s.productRepo.Create(*nProduct)
	if createErr != nil {
		return nil, createErr
	}
	return &createdProduct, nil
}

func (s *Service) uploadImages(product entity.Product, files map[string]*multipart.File) []string {
	var mu sync.Mutex
	imageUrls := make([]string, 0, len(files))
	var wg sync.WaitGroup
	wg.Add(len(files))
	for key, file := range files {
		content, _ := io.ReadAll(*file)
		buffer := bytes.NewBuffer(content)
		(*file).Close()
		go func(key string, buffer *bytes.Buffer) {
			defer wg.Done()
			imageDest, uploadErr := s.upload(product, key, buffer)
			if uploadErr != nil {
				return
			}
			mu.Lock()
			imageUrls = append(imageUrls, imageDest)
			mu.Unlock()
		}(key, buffer)
	}
	wg.Wait()
	productImages := make([]entity.ProductImage, 0, len(imageUrls))
	for _, imageUrl := range imageUrls {
		productImages = append(productImages, entity.ProductImage{
			ProductId: product.Id,
			ImageUrl:  imageUrl,
			IsPrimary: true,
		})
	}
	s.productRepo.AddProductImages(productImages)
	return imageUrls
}

func (s *Service) upload(product entity.Product, key string, buffer *bytes.Buffer) (string, err.ApplicationError) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, partErr := writer.CreateFormFile("image_file", key)
	if partErr != nil {
		s.logger.Error().Err(partErr).Msg("Error occurred when create part file")
		return "", err.CommonError()
	}
	_, copyErr := io.Copy(part, buffer)
	if copyErr != nil {
		s.logger.Err(copyErr).Msg("Error occurred when copying file")
		return "", err.CommonError()
	}

	writer.WriteField("alt", fmt.Sprintf("%s image", product.Name))
	writer.WriteField("owner_id", product.Id.String())
	closeWriterErr := writer.Close()
	if closeWriterErr != nil {
		s.logger.Error().Err(closeWriterErr).Msg("Error occurred When closing writer")
		return "", err.CommonError()
	}
	req, createReqErr := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/upload", s.httpEndpoint.UploadServerUrl), body)
	if createReqErr != nil {
		s.logger.Error().Err(createReqErr).Msg("Error occurred When creating uploading request")
		return "", err.CommonError()
	}
	client := http.Client{Timeout: 30 * time.Second}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, doErr := client.Do(req)
	if doErr != nil {
		s.logger.Error().Err(doErr).Msg("Error occurred when sending request")
		return "", err.CommonError()
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		s.logger.Error().Msg("Error occurred when uploading image")
		return "", err.CommonError()
	}

	id := res.Header.Get("Location")[strings.LastIndex(res.Header.Get("Location"), "/")+1:]
	return fmt.Sprintf("%s/%s", s.httpEndpoint.UploadServerUrl, id), nil
}

func (s *Service) UpdateProduct(id uuid.UUID, name, description, sku string, status entity.ProductStatus, userManual string) (*entity.Product, err.ApplicationError) {
	product, fetchErr := s.productRepo.Get(id)
	if product == nil || fetchErr != nil {
		if fetchErr != nil {
			s.logger.Error().Err(fetchErr).Msg("Error occurred when fetching product")
		}
		return nil, err.NewProductError(404, "Product not found", "", nil)
	}
	product.Update(name, description, sku, status, userManual)
	updatedProduct, updateErr := s.productRepo.Update(*product)
	if updateErr != nil {
		return nil, updateErr
	}

	return &updatedProduct, nil
}

func (s *Service) DeleteProduct(id uuid.UUID) err.ApplicationError {
	return err.CommonError()
}
func (s *Service) CheckProductQuantity(productId, colorId, sizeId uuid.UUID, requireQuantity int) (bool, err.ApplicationError) {
	currentQuantity, getQuantityErr := s.inventoryRepo.GetQuantity(productId, sizeId, colorId)
	if getQuantityErr != nil {
		s.logger.Error().Err(getQuantityErr).Msg("")
		return false, err.CommonError()
	}

	if currentQuantity < requireQuantity {
		return false, err.NewProductError(400, "Not enough quantity", "Not enought quantity", nil)
	}

	return true, nil
}

func (s *Service) GetProductByIds(productIds []uuid.UUID) []entity.Product {
	products, getErr := s.productRepo.GetByIds(productIds)
	if getErr != nil {
		return []entity.Product{}
	}
	return products
}
func (s *Service) GetQuantity(productId, sizeId, colorId uuid.UUID) (int, err.ApplicationError) {
	quantity, getQuantityErr := s.inventoryRepo.GetQuantity(productId, sizeId, colorId)
	if getQuantityErr != nil {
		s.logger.Error().Err(getQuantityErr).Msg("")
		return 0, err.NewProductError(500, "error occurred", "error occurred", nil)
	}
	return quantity, nil
}

func (s *Service) UploadProductColor(productId, colorId uuid.UUID, file *multipart.FileHeader) err.ApplicationError {
	product, getErr := s.productRepo.Get(productId)
	if getErr != nil {
		s.logger.Error().Err(getErr).Msg("")
		return err.CommonError()
	}
	openedFile, openErr := file.Open()
	if openErr != nil {
		s.logger.Error().Err(openErr).Msg("")
		return err.CommonError()
	}
	defer openedFile.Close()
	content, _ := io.ReadAll(openedFile)
	buffer := bytes.NewBuffer(content)
	dest, uploadErr := s.upload(*product, file.Filename, buffer)
	if uploadErr != nil {
		return uploadErr
	}
	productColorImage := &entity.ProductImage{
		Id:        uuid.New(),
		ProductId: productId,
		ColorId:   colorId,
		ImageUrl:  dest,
		IsPublic:  true,
		IsPrimary: false,
	}
	return s.productRepo.AddProductImages([]entity.ProductImage{*productColorImage})
}
