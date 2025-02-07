package product

import (
	"context"
	"database/sql"

	"mime/multipart"
	"post-backend/internal/custom"
	"post-backend/internal/helper"
	"post-backend/internal/notification"
	"post-backend/internal/setting"
	stockhistory "post-backend/internal/stock_history"
	"sync"

	"github.com/gosimple/slug"
)

type ProductService interface {
	Create(ctx context.Context, input CreateProductInput, productImagesFile map[string]multipart.File, userId int) (Product, error)
	Update(ctx context.Context, inputData UpdateProductInput, inputId GetProductInput) (Product, error)
	Delete(ctx context.Context, input GetProductInput) error
	Get(ctx context.Context, input GetProductInput) (Product, error)
	GetBySlug(ctx context.Context, input SlugProductInput) (Product, error)
	GetAll(ctx context.Context) ([]Product, error)
	CreateImage(ctx context.Context, input GetProductInput, productImagesFile map[string]multipart.File) error
	DeleteImage(ctx context.Context, input GetProductImageInput) error
	SetLogo(ctx context.Context, inputProductId GetProductInput, inputProductImageId GetProductImageInput) error
	UpdateStock(ctx context.Context, inputProductId GetProductInput, inputData UpdateStockProductInput, userId int) error
}

type ProductServiceImpl struct {
	DB                     *sql.DB
	ProductRepository      ProductRepository
	StockHistoryRepository stockhistory.StockHistoryRepository
	NotificationRepository notification.NotificationRepository
	SettingRepository      setting.SettingRepository
}

func NewProductService(DB *sql.DB, productRepository ProductRepository, stockHistoryRepository stockhistory.StockHistoryRepository, notificationRepository notification.NotificationRepository, settingRepository setting.SettingRepository) ProductService {
	return &ProductServiceImpl{
		DB:                     DB,
		ProductRepository:      productRepository,
		StockHistoryRepository: stockHistoryRepository,
		NotificationRepository: notificationRepository,
		SettingRepository:      settingRepository,
	}
}

func (p *ProductServiceImpl) Create(ctx context.Context, input CreateProductInput, productImagesFile map[string]multipart.File, userId int) (Product, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return Product{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	isProductExist, err := p.ProductRepository.FindBySlug(ctx, tx, slug.Make(input.Name))
	if err != nil && err != custom.ErrNotFound {
		return Product{}, err
	}
	if isProductExist.Id > 0 {
		return Product{}, custom.ErrAlreadyExists
	}

	product := Product{
		CategoryId:  input.CategoryId,
		Name:        input.Name,
		Slug:        slug.Make(input.Name),
		Price:       input.Price,
		Stock:       input.Stock,
		Description: input.Description,
		Status:      input.Status,
	}
	product, err = p.ProductRepository.Insert(ctx, tx, product)
	if err != nil {
		return Product{}, err
	}

	stockHistory := stockhistory.StockHistory{
		ProductId:   product.Id,
		Type:        "+",
		Qty:         product.Stock,
		StockBefore: 0,
		StockAfter:  product.Stock,
		UserId:      userId,
	}
	_, err = p.StockHistoryRepository.Insert(ctx, tx, stockHistory)
	if err != nil {
		return Product{}, err
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	isFirstFile := true

	errChan := make(chan error, len(productImagesFile))

	for fileName, file := range productImagesFile {
		wg.Add(1)
		go func(fileNameInner string, fileInner multipart.File) {
			defer wg.Done()
			productImage := ProductImages{
				ProductId: product.Id,
			}

			mu.Lock()
			currentIsFirst := isFirstFile
			if isFirstFile {
				isFirstFile = false
			}
			mu.Unlock()

			if currentIsFirst {
				productImage.IsLogo = true
			}

			insertProductImages(*p, "product", fileNameInner, fileInner, errChan, &productImage, ctx, tx)
		}(fileName, file)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return Product{}, err
		}
	}

	return product, nil
}

func (p *ProductServiceImpl) CreateImage(ctx context.Context, input GetProductInput, productImagesFile map[string]multipart.File) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.HandleTransaction(tx, &err)

	product, err := p.ProductRepository.FindById(ctx, tx, input.Id)
	if err != nil {
		return err
	}
	if product.Id <= 0 {
		return custom.ErrNotFound
	}

	var wg sync.WaitGroup
	errChan := make(chan error)
	productImage := ProductImages{
		ProductId: input.Id,
		IsLogo:    false,
	}

	for fileName, file := range productImagesFile {
		wg.Add(1)
		go func(fileNameInner string, fileInner multipart.File) {
			defer wg.Done()
			insertProductImages(*p, "product", fileNameInner, fileInner, errChan, &productImage, ctx, tx)
		}(fileName, file)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *ProductServiceImpl) Delete(ctx context.Context, input GetProductInput) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.HandleTransaction(tx, &err)

	product, err := p.ProductRepository.FindById(ctx, tx, input.Id)
	if err != nil {
		return err
	}
	if product.Id <= 0 {
		return custom.ErrNotFound
	}

	err = p.ProductRepository.Delete(ctx, tx, input.Id)
	if err != nil {
		return err
	}

	for _, image := range product.ProductImages {
		err = deleteImage("product", image.ImageUrl)

		if err != nil {
			return err
		}
	}

	return nil
}

