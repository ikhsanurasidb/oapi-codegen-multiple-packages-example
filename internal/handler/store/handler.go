package store

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	gen_store "github.com/oapi-codegen-multiple-packages-example/internal/gen/store"
	"github.com/oapi-codegen-multiple-packages-example/internal/service/store"
)

type Handler struct {
	service store.Service
}

func NewHandler(svc store.Service) *Handler {
	return &Handler{
		service: svc,
	}
}

func (h *Handler) GetInventory(c *gin.Context) {
	ctx := c.Request.Context()

	inventory, err := h.service.GetInventory(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get inventory"})
		return
	}

	c.JSON(http.StatusOK, inventory)
}

func (h *Handler) PlaceOrder(c *gin.Context) {
	var order gen_store.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order data"})
		return
	}

	ctx := c.Request.Context()
	createdOrder, err := h.service.CreateOrder(ctx, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdOrder)
}

func (h *Handler) DeleteOrder(c *gin.Context, orderID int64) {
	ctx := c.Request.Context()

	if err := h.service.DeleteOrder(ctx, orderID); err != nil {
		if errors.Is(err, errors.New("order not found")) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) GetOrderById(c *gin.Context, orderID int64) {
	ctx := c.Request.Context()

	order, err := h.service.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, errors.New("order not found")) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order"})
		return
	}

	c.JSON(http.StatusOK, order)
}
