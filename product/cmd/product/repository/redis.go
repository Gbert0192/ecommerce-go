package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"product/models"
	"time"
)

var (
	cacheKeyProductInfo         = "product:%d"
	cacheKeyProductCategoryInfo = "product_category:%d"
)

func (r *ProductRepository) GetProductByIDFromRedis(ctx context.Context, productID int64) (*models.Product, error) {

	cacheKey := fmt.Sprintf(cacheKeyProductInfo, productID)
	var product models.Product
	productStr, err := r.Redis.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(productStr), &product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetProductCategoryByIDFromRedis(ctx context.Context, productCategoryID int64) (*models.ProductCategory, error) {
	cacheKey := fmt.Sprintf(cacheKeyProductCategoryInfo, productCategoryID)
	var productCategory models.ProductCategory
	productCategoryStr, err := r.Redis.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(productCategoryStr), &productCategory)
	if err != nil {
		return nil, err
	}
	return &productCategory, nil
}

func (r *ProductRepository) SetProductByID(ctx context.Context, product *models.Product, productID int64) error {
	cacheKey := fmt.Sprintf(cacheKeyProductInfo, productID)
	productJSON, err := json.Marshal(product)
	if err != nil {
		return err
	}
	err = r.Redis.SetEx(ctx, cacheKey, productJSON, 10*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) SetProductCategoryByID(ctx context.Context, product *models.ProductCategory, productCategoryID int64) error {
	cacheKey := fmt.Sprintf(cacheKeyProductCategoryInfo, productCategoryID)
	productJSON, err := json.Marshal(product)
	if err != nil {
		return err
	}
	err = r.Redis.SetEx(ctx, cacheKey, productJSON, 1*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}
