package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/product/internal/models"
)

func (s *productService) GetProductByID(ctx *gin.Context, id string) (models.Product, error) {
	Product, err := s.ProductRepository.GetProductByID(id)
	if err != nil {
		return models.Product{}, err
	}

	return Product, nil
}


func (s *productService) GetProductsByStoreID(ctx *gin.Context, storeID string) ([]models.Product, error) {
	Products, err := s.ProductRepository.GetProductsByStoreID(storeID)
	if err != nil {
		return []models.Product{}, err
	}

	return Products, nil
}