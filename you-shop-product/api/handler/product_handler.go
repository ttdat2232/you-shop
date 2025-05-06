package handler

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/TechwizsonORG/product-service/api/handler/utility"
	"github.com/TechwizsonORG/product-service/api/middleware"
	"github.com/TechwizsonORG/product-service/api/model"
	productModel "github.com/TechwizsonORG/product-service/api/model/product"
	configModel "github.com/TechwizsonORG/product-service/config/model"
	appErr "github.com/TechwizsonORG/product-service/err"
	"github.com/TechwizsonORG/product-service/usecase/inventory"
	inventoryModel "github.com/TechwizsonORG/product-service/usecase/inventory/model"
	"github.com/TechwizsonORG/product-service/usecase/product"
	"github.com/TechwizsonORG/product-service/usecase/rpc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
)

type ProductHandler struct {
	rpcService        rpc.RpcInterface
	productService    product.UseCase
	rpcServerEndpoint configModel.RpcServerEndpoint
	logger            zerolog.Logger
	inventoryService  inventory.InventoryUseCase
}

func NewProductHandler(productService product.UseCase, rpcService rpc.RpcInterface, rpcServerEndpoint configModel.RpcServerEndpoint, logger zerolog.Logger, inventoryService inventory.InventoryUseCase) *ProductHandler {
	logger = logger.With().Str("Handler", "product").Logger()
	return &ProductHandler{
		productService:    productService,
		rpcService:        rpcService,
		rpcServerEndpoint: rpcServerEndpoint,
		logger:            logger,
		inventoryService:  inventoryService,
	}
}

// ProductRoutes is a function that registers the product routes
func (p *ProductHandler) ProductRoutes(router *gin.RouterGroup) {
	productGroup := router.Group("/products")

	productGroup.GET("", p.getProducts)
	productGroup.GET("/quantity", p.getQuantity)
	productGroup.GET(":id", p.getProduct)
	productGroup.POST("", middleware.AuthorizationMiddleware([]string{"admin"}, nil), p.addProduct)
	productGroup.POST("/:id/color", middleware.AuthorizationMiddleware([]string{"admin"}, nil), p.uploadProductColorImage)
	productGroup.POST("/:id/inventory", middleware.AuthorizationMiddleware([]string{"admin"}, nil), p.addProductInventory)
	productGroup.PUT(":id", middleware.AuthorizationMiddleware([]string{"admin"}, nil), p.updateProduct)
	productGroup.PUT("/:id/inventory", middleware.AuthorizationMiddleware([]string{"admin"}, nil), p.updateProductInventory)
}

// GetProduct godoc
//
//	@Summary	Get products data
//	@Tags		products
//	@Produce	json
//	@Param		page		query		int												false	"page number. Default is 1"			Format(int)
//	@Param		page_size	query		int												false	"page_size number. Default is 10"	Format(int)
//	@Failure	400			{object}	model.ApiResponse{data=appErr.ValidationError}	"page or page_size is not a positive number"
//	@Router		/products [get]
func (p *ProductHandler) getProducts(c *gin.Context) {
	isSuccess, validationErr := utility.PaginationValidator(c)
	if !isSuccess {
		c.Errors = append(c.Errors, &gin.Error{Err: validationErr})
		return
	}

	page, pageSize := utility.GetPaginationQuery(c)

	count, products := p.productService.GetProducts(page, pageSize)

	results := []productModel.Product{}
	Ids := make([]string, len(results))
	for _, product := range products {
		productModel := productModel.FromEntity(product)
		results = append(results, productModel)
		Ids = append(Ids, productModel.Id.String())
	}
	jsonReq, _ := json.Marshal(Ids)

	priceJsonRes := make(chan string)
	imageJsonRes := make(chan string)
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	go func() {
		priceJsonRes <- p.rpcService.Req(p.rpcServerEndpoint.ProductsPrice, string(jsonReq))
	}()

	go func() {
		imageJsonRes <- p.rpcService.Req(p.rpcServerEndpoint.GetImageByIds, string(jsonReq))
	}()

	var pricesMap map[string]float64
	if err := json.Unmarshal([]byte(<-priceJsonRes), &pricesMap); err != nil {
		p.logger.Error().Err(err).Msg("Couldn't unmarshal prices")
	} else {
		for i := range results {
			results[i].Price = pricesMap[results[i].Id.String()]
		}
	}

	var imagesMap map[string][]string
	if unmarshalErr := json.Unmarshal([]byte(<-imageJsonRes), &imagesMap); unmarshalErr != nil {
		p.logger.Error().Err(unmarshalErr).Msg("Couldn't unmarshal prices")
	} else {
		for i := range results {
			results[i].Images = imagesMap[results[i].Id.String()]
		}
	}

	c.JSON(http.StatusOK, model.SuccessResponse(model.NewPaginationResponse(page, pageSize, count, results)))
}

