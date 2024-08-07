package service

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/product/internal/models"
	"github.com/tanush-128/openzo_backend/product/internal/repository"
)

// InventoryService defines the interface for the inventory service.
type InventoryService interface {
	CreateTransaction(ctx *gin.Context, transaction *models.InventoryTransaction) (*models.InventoryTransaction, error)
	GetTransactionByID(ctx *gin.Context,id string) (*models.InventoryTransaction, error)
	UpdateTransaction(ctx *gin.Context,transaction *models.InventoryTransaction) (*models.InventoryTransaction, error)
	DeleteTransaction(ctx *gin.Context,id string) error
	GetAllTransactionsByProductID(ctx *gin.Context,productID string) ([]models.InventoryTransaction, error)
}

type inventoryService struct {
	repo repository.InventoryTransactionRepository
}

// NewInventoryService creates a new instance of InventoryService.
func NewInventoryService(repo repository.InventoryTransactionRepository) InventoryService {
	return &inventoryService{repo: repo}
}

// CreateTransaction creates a new inventory transaction.
func (s *inventoryService) CreateTransaction(ctx *gin.Context,transaction *models.InventoryTransaction) (*models.InventoryTransaction, error) {

	err := s.repo.Create(transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

// GetTransactionByID retrieves an inventory transaction by its ID.
func (s *inventoryService) GetTransactionByID(ctx *gin.Context,id string) (*models.InventoryTransaction, error) {
	return s.repo.GetByID(id)
}

// UpdateTransaction updates an existing inventory transaction.
func (s *inventoryService) UpdateTransaction(ctx *gin.Context,transaction *models.InventoryTransaction) (*models.InventoryTransaction, error) {
	transaction, err := s.repo.GetByID(transaction.ID)
	if err != nil {
		return nil, err
	}
	if transaction == nil {
		return nil, errors.New("transaction not found")
	}

	err = s.repo.Update(transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

// DeleteTransaction deletes an inventory transaction by its ID.
func (s *inventoryService) DeleteTransaction(ctx *gin.Context,id string) error {
	return s.repo.Delete(id)
}

// GetAllTransactionsByProductID retrieves all inventory transactions for a given product ID.
func (s *inventoryService) GetAllTransactionsByProductID(ctx *gin.Context,productID string) ([]models.InventoryTransaction, error) {
	return s.repo.GetAllByProductID(productID)
}
