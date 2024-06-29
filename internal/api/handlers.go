package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/product/internal/models"
	"github.com/tanush-128/openzo_backend/product/internal/service"
	"github.com/tanush-128/openzo_backend/product/internal/utils"
)

type Handler struct {
	ProductService service.ProductService
}

func NewHandler(ProductService *service.ProductService) *Handler {
	return &Handler{ProductService: *ProductService}
}

type ProductDisplayOrderUpdate struct {
	ProductID    string `json:"product_id"`
	DisplayOrder int    `json:"display_order"`
}

type BatchUpdateRequest struct {
	Updates []ProductDisplayOrderUpdate `json:"updates"`
}

func (h *Handler) CreateProduct(ctx *gin.Context) {
	var product models.Product

	product.Name = ctx.PostForm("name")
	product.Description = ctx.PostForm("description")
	product.QuantityUnit = ctx.PostForm("quantity_unit")
	product.MRP = utils.StringToInt(ctx.PostForm("mrp"))
	product.MSRP = utils.StringToInt(ctx.PostForm("msrp"))
	product.DiscountPrice = utils.StringToInt(ctx.PostForm("discount_price"))
	product.Barcode = ctx.PostForm("barcode")
	product.StoreID = ctx.PostForm("store_id")
	product.Category = ctx.PostForm("category")
	product.Quantity = utils.StringToInt(ctx.PostForm("quantity"))
	product.Brand = ctx.PostForm("brand")
	product.CriticalQuantity = utils.StringToInt(ctx.PostForm("critical_quantity"))
	product.Type = ctx.PostForm("type")
	product.VegType = ctx.PostForm("veg_type")
	product.Servers = utils.StringToInt(ctx.PostForm("servers"))
	product.CustomCode = ctx.PostForm("custom_code")
	product.SizeVariants = []models.SizeVariant{}
	product.ColorVariants = []models.ColorVariant{}
	product.OutOfStock = ctx.PostForm("out_of_stock") == "true"

	json.Unmarshal([]byte(ctx.PostForm("size_variants")), &product.SizeVariants)
	json.Unmarshal([]byte(ctx.PostForm("color_variants")), &product.ColorVariants)

	createdProduct, err := h.ProductService.CreateProduct(ctx, product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdProduct)

}

func (h *Handler) GetProductByID(ctx *gin.Context) {
	id := ctx.Param("id")

	Product, err := h.ProductService.GetProductByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Product)
}

func (h *Handler) GetProductsByStoreID(ctx *gin.Context) {
	storeID := ctx.Param("id")

	Products, err := h.ProductService.GetProductsByStoreID(ctx, storeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, Products)

}

func (h *Handler) UpdateProduct(ctx *gin.Context) {

	var product models.Product
	product.ID = ctx.PostForm("id")
	product.Name = ctx.PostForm("name")
	product.Description = ctx.PostForm("description")
	product.QuantityUnit = ctx.PostForm("quantity_unit")
	product.MRP = utils.StringToInt(ctx.PostForm("mrp"))
	product.MSRP = utils.StringToInt(ctx.PostForm("msrp"))
	product.DiscountPrice = utils.StringToInt(ctx.PostForm("discount_price"))
	product.Barcode = ctx.PostForm("barcode")
	product.StoreID = ctx.PostForm("store_id")
	product.Category = ctx.PostForm("category")
	product.Quantity = utils.StringToInt(ctx.PostForm("quantity"))
	product.Brand = ctx.PostForm("brand")
	product.CriticalQuantity = utils.StringToInt(ctx.PostForm("critical_quantity"))
	product.CustomCode = ctx.PostForm("custom_code")
	product.Type = ctx.PostForm("type")
	product.VegType = ctx.PostForm("veg_type")
	product.Servers = utils.StringToInt(ctx.PostForm("servers"))
	product.SizeVariants = []models.SizeVariant{}
	product.ColorVariants = []models.ColorVariant{}
	product.Images = []models.ProductImage{}

	product.OutOfStock = ctx.PostForm("out_of_stock") == "true"


	json.Unmarshal([]byte(ctx.PostForm("product_images")), &product.Images)
	json.Unmarshal([]byte(ctx.PostForm("size_variants")), &product.SizeVariants)
	json.Unmarshal([]byte(ctx.PostForm("color_variants")), &product.ColorVariants)

	updatedProduct, err := h.ProductService.UpdateProduct(ctx, product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedProduct)
}

func (h *Handler) UpdateDisplayOrder(ctx *gin.Context) {
	id := ctx.Param("id")
	displayOrder := utils.StringToInt(ctx.Query("display_order"))

	err := h.ProductService.UpdateDisplayOrder(ctx, id, displayOrder)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product display order updated successfully"})
}

func (h *Handler) BatchUpdateDisplayOrder(c *gin.Context) {
	var req BatchUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products := make([]models.Product, len(req.Updates))
	for i, update := range req.Updates {
		products[i] = models.Product{
			ID:           update.ProductID,
			DisplayOrder: update.DisplayOrder,
		}
	}

	if err := h.ProductService.BatchUpdateDisplayOrder(c, products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product display orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *Handler) ChangeProductQuantity(ctx *gin.Context) {
	id := ctx.Param("id")
	quantity := utils.StringToInt(ctx.Query("quantity"))

	err := h.ProductService.ChangeProductQuantity(ctx, id, quantity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product quantity updated successfully"})
}

func (h *Handler) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.ProductService.DeleteProduct(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