func (p *ProductServiceImpl) DeleteImage(ctx context.Context, input GetProductImageInput) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.HandleTransaction(tx, &err)

	productImage, err := p.ProductRepository.FindImageById(ctx, tx, input.Id)
	if err != nil {
		return err
	}
	if productImage.Id <= 0 {
		return custom.ErrNotFound
	}

	err = p.ProductRepository.DeleteImage(ctx, tx, input.Id)
	if err != nil {
		return err
	}

	if productImage.IsLogo {
		lastestProductImage, err := p.ProductRepository.FindLatestImage(ctx, tx, productImage.ProductId)
		if err != nil {
			return err
		}

		err = p.ProductRepository.UpdateImage(ctx, tx, lastestProductImage.Id)
		if err != nil {
			return err
		}
	}

	err = deleteImage("product", productImage.ImageUrl)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductServiceImpl) Get(ctx context.Context, input GetProductInput) (Product, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return Product{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	product, err := p.ProductRepository.FindById(ctx, tx, input.Id)
	if err != nil {
		return Product{}, err
	}
	if product.Id <= 0 {
		return Product{}, custom.ErrNotFound
	}

	return product, nil
}

func (p *ProductServiceImpl) GetBySlug(ctx context.Context, input SlugProductInput) (Product, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return Product{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	product, err := p.ProductRepository.FindBySlug(ctx, tx, input.Slug)
	if err != nil {
		return Product{}, err
	}
	if product.Id <= 0 {
		return Product{}, custom.ErrNotFound
	}

	return product, nil
}

func (p *ProductServiceImpl) GetAll(ctx context.Context) ([]Product, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return []Product{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	products, err := p.ProductRepository.FindAll(ctx, tx)
	if err != nil {
		return []Product{}, err
	}
	if len(products) <= 0 {
		return []Product{}, nil
	}

	return products, nil
}

func (p *ProductServiceImpl) SetLogo(ctx context.Context, inputProductId GetProductInput, inputProductImageId GetProductImageInput) error {
	tx, err := p.DB.Begin()
	if err != nil {
		return err
	}
	defer helper.HandleTransaction(tx, &err)

	product, err := p.ProductRepository.FindById(ctx, tx, inputProductId.Id)
	if err != nil {
		return err
	}
	if product.Id <= 0 {
		return custom.ErrNotFound
	}

	productImages, err := p.ProductRepository.FindImageById(ctx, tx, inputProductImageId.Id)
	if err != nil {
		return err
	}
	if productImages.Id <= 0 {
		return custom.ErrNotFound
	}

	err = p.ProductRepository.UpdateImagesLogoFalse(ctx, tx, inputProductId.Id)
	if err != nil {
		return err
	}

	err = p.ProductRepository.UpdateImage(ctx, tx, inputProductImageId.Id)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductServiceImpl) Update(ctx context.Context, inputData UpdateProductInput, inputId GetProductInput) (Product, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return Product{}, err
	}
	defer helper.HandleTransaction(tx, &err)

	product, err := p.ProductRepository.FindById(ctx, tx, inputId.Id)
	if err != nil {
		return Product{}, err
	}
	if product.Id <= 0 {
		return Product{}, custom.ErrNotFound
	}

	isSlugExist, err := p.ProductRepository.FindBySlug(ctx, tx, slug.Make(inputData.Name))
	if err != nil && err != custom.ErrNotFound {
		return Product{}, err
	}
	if isSlugExist.Id > 0 && isSlugExist.Id != product.Id {
		return Product{}, custom.ErrAlreadyExists
	}

	product.CategoryId = inputData.CategoryId
	product.Name = inputData.Name
	product.Slug = slug.Make(inputData.Name)
	product.Price = inputData.Price
	product.Stock = inputData.Stock
	product.Description = inputData.Description
	product.Status = inputData.Status

	product, err = p.ProductRepository.Update(ctx, tx, product)
	if err != nil {
		return Product{}, err
	}

	err = CheckProductStock(p.SettingRepository, p.NotificationRepository, ctx, tx, product)
	if err != nil {
		return Product{}, err
	}

	return product, nil
}

func (p *ProductServiceImpl) UpdateStock(ctx context.Context, inputProductId GetProductInput, inputData UpdateStockProductInput, userId int) error {
	tx, err := p.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer helper.HandleTransaction(tx, &err)

	product, err := p.ProductRepository.FindById(ctx, tx, inputProductId.Id)
	if err != nil {
		return err
	}
	if product.Id < 0 {
		return custom.ErrNotFound
	}

	newStock := 0

	switch inputData.Type {
	case "+":
		newStock = product.Stock + inputData.Qty
	case "-":
		newStock = product.Stock - inputData.Qty
	default:
		return custom.ErrInternal
	}

	if newStock < 0 {
		return custom.ErrInsufficientStock
	}

	err = p.ProductRepository.UpdateStock(ctx, tx, inputProductId.Id, newStock)
	if err != nil {
		return err
	}
	product.Stock = newStock

	stockHistory := stockhistory.StockHistory{
		ProductId:   inputProductId.Id,
		Type:        inputData.Type,
		Qty:         inputData.Qty,
		StockBefore: product.Stock,
		StockAfter:  newStock,
		UserId:      userId,
	}
	_, err = p.StockHistoryRepository.Insert(ctx, tx, stockHistory)
	if err != nil {
		return err
	}

	err = CheckProductStock(p.SettingRepository, p.NotificationRepository, ctx, tx, product)
	if err != nil {
		return err
	}

	return nil
}
