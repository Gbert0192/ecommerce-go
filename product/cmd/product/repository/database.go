package repository

import (
	"context"
	"fmt"
	"product/models"
)

func (r *ProductRepository) GetProductByID(ctx context.Context, id int64) (*models.Product, error) {
	var product models.Product
	row := r.Database.WithContext(ctx).First(&product, id)
	if row.Error != nil {
		return nil, row.Error
	}
	return &product, nil
}

func (r *ProductRepository) FindProductId(ctx context.Context, productId int64) (*models.Product, error) {
	var product models.Product
	err := r.Database.WithContext(ctx).Table("product").Where("id = ?", productId).Last(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) FindProductCategory(ctx context.Context, productCategoryID int) (*models.ProductCategory, error) {
	var productCategory models.ProductCategory
	err := r.Database.WithContext(ctx).Table("product_category").Where("id = ?", productCategoryID).Last(&productCategory).Error
	if err != nil {
		return nil, err
	}
	return &productCategory, nil
}

func (r *ProductRepository) InsertNewProduct(ctx context.Context, product *models.Product) (int64, error) {
	err := r.Database.WithContext(ctx).Table("product").Create(product).Error
	if err != nil {
		return 0, err
	}
	return product.ID, nil
}

func (r *ProductRepository) InsertNewProductCategory(ctx context.Context, productCategory *models.ProductCategory) (int64, error) {
	err := r.Database.WithContext(ctx).Table("product_category").Create(productCategory).Error
	if err != nil {
		return 0, nil
	}
	return productCategory.ID, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	err := r.Database.WithContext(ctx).Table("product").Save(product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *ProductRepository) UpdateProductCategory(ctx context.Context, productCategory *models.ProductCategory) (*models.ProductCategory, error) {
	err := r.Database.WithContext(ctx).Table("product_category").Save(productCategory).Error
	if err != nil {
		return nil, err
	}
	return productCategory, nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, productId int64) error {
	err := r.Database.WithContext(ctx).Table("product").Delete(&models.Product{}, productId).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) DeleteProductCategory(ctx context.Context, productCategoryId int64) error {
	err := r.Database.WithContext(ctx).Table("product_category").Delete(&models.ProductCategory{}, productCategoryId).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ProductRepository) SearchProduct(ctx context.Context, params models.SearchProductParameter) ([]models.Product, int, error) {
	var products []models.Product
	var totalCount int64

	query := r.Database.WithContext(ctx).Table("product").
		Select("product.id, product.name, product.description, product.price, product.stock, product.category_id, product_category.name As category").
		Joins("JOIN product_category ON product.category_id = product_category.id")

	//FILTERING
	if params.Name != "" {
		query = query.Where("product.name ILIKE ?", "%"+params.Name+"%")
	}

	if params.Category != "" {
		query = query.Where("product_category.name = ?", params.Category)
	}
	if params.MinPrice > 0 {
		query = query.Where("product.price >= ?", params.MinPrice)
	}
	if params.MaxPrice > 0 {
		query = query.Where("product.price <= ?", params.MaxPrice)
	}

	//getting total count
	query.Model(&models.Product{}).Count(&totalCount)

	if params.OrderBy == "" {
		params.OrderBy = "product.name"
	}

	if params.Sort == "" || (params.Sort != "ASC" && params.Sort != "DESC") {
		params.Sort = "ASC"
	}
	orderBy := fmt.Sprintf("%s %s", params.OrderBy, params.Sort)
	query = query.Order(orderBy)

	offset := (params.Page - 1) * params.PageSize
	query = query.Limit(params.PageSize).Offset(offset)

	err := query.Scan(&products).Error
	if err != nil {
		return nil, 0, err

	}

	return products, int(totalCount), nil

}