// GetProduct godoc
//
//	@Summary	Get product by id
//	@Tags		products
//	@Produce	json
//	@Param		page	path		string	true	"product id in uuid format. Eg: ddb1fdef-2ffb-44a5-a833-fab7b4d60355"
//	@Failure	400		{object}	model.ApiResponse{data=appErr.ProductError}
//	@Router		/products/{id} [get]
func (p *ProductHandler) getProduct(c *gin.Context) {
	productId := c.Param("id")
	product, err := p.productService.GetProduct(productId)

	if err != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: err})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse(productModel.FromEntity(*product)))
}

// AddProduct godoc
//
//	@Summary	Add product
//	@Tags		products
//	@Router		/products [post]
//
//	@Param		name			formData	string	true	"name of product"
//	@Param		sku				formData	string	true	"sku of product"
//	@Param		user_manual		formData	string	true	"user manual of product"
//	@Param		product_images	formData	file	false	"images of product. Max is 10 images"
//	@Param		thumbnail		formData	file	true	"thumbnail of product"
//	@Failure	500,400			{object}	model.ApiResponse{data=appErr.ValidationError}
//	@Success	201
//	@Header		201	{string}	Location	"/api/v1/products/ddb1fdef-2ffb-44a5-a833-fab7b4d60355"
func (p *ProductHandler) addProduct(c *gin.Context) {
	req := c.Request
	name := req.FormValue("name")
	description := req.FormValue("description")
	sku := req.FormValue("sku")
	userManual := req.FormValue("user_manual")

	req.ParseMultipartForm(100 << 20)

	fileheaders := req.MultipartForm.File["product_images"]
	if len(fileheaders) > 10 {
		c.Errors = append(c.Errors, &gin.Error{Err: appErr.NewProductError(400, "Maximum 10 images are allowed", "", nil)})
		return
	}
	files := map[string]*multipart.File{}
	for _, header := range fileheaders {
		file, openErr := header.Open()
		if openErr != nil {
			p.logger.Error().Err(openErr).Msg("false when opening file")
			c.Errors = append(c.Errors, &gin.Error{Err: appErr.NewProductError(500, "failed when opening file", "", nil)})
			return
		}
		files[header.Filename] = &file
	}

	thumbnailFile := req.MultipartForm.File["thumbnail"]
	if len(thumbnailFile) == 0 {
		c.Errors = append(c.Errors, &gin.Error{Err: appErr.NewProductError(400, "Thumbnail is required", "", nil)})
		return
	}

	thumbnailReq := thumbnailFile[0]

	nProduct, createProductErr := p.productService.CreateProduct(name, description, sku, userManual, files, thumbnailReq)
	if createProductErr != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: createProductErr})
		return
	}
	c.Header("Location", fmt.Sprintf("%s/%s", c.Request.URL.Path, nProduct.Id.String()))
	c.Status(201)
}

// UpdateProduct godoc
//
//	@Summary	Update product
//	@Tags		products
//	@Accept		json
//	@Param		id				path		string											true	"product id, Eg: ddb1fdef-2ffb-44a5-a833-fab7b4d60355 "
//	@Param		updateProduct	body		productModel.UpdateProduct						true	"product update body"
//	@Failure	400				{object}	model.ApiResponse{data=appErr.ValidationError}	"Cannot found param Id"
//	@Failure	400				{object}	model.ApiResponse{data=appErr.ValidationError}	"Cannot parse Id"
//	@Failure	500				{object}	model.ApiResponse{data=appErr.ValidationError}	"Cannot parse request body"
//	@Router		/products/{id} [PUT]
func (p *ProductHandler) updateProduct(c *gin.Context) {
	var updateProduct productModel.UpdateProduct
	paramId, ok := c.Params.Get("id")
	if !ok {
		c.Errors = append(c.Errors, &gin.Error{Err: appErr.NewProductError(400, "Cannot found param Id", "Cannot found param Id", nil)})
		return
	}

	id, parseIdErr := uuid.Parse(paramId)
	if parseIdErr != nil {

		c.Errors = append(c.Errors, &gin.Error{Err: appErr.NewProductError(400, "Cannot parse Id", "Cannot parse Id", nil)})
		return
	}
	bindJsonErr := c.BindJSON(&updateProduct)

	if bindJsonErr != nil {
		errMsg := "Cannot parse request body"
		c.Errors = append(c.Errors, &gin.Error{Err: appErr.NewProductError(500, errMsg, errMsg, nil)})
		return
	}

	updatedProduct, updateErr := p.productService.UpdateProduct(id, updateProduct.Name, updateProduct.Description, updateProduct.Sku, updateProduct.Status, updateProduct.UserManual)
	if updateErr != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: updateErr})
		return
	}
	res := model.SuccessResponse(productModel.FromEntity(*updatedProduct))
	c.JSON(res.Code, res)
}

