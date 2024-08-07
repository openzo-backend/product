package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/product/internal/models"
	"github.com/tanush-128/openzo_backend/product/internal/service"
)

type InventoryHandler struct {
	InventoryService service.InventoryService
}

func NewInventoryHandler(InventoryService *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{InventoryService: *InventoryService}
}

func (h *InventoryHandler) CreateInventoryTransaction(ctx *gin.Context) {
	var inventoryTransaction models.InventoryTransaction

	err := ctx.ShouldBindJSON(&inventoryTransaction)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdInventoryTransaction, err := h.InventoryService.CreateTransaction(ctx, &inventoryTransaction)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdInventoryTransaction)
}

func (h *InventoryHandler) GetInventoryTransactionByID(ctx *gin.Context) {
	id := ctx.Param("id")

	inventoryTransaction, err := h.InventoryService.GetTransactionByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if inventoryTransaction == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	ctx.JSON(http.StatusOK, inventoryTransaction)
}

func (h *InventoryHandler) UpdateInventoryTransaction(ctx *gin.Context) {
	var inventoryTransaction models.InventoryTransaction

	err := ctx.ShouldBindJSON(&inventoryTransaction)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedInventoryTransaction, err := h.InventoryService.UpdateTransaction(ctx, &inventoryTransaction)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedInventoryTransaction)
}

func (h *InventoryHandler) DeleteInventoryTransaction(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.InventoryService.DeleteTransaction(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Transaction deleted"})
}

func (h *InventoryHandler) GetAllTransactionsByProductID(ctx *gin.Context) {
	productID := ctx.Param("product_id")

	inventoryTransactions, err := h.InventoryService.GetAllTransactionsByProductID(ctx, productID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, inventoryTransactions)
}
