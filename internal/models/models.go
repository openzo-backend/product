package models

type Product struct {
	ID              string         `json:"id" gorm:"primaryKey"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	QuantityUnit    string         `json:"quantity_unit" default:"Peice"`
	MRP             int            `json:"mrp"`
	MSRP            int            `json:"msrp"`
	DiscountPrice   int            `json:"discount_price"`
	Quantity        int            `json:"quantity"`
	Images          []ProductImage `json:"images" gorm:"foreignKey:ProductID;references:ID"`
	Brand           string         `json:"brand"`
	CrticalQuantity int            `json:"critical_quantity"`
	CustomCode      string         `json:"custom_code"`
	SizeVariants    []SizeVariant  `json:"size_variants"`
	ColorVariants   []ColorVariant `json:"color_variants"`
	Category        string         `json:"category"`
	Barcode         string         `json:"barcode"`
	StoreID         string         `json:"store_id"`
}

type SizeVariant struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	ProductID string `json:"product_id" gorm:"size:36;index"`
	Size      string `json:"size"`
	Price     int    `json:"price"`
	Quantity  int    `json:"quantity"`
}

type ColorVariant struct {
	ID        int    `json:"id"`
	ProductID string `json:"product_id" gorm:"size:36;index"`
	Color     string `json:"color"`
	Price     int    `json:"price"`
	Quantity  int    `json:"quantity"`
}

type ProductImage struct {
	ID        int    `json:"id"`
	ProductID string `json:"product_id" gorm:"size:36;index"`
	Image     string `json:"image" gorm:"type:text"`
}
