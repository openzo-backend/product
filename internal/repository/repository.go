package repository

import (
	"github.com/google/uuid"
	"github.com/tanush-128/openzo_backend/product/internal/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(Product models.Product) (models.Product, error)
	GetProductByID(id string) (models.Product, error)
	GetProductsByStoreID(id string) ([]models.Product, error)
	UpdateProduct(Product models.Product) (models.Product, error)
	ChangeProductQuantity(id string, quantity int) error
	DeleteProduct(id string) error
	// Add more methods for other Product operations (GetProductByEmail, UpdateProduct, etc.)

}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {

	return &productRepository{db: db}
}

func (r *productRepository) CreateProduct(Product models.Product) (models.Product, error) {
	Product.ID = uuid.New().String()
	tx := r.db.Create(&Product)

	if tx.Error != nil {
		return models.Product{}, tx.Error
	}

	return Product, nil
}

func (r *productRepository) GetProductByID(id string) (models.Product, error) {
	var Product models.Product
	tx := r.db.Preload("Images").Preload("SizeVariants").Preload("ColorVariants").Where("id = ?", id).First(&Product)
	if tx.Error != nil {
		return models.Product{}, tx.Error
	}

	return Product, nil
}

func (r *productRepository) GetProductsByStoreID(id string) ([]models.Product, error) {
	var Products []models.Product
	tx := r.db.Preload("Images").Preload("SizeVariants").Preload("ColorVariants").Where("store_id = ?", id).Find(&Products)
	if tx.Error != nil {
		return []models.Product{}, tx.Error
	}

	return Products, nil
}

func (r *productRepository) DeleteProduct(id string) error {
	var Product models.Product
	tx := r.db.Where("id = ?", id).Delete(&Product)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *productRepository) UpdateProduct(Product models.Product) (models.Product, error) {
	var product models.Product
	tx := r.db.Where("id = ?", Product.ID).First(&product)
	if tx.Error != nil {
		return models.Product{}, tx.Error
	}

	tx = r.db.Save(&Product)
	if tx.Error != nil {
		return models.Product{}, tx.Error
	}
	Product.Images = product.Images

	return Product, nil
}

func (r *productRepository) ChangeProductQuantity(id string, quantity int) error {
	var Product models.Product
	tx := r.db.Where("id = ?", id).First(&Product)
	if tx.Error != nil {
		return tx.Error
	}

	Product.Quantity = Product.Quantity - quantity
	tx = r.db.Save(&Product)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// Implement other repository methods (GetProductByID, GetProductByEmail, UpdateProduct, etc.) with proper error handling
