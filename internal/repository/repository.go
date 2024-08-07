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
	GetPostByPincode(pincode string) ([]ProductWithStore, error)
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
		Preload("InventoryTransactions").
		Where("store_id = ?", storeID).
		Order("category ASC, display_order ASC").
		Find(&products)
	if tx.Error != nil {
		return []models.Product{}, tx.Error
	}

	for i, product := range products {
		var totalQuantity int
		for _, transaction := range product.InventoryTransactions {
			totalQuantity += transaction.Quantity
		}
		products[i].Quantity = totalQuantity
	}

	// 	type InventoryTransaction struct {
	// 	ID        string `json:"id" gorm:"primaryKey"`
	// 	ProductID string `json:"product_id" gorm:"size:36;index"`
	// 	Quantity  int    `json:"quantity" gorm:"not null"`
	// 	Price     int    `json:"price" gorm:"not null"`

	// 	// TransactionType can be one of the following:
	// 	// 1. INVENTORY_ADJUSTMENT
	// 	// 2. PURCHASE
	// 	// 3. SALE
	// 	// 4. RETURN
	// 	TransactionType string `json:"transaction_type" gorm:"not null"`
	// 	Description     string `json:"description" gorm:"type:text"`
	// 	CreatedAt       time.Time
	// }

	// var inventoryTransactions []models.InventoryTransaction
	// subQuery := r.db.Model(&models.Product{}).Select("id").Where("store_id = ?", storeID)
	// tx = r.db.Where("product_id IN (?)", subQuery).Find(&inventoryTransactions)

	// if tx.Error != nil {
	// 	return []models.Product{}, tx.Error
	// }

	// for i, product := range products {
	// 	for _, transaction := range inventoryTransactions {
	// 		if product.ID == transaction.ProductID {
	// 			fmt.Println("transaction.ProductID %+v", transaction)
	// 			products[i].Quantity += transaction.Quantity

	// 			products[i].InventoryTransactions = append(products[i].InventoryTransactions, transaction)
	// 		}
	// 	}

	// }

	return products, nil
}

type StoreBasicDetails struct {
	StoreId    string `json:"storee_id"`
	StoreName  string `json:"store_name"`
	StoreImage string `json:"store_image"`

	StoreAddress string `json:"store_address"`

	StoreCategory    string `json:"store_category" gorm:"default:general_store'"`
	StoreSubCategory string `json:"store_sub_category" gorm:"default:general_store'"`

	StoreDescription string  `json:"store_description"`
	StoreRating      float64 `json:"store_rating" gorm:"default:0"`
	StoreReviewCount int     `json:"store_review_count" gorm:"default:0"`
}

type ProductWithStore struct {
	models.Product
	StoreBasicDetails
}

func (r *productRepository) GetPostByPincode(pincode string) ([]ProductWithStore, error) {
	var products []ProductWithStore

	// product of type == post is a Post
	// product has field store_id which is the id of the store, store has field pincode
	// so we can join product and store on store_id and then filter by store.pincode

	tx := r.db.
		Model(&models.Product{}).
		Select("products.*, stores.id as storee_id, stores.name as store_name, stores.image as store_image, stores.address as store_address, stores.category as store_category, stores.sub_category as store_sub_category, stores.description as store_description, stores.rating as store_rating, stores.review_count as store_review_count").
		Preload("Images").
		Where("products.type = ?", "post").
		Joins("JOIN stores ON products.store_id = stores.id").
		Where("stores.pincode = ? ", pincode).
		Order("products.created_at DESC").
		Find(&products)

	if tx.Error != nil {
		return []ProductWithStore{}, tx.Error
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
