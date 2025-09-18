package models

type Product struct {
	ID          int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Stock       int64  `json:"stock"`
	CategoryId  int64  `json:"category_id"`
}

type ProductCategory struct {
	ID   int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`
}

type ProductCategoryManagementParameter struct {
	Action string `json:"action"`
	ProductCategory
}

type ProductManagementParameter struct {
	Action string `json:"action"`
	Product
}

type SearchProductParameter struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	MinPrice int    `json:"minPrice"`
	MaxPrice int    `json:"maxPrice"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	OrderBy  string `json:"orderBy"`
	Sort     string `json:"sort"`
}

type SearchProductResponse struct {
	Products    []Product `json:"products"`
	Page        int       `json:"page"`
	PageSize    int       `json:"pageSize"`
	TotalCount  int       `json:"totalCount"`
	TotalPages  int       `json:"totalPages"`
	NextPageUrl *string   `json:"nextPageUrl"`
}
