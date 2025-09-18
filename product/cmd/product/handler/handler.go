package handler

import (
	"fmt"
	"net/http"
	"product/cmd/product/usecase"
	"product/infrastructure/log"
	"product/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ProductHandler struct {
	ProductUseCase usecase.ProductUseCase
}

func NewProductHandler(productUseCase *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		ProductUseCase: *productUseCase,
	}
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	productId := c.Param("id")

	productIdConvert, err := strconv.ParseInt(productId, 10, 64)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"product ID": productId,
		}).Errorf("strconv.ParseInt got error : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Missing Param",
		})
		return
	}
	product, err := h.ProductUseCase.GetProductById(c.Request.Context(), productIdConvert)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"productId": productId,
		}).Errorf("h.ProductUseCase.GetProductById got error :%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_message": err,
		})
		return
	}
	if product.ID == 0 {
		log.Logger.WithFields(logrus.Fields{
			"productId": productId,
		}).Info("product ID not found")
		c.JSON(http.StatusOK, gin.H{
			"message": "Product Not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Succesfully get Product",
		"product": product,
	})
}

func (h *ProductHandler) GetProductCategory(c *gin.Context) {
	productCategoryId := c.Param("id")

	productCategoryIdConvert, err := strconv.Atoi(productCategoryId)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"product ID": productCategoryId,
		}).Errorf("strconv.Atoi got error : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Missing Param",
		})
		return
	}
	productCategory, err := h.ProductUseCase.GetProductCategoryById(c.Request.Context(), productCategoryIdConvert)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"productCategoryId": productCategoryId,
		}).Errorf("h.ProductUseCase.GetProductCategoryById got error :%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_message": "Internal Server Error",
		})
		return
	}
	if productCategory.ID == 0 {
		log.Logger.WithFields(logrus.Fields{
			"productCategoryId": productCategoryId,
		}).Info("product ID not found")
		c.JSON(http.StatusOK, gin.H{
			"message": "Product Not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Succesfully get Product",
		"product": productCategory,
	})
}

func (h *ProductHandler) ProductManagement(c *gin.Context) {
	var param models.ProductManagementParameter
	if err := c.ShouldBindBodyWithJSON(&param); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Input",
		})
		return
	}

	if param.Action == "" {
		log.Logger.Error("missing parameter action")
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Missing Required Parameter",
		})
		return
	}

	switch param.Action {
	case "add":
		if param.ID != 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid Request - product Id is not empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})
			return
		}
		productId, err := h.ProductUseCase.CreateProduct(c.Request.Context(), &param.Product)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUseCase.CreateProductCategory got error : %v ", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully create new product : %d", productId),
		})

	case "edit":
		if param.ID == 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid Request - product id is empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})
			return
		}
		product, err := h.ProductUseCase.UpdateProduct(c.Request.Context(), &param.Product)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUseCase.UpdateProduct got error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "successfull Edit Product",
			"product": product,
		})

	case "delete":
		if param.ID == 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid Request - product id is empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})
			return
		}
		err := h.ProductUseCase.DeleteProduct(c.Request.Context(), param.ID)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUseCase.DeleteProduct got error : %v ", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfull deleted product with id : %d", param.ID),
		})

	default:
		log.Logger.Error("invalid action")
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Action",
		})
		return
	}
}

func (h *ProductHandler) ProductCategoryManagement(c *gin.Context) {
	var param models.ProductCategoryManagementParameter
	if err := c.ShouldBindBodyWithJSON(&param); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Input",
		})
		return
	}

	if param.Action == "" {
		log.Logger.Error("missing parameter action")
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Missing Required Parameter",
		})
		return
	}

	switch param.Action {
	case "add":
		if param.ID != 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid Request - product category is not empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})
			return
		}
		productCategoryID, err := h.ProductUseCase.CreateProductCategory(c.Request.Context(), &param.ProductCategory)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUseCase.CreateProductCategory got error : %v ", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully create new product category: %d", productCategoryID),
		})

	case "edit":
		if param.ID == 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid Request - product id is empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})
			return
		}
		productCategory, err := h.ProductUseCase.UpdateProductCategory(c.Request.Context(), &param.ProductCategory)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUseCase.UpdateProductCategory got error : %v ", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":         "Success Edit Product",
			"productCategory": productCategory,
		})
		return
	case "delete":
		if param.ID == 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid Request - product id is empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})
			return
		}
		err := h.ProductUseCase.DeleteProductCategory(c.Request.Context(), param.ID)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUseCase.DeleteProductCategory got error : %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Product Category Id %d succesfully deleted", param.ID),
		})

	default:
		log.Logger.Error("invalid action")
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Action",
		})
		return
	}
}

func (h *ProductHandler) SearchProduct(c *gin.Context) {
	name := c.Query("name")
	category := c.Query("category")

	minPriceStr, _ := strconv.ParseFloat(c.Query("minPrice"), 64)
	maxPriceStr, _ := strconv.ParseFloat(c.Query("maxPrice"), 64)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	orderBy := c.Query("orderBy")
	sort := c.Query("sort")

	products, totalCount, err := h.ProductUseCase.SearchProduct(c.Request.Context(), models.SearchProductParameter{
		Name:     name,
		Category: category,
		MinPrice: int(minPriceStr),
		MaxPrice: int(maxPriceStr),
		Page:     page,
		PageSize: pageSize,
		OrderBy:  orderBy,
		Sort:     sort,
	})
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"param": c.Params,
		}).Errorf("h.ProductUseCase.SearchProduct got error : %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_message": err,
		})
		return
	}

	totalPages := (totalCount + pageSize - 1) / pageSize

	var nextPageUrl *string
	if page < totalPages {
		url := fmt.Sprintf("/products/search?name=%s&category=%s&minPrice=%f&maxPrice=%f&page=%d&pageSize=%d",
			c.Request.Host, name, category, minPriceStr, maxPriceStr, page+1, pageSize)
		nextPageUrl = &url
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully get products",
		"data": models.SearchProductResponse{
			Products:    products,
			Page:        page,
			PageSize:    pageSize,
			TotalCount:  totalCount,
			TotalPages:  totalPages,
			NextPageUrl: nextPageUrl,
		},
	})

}
