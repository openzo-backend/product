package models

import "time"

type Product struct {
	ID              string         `json:"id" gorm:"primaryKey`
	StoreID         string         `json:"store_id" gorm:"size:36;not null"`
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
	Name            string         `json:"name" gorm:"not null"`
	Description     string         `json:"description" gorm:"type:text"`
	QuantityUnit    string         `json:"quantity_unit" gorm:"default:'Piece';not null"`
	MRP             int            `json:"mrp" gorm:"not null"`
	DiscountPrice   int            `json:"discount_price" gorm:"default:0"`
	Images          []ProductImage `json:"images"`
	Brand           string         `json:"brand"`
	Barcode         string         `json:"barcode" gorm:"index;size:36"`
	Category        string         `json:"category"`
	DisplayOrder    int            `json:"display_order" gorm:"default:0"`
	SizeVariants    []SizeVariant  `json:"size_variants" `
	ColorVariants   []ColorVariant `json:"color_variants"`
	Type            string         `json:"type,omitempty"`
	MetaDescription string         `json:"meta_description,omitempty"`
	MetaTags        string         `json:"meta_tags,omitempty"`
	VegType         string         `json:"veg_type,omitempty"`
	Servers         int            `json:"servers,omitempty"`
	OutOfStock      bool           `json:"out_of_stock" gorm:"default:false"`
	ProductPrivate
}

type InventoryTransaction struct {
	ID        string `json:"id" gorm:"primaryKey"`
	ProductID string `json:"product_id" gorm:"size:36;index"`
	Quantity  int    `json:"quantity" gorm:"not null"`
	Price     int    `json:"price" gorm:"not null"`

	// TransactionType can be one of the following:
	// 1. INVENTORY_ADJUSTMENT
	// 2. PURCHASE
	// 3. SALE
	// 4. RETURN
	TransactionType string    `json:"transaction_type" gorm:"not null"`
	Description     string    `json:"description" gorm:"type:text"`
	CreatedAt       time.Time `json:"created_at"`
}

type ProductPrivate struct {
	MSRP                  int                    `json:"msrp,omitempty"`
	Quantity              int                    `json:"quantity,omitempty"`
	CriticalQuantity      int                    `json:"critical_quantity,omitempty"`
	CustomCode            string                 `json:"custom_code,omitempty"`
	InventoryTransactions []InventoryTransaction `json:"inventory_transactions,omitempty"`
}

type SizeVariant struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID string `json:"product_id" gorm:"size:36;index"`
	Size      string `json:"size" gorm:"not null"`
	Price     int    `json:"price" gorm:"not null"`
	Quantity  int    `json:"quantity" gorm:"not null"`
}

type ColorVariant struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID string `json:"product_id" gorm:"size:36;index"`
	Color     string `json:"color" gorm:"not null"`
	Price     int    `json:"price" gorm:"not null"`
	Quantity  int    `json:"quantity" gorm:"not null"`
}

type ProductImage struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID string `json:"product_id" gorm:"size:36;index"`
	Image     string `json:"image" gorm:"type:text"`
}
