package service

import (
	"context"
	"product/cmd/product/repository"
	"product/infrastructure/log"
	"product/models"

	"github.com/sirupsen/logrus"
)

type ProductService struct {
	ProductRepository repository.ProductRepository
}

func NewProductRepository(productRepository repository.ProductRepository) *ProductService {
	return &ProductService{
		ProductRepository: productRepository,
	}
}

func (s *ProductService) GetProductById(ctx context.Context, productId int64) (*models.Product, error) {
	//redis
	product, err := s.ProductRepository.GetProductByIDFromRedis(ctx, productId)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"productId": productId,
		}).Errorf("s.ProductRepository.GetProductByIDFromRedis got error : %v", err)
	}
	if product != nil && product.ID != 0 {
		return product, nil
	}

	//db
	product, err = s.ProductRepository.FindProductId(ctx, productId)
	if err != nil {
		return nil, err
	}

	// Create background context for cache update (won't be cancelled when request ends)
	go func(product *models.Product, productId int64) {
		log.Logger.WithFields(logrus.Fields{
			"productId": productId,
		}).Info("Starting cache update goroutine")

		// Use background context to prevent cancellation
		cacheCtx := context.Background()
		errConCurrent := s.ProductRepository.SetProductByID(cacheCtx, product, productId)
		if errConCurrent != nil {
			log.Logger.WithFields(logrus.Fields{
				"product":   product,
				"productId": productId,
			}).Errorf("s.ProductRepository.SetProductByID got error: %v", errConCurrent)
		} else {
			log.Logger.WithFields(logrus.Fields{
				"productId": productId,
			}).Info("Successfully cached product in Redis")
		}
	}(product, productId)
	return product, nil
}

func (s *ProductService) GetProductCategoryById(ctx context.Context, productCategoryId int) (*models.ProductCategory, error) {
	productCategory, err := s.ProductRepository.FindProductCategory(ctx, productCategoryId)
	if err != nil {
		return nil, err
	}
	return productCategory, nil
}

func (s *ProductService) CreateNewProduct(ctx context.Context, param *models.Product) (int64, error) {
	productId, err := s.ProductRepository.InsertNewProduct(ctx, param)
	if err != nil {
		return 0, err
	}
	return productId, nil
}

func (s *ProductService) CreateNewProductCategory(ctx context.Context, param *models.ProductCategory) (int64, error) {
	productCategoryId, err := s.ProductRepository.InsertNewProductCategory(ctx, param)
	if err != nil {
		return 0, nil
	}
	return productCategoryId, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	productData, err := s.ProductRepository.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}
	return productData, nil
}

func (s *ProductService) UpdateProductCategory(ctx context.Context, productCategory *models.ProductCategory) (*models.ProductCategory, error) {
	productCategoryData, err := s.ProductRepository.UpdateProductCategory(ctx, productCategory)
	if err != nil {
		return nil, err
	}
	return productCategoryData, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, productId int64) error {
	err := s.ProductRepository.DeleteProduct(ctx, productId)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductService) DeleteProductCategory(ctx context.Context, productId int64) error {
	err := s.ProductRepository.DeleteProductCategory(ctx, productId)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProductService) SearchProduct(ctx context.Context, param models.SearchProductParameter) ([]models.Product, int, error) {
	products, totalCount, err := s.ProductRepository.SearchProduct(ctx, param)
	if err != nil {
		return nil, 0, err
	}
	return products, totalCount, nil
}
