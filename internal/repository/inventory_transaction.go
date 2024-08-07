package repository

import (
	"github.com/google/uuid"
	"github.com/tanush-128/openzo_backend/product/internal/models"
	"gorm.io/gorm"
)

// InventoryTransaction represents the inventory transaction model.

// InventoryTransactionRepository defines the interface for inventory transaction repository.
type InventoryTransactionRepository interface {
	Create(transaction *models.InventoryTransaction) error
	GetByID(id string) (*models.InventoryTransaction, error)
	Update(transaction *models.InventoryTransaction) error
	Delete(id string) error
	GetAllByProductID(productID string) ([]models.InventoryTransaction, error)
}

type inventoryTransactionRepository struct {
	db *gorm.DB
}

// NewInventoryTransactionRepository creates a new instance of InventoryTransactionRepository.
func NewInventoryTransactionRepository(db *gorm.DB) InventoryTransactionRepository {
	return &inventoryTransactionRepository{db: db}
}

// Create inserts a new inventory transaction into the database.
func (r *inventoryTransactionRepository) Create(transaction *models.InventoryTransaction) error {
	transaction.ID = uuid.New().String()
	return r.db.Create(transaction).Error
}

// GetByID retrieves an inventory transaction by its ID.
func (r *inventoryTransactionRepository) GetByID(id string) (*models.InventoryTransaction, error) {
	var transaction models.InventoryTransaction
	if err := r.db.First(&transaction, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

// Update modifies an existing inventory transaction.
func (r *inventoryTransactionRepository) Update(transaction *models.InventoryTransaction) error {
	return r.db.Save(transaction).Error
}

// Delete removes an inventory transaction by its ID.
func (r *inventoryTransactionRepository) Delete(id string) error {
	return r.db.Delete(&models.InventoryTransaction{}, "id = ?", id).Error
}

// GetAllByProductID retrieves all inventory transactions for a given product ID.
func (r *inventoryTransactionRepository) GetAllByProductID(productID string) ([]models.InventoryTransaction, error) {
	var transactions []models.InventoryTransaction
	if err := r.db.Where("product_id = ?", productID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
