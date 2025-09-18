package usecase

import (
	"context"
	"errors"
	"product/cmd/product/service"
	"product/infrastructure/log"
	"product/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductUseCase struct {
	ProductService service.ProductService
}

func NewProductUseCase(productService service.ProductService) *ProductUseCase {
	return &ProductUseCase{
		ProductService: productService,
	}
}

func (uc *ProductUseCase) GetProductById(ctx context.Context, productId int64) (*models.Product, error) {
	product, err := uc.ProductService.GetProductById(ctx, productId)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (uc *ProductUseCase) GetProductCategoryById(ctx context.Context, productCategoryId int) (*models.ProductCategory, error) {
	productCategory, err := uc.ProductService.ProductRepository.FindProductCategory(ctx, productCategoryId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &models.ProductCategory{}, nil
		}
		return nil, err
	}
	return productCategory, nil
}

func (uc *ProductUseCase) CreateProduct(ctx context.Context, param *models.Product) (int64, error) {
	productId, err := uc.ProductService.ProductRepository.InsertNewProduct(ctx, param)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"name":     param.Name,
			"category": param.CategoryId,
		}).Errorf("uc.ProductService.ProductRepository.InsertNewProduct got error : %v", err)
		return 0, err
	}
	return productId, nil
}

func (uc *ProductUseCase) CreateProductCategory(ctx context.Context, param *models.ProductCategory) (int64, error) {
	productCategoryId, err := uc.ProductService.ProductRepository.InsertNewProductCategory(ctx, param)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"name": param.Name,
		}).Errorf("uc.ProductService.ProductRepository.InsertNewProductCategory got error : %v", err)
		return 0, err
	}
	return productCategoryId, nil
}

func (uc *ProductUseCase) UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	product, err := uc.ProductService.ProductRepository.UpdateProduct(ctx, product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (uc *ProductUseCase) UpdateProductCategory(ctx context.Context, productCategory *models.ProductCategory) (*models.ProductCategory, error) {
	productCategory, err := uc.ProductService.ProductRepository.UpdateProductCategory(ctx, productCategory)
	if err != nil {
		return nil, err
	}
	return productCategory, nil
}

func (uc *ProductUseCase) DeleteProduct(ctx context.Context, productId int64) error {
	err := uc.ProductService.ProductRepository.DeleteProduct(ctx, productId)
	if err != nil {
		return err
	}
	return nil
}

func (uc *ProductUseCase) DeleteProductCategory(ctx context.Context, productCategoryId int64) error {
	err := uc.ProductService.ProductRepository.DeleteProductCategory(ctx, productCategoryId)
	if err != nil {
		return err
	}
	return nil
}

func (uc *ProductUseCase) SearchProduct(ctx context.Context, param models.SearchProductParameter) ([]models.Product, int, error) {
	product, totalCount, err := uc.ProductService.SearchProduct(ctx, param)

	if err != nil {
		return nil, 0, err
	}
	return product, totalCount, nil
}
