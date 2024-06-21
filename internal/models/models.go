package models

type Product struct {
	ID      string `json:"id" gorm:"primaryKey"`
	StoreID string `json:"store_id"`

	Name          string         `json:"name"`
	Description   string         `json:"description"`
	QuantityUnit  string         `json:"quantity_unit" default:"Peice"`
	MRP           int            `json:"mrp"`
	DiscountPrice int            `json:"discount_price" gorm:"default:0"`
	Images        []ProductImage `json:"images" gorm:"foreignKey:ProductID;references:ID"`
	Brand         string         `json:"brand"`
	Barcode       string         `json:"barcode"`

	Category      string         `json:"category"`
	SizeVariants  []SizeVariant  `json:"size_variants"`
	ColorVariants []ColorVariant `json:"color_variants"`

	Type string `json:"type,omitempty"`

	VegType    string `json:"veg_type,omitempty"`
	Servers    int    `json:"servers,omitempty"`
	OutOfStock bool   `json:"out_of_stock" gorm:"default:false"`
	ProductPrivate
}

type ProductPrivate struct {
	MSRP            int    `json:"msrp,,omitempty"`
	Quantity        int    `json:"quantity,omitempty"`
	CrticalQuantity int    `json:"critical_quantity,omitempty"`
	CustomCode      string `json:"custom_code,omitempty"`
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
