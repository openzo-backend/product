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
	UpdateDisplayOrder(id string, displayOrder int) error
	ChangeProductQuantity(id string, quantity int) error
	BatchUpdateDisplayOrder(updates []models.Product) error
	DeleteProduct(id string) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {

	return &productRepository{db: db}
}

func (r *productRepository) CreateProduct(product models.Product) (models.Product, error) {
	product.ID = uuid.New().String()

	// Set default display order if not provided
	if product.DisplayOrder == 0 {
		var maxDisplayOrder int
		r.db.Model(&models.Product{}).
			Where("store_id = ?", product.StoreID).
			Select("COALESCE(MAX(display_order), 0)").
			Row().Scan(&maxDisplayOrder)
		product.DisplayOrder = maxDisplayOrder + 1
	}

	tx := r.db.Create(&product)
	if tx.Error != nil {
		return models.Product{}, tx.Error
	}

	return product, nil
}

func (r *productRepository) GetProductByID(id string) (models.Product, error) {
	var Product models.Product
	tx := r.db.Preload("Images").Preload("SizeVariants").Preload("ColorVariants").Where("id = ?", id).First(&Product)
	if tx.Error != nil {
		return models.Product{}, tx.Error
	}

	return Product, nil
}

func (r *productRepository) GetProductsByStoreID(storeID string) ([]models.Product, error) {
	var products []models.Product
	tx := r.db.Preload("Images").
		Preload("SizeVariants").
		Preload("ColorVariants").
		Where("store_id = ?", storeID).
		Order("category ASC, display_order ASC").
		Find(&products)
	if tx.Error != nil {
		return []models.Product{}, tx.Error
	}

	return products, nil
}

func (r *productRepository) DeleteProduct(id string) error {

	tx := r.db.Where("id = ?", id).Delete(&models.ProductImage{})
	if tx.Error != nil {
		return tx.Error
	}

	tx = r.db.Where("product_id = ?", id).Delete(&models.SizeVariant{})
	if tx.Error != nil {
		return tx.Error
	}

	tx = r.db.Where("product_id = ?", id).Delete(&models.ColorVariant{})
	if tx.Error != nil {
		return tx.Error
	}

	var product models.Product
	if err := r.db.Where("id = ?", id).First(&product).Error; err != nil {
		return err
	}

	tx = r.db.Where("id = ?", id).Delete(&models.Product{})
	if tx.Error != nil {
		return tx.Error
	}

	// Adjust the display order of remaining products
	r.db.Exec("UPDATE products SET display_order = display_order - 1 WHERE store_id = ? AND display_order > ?", product.StoreID, product.DisplayOrder)
	return nil
}

func (r *productRepository) UpdateProduct(Product models.Product) (models.Product, error) {

	// update all images

	tx := r.db.Model(&Product).Association("Images").Replace(Product.Images)
	if tx != nil {
		return models.Product{}, tx
	}

	tx = r.db.Model(&Product).Association("SizeVariants").Replace(Product.SizeVariants)
	if tx != nil {
		return models.Product{}, tx
	}

	tx = r.db.Model(&Product).Association("ColorVariants").Replace(Product.ColorVariants)
	if tx != nil {
		return models.Product{}, tx
	}

	tx1 := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&Product)
	if tx1.Error != nil {
		return models.Product{}, tx1.Error
	}

	return Product, nil
}

func (r *productRepository) UpdateDisplayOrder(id string, displayOrder int) error {
	tx := r.db.Model(&models.Product{}).
		Where("id = ?", id).
		Update("display_order", displayOrder)

	return tx.Error
}

func (r *productRepository) BatchUpdateDisplayOrder(updates []models.Product) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, product := range updates {
			if err := tx.Model(&models.Product{}).Where("id = ?", product.ID).Update("display_order", product.DisplayOrder).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *productRepository) ChangeProductQuantity(id string, quantity int) error {
	var Product models.Product
	tx := r.db.Where("id = ?", id).First(&Product)
	if tx.Error != nil {
		return tx.Error
	}

	Product.Quantity = quantity
	tx = r.db.Save(&Product)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// Implement other repository methods (GetProductByID, GetProductByEmail, UpdateProduct, etc.) with proper error handling