func (p *ProductHandler) getQuantity(c *gin.Context) {
	productId, uuidParseErr := uuid.Parse(c.Query("productId"))
	if uuidParseErr != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: appErr.NewProductError(400, "couldn't parse product id", "couldn't parse product id", nil)})
		return
	}
	sizeId, uuidParseErr := uuid.Parse(c.Query("sizeId"))
	if uuidParseErr != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: appErr.NewProductError(400, "couldn't parse size id", "couldn't parse size id", nil)})
		return
	}
	colorId, uuidParseErr := uuid.Parse(c.Query("colorId"))
	if uuidParseErr != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: appErr.NewProductError(400, "couldn't parse color id", "couldn't parse color id", nil)})
		return
	}

	quantity, getErr := p.productService.GetQuantity(productId, sizeId, colorId)
	if getErr != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: getErr})
		return
	}

	c.JSON(200, model.SuccessResponse(quantity))
}

func (p *ProductHandler) uploadProductColorImage(c *gin.Context) {

	productId, parseErr := uuid.Parse(c.Param("id"))
	if parseErr != nil {
		p.logger.Error().Err(parseErr).Msg("")
		c.Errors = append(c.Errors, &gin.Error{Err: appErr.NewProductError(400, "couldn't parse product id", "couldn't parse product id", nil)})
		return
	}
	colorId, parseErr := uuid.Parse(c.Request.FormValue("colorId"))
	if parseErr != nil {
		p.logger.Error().Err(parseErr).Msg("")
		c.Errors = append(c.Errors, &gin.Error{Err: appErr.NewProductError(400, "couldn't parse color id", "couldn't parse color id", nil)})
		return
	}
	file := c.Request.MultipartForm.File["image"]
	if len(file) == 0 || len(file) > 1 {
		c.Errors = append(c.Errors, &gin.Error{Err: appErr.NewProductError(400, "required 1 image file", "required 1 image file", nil)})
		return
	}

	p.productService.UploadProductColor(productId, colorId, file[0])
}

func (p *ProductHandler) addProductInventory(c *gin.Context) {
	var createProductInventory productModel.CreateProductInventory
	productId, parseErr := uuid.Parse(c.Param("id"))
	if parseErr != nil {
		p.logger.Error().Err(parseErr).Msg("")
		c.Errors = append(c.Errors, &gin.Error{Err: appErr.NewProductError(400, "couldn't parse product id", "couldn't parse product id", nil)})
		return
	}
	bindErr := c.BindJSON(&createProductInventory)
	if bindErr != nil {
		p.logger.Error().Err(bindErr).Msg("")
		c.Errors = append(c.Errors, &gin.Error{Err: appErr.NewProductError(400, "couldn't parse body", "couldn't parse body", nil)})
		return
	}

	createInventories := make([]inventoryModel.CreateInventory, 0, len(createProductInventory.Inventories))
	for _, inventory := range createProductInventory.Inventories {
		if inventory.Price < 0 {
			c.Errors = append(c.Errors, &gin.Error{Err: appErr.NewProductError(400, "price cannot be a negative number", "price cannot be a negative number", nil)})
			return
		}
		createInventories = append(createInventories, inventoryModel.CreateInventory{
			ProductId: productId,
			ColorId:   inventory.ColorId,
			SizeId:    inventory.SizeId,
			Price:     inventory.Price,
			Quantity:  inventory.Quantity,
		})
	}
	if createErr := p.inventoryService.AddInventories(createInventories); createErr != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: createErr})
		return
	}
	c.JSON(200, model.SuccessResponse("Inventories added"))
}

func (p *ProductHandler) updateProductInventory(c *gin.Context) {
	productId, parseErr := uuid.Parse(c.Param("id"))
	if parseErr != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: appErr.NewProductError(400, "couldn't parse id", "couldn't parse id", nil)})
		return
	}
	var updateProductInventory productModel.UpdateProductInventoryRequest
	bindErr := c.BindJSON(&updateProductInventory)
	if bindErr != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: appErr.NewProductError(400, "couldn't parse body", "couldn't parse body", nil)})
		return
	}
	updateErr := p.inventoryService.UpdateInventory(productId, updateProductInventory.ColorId, updateProductInventory.SizeId, updateProductInventory.Quantity, updateProductInventory.Price)
	if updateErr != nil {
		c.Errors = append(c.Errors, &gin.Error{Err: updateErr})
		return
	}
}
