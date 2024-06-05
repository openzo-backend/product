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
	product.CrticalQuantity = utils.StringToInt(ctx.PostForm("critical_quantity"))
	product.Type = ctx.PostForm("type")
	product.VegType = ctx.PostForm("veg_type")
	product.Servers = utils.StringToInt(ctx.PostForm("servers"))
	product.CustomCode = ctx.PostForm("custom_code")
	product.SizeVariants = []models.SizeVariant{}
	product.ColorVariants = []models.ColorVariant{}

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
	product.ID = ctx.Param("id")
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
	product.CrticalQuantity = utils.StringToInt(ctx.PostForm("critical_quantity"))
	product.CustomCode = ctx.PostForm("custom_code")
	product.Type = ctx.PostForm("type")
	product.VegType = ctx.PostForm("veg_type")
	product.Servers = utils.StringToInt(ctx.PostForm("servers"))
	product.SizeVariants = []models.SizeVariant{}
	product.ColorVariants = []models.ColorVariant{}
	product.Images = []models.ProductImage{}

	json.Unmarshal([]byte(ctx.PostForm("size_variants")), &product.SizeVariants)
	json.Unmarshal([]byte(ctx.PostForm("color_variants")), &product.ColorVariants)

	updatedProduct, err := h.ProductService.UpdateProduct(ctx, product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedProduct)
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
