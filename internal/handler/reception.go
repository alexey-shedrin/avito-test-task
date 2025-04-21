package handler

import (
	"log"
	"strings"

	"github.com/alexey-shedrin/avito-test-task/internal/middleware"
	"github.com/alexey-shedrin/avito-test-task/internal/model/dto/request"
	"github.com/alexey-shedrin/avito-test-task/internal/model/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	InvalidPvzId       = "invalid pvz id"
	InvalidPvzIdOrType = "invalid pvz id or type"
)

type ReceptionService interface {
	CreateReception(reception *entity.Reception) (*entity.Reception, error)
	CreateProduct(product *entity.Product, pvzID uuid.UUID) (*entity.Product, error)
	DeleteLastProduct(pvzID uuid.UUID) error
	CloseLastReception(pvzID uuid.UUID) (*entity.Reception, error)
}

func (h *Handler) PostReceptions(c *gin.Context) {
	log.SetPrefix("handler.PostReceptions")

	middleware.Auth(entity.EmployeeRole)(c)
	if c.IsAborted() {
		return
	}

	var req request.Reception
	if err := c.ShouldBindJSON(&req); err != nil {
		if strings.Contains(err.Error(), "Field validation") {
			c.JSON(400, gin.H{"error": InvalidPvzId})

			return
		}

		c.JSON(400, gin.H{"error": err.Error()})

		return
	}

	reception := &entity.Reception{
		PvzId: req.PvzId,
	}

	reception, err := h.receptionService.CreateReception(reception)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})

		return
	}

	c.JSON(201, reception.ToResponse())
}

func (h *Handler) PostProducts(c *gin.Context) {
	log.SetPrefix("handler.PostProducts")

	middleware.Auth(entity.EmployeeRole)(c)
	if c.IsAborted() {
		return
	}

	var req request.CreateProduct
	if err := c.ShouldBindJSON(&req); err != nil {
		if strings.Contains(err.Error(), "Field validation") {
			c.JSON(400, gin.H{"error": InvalidPvzIdOrType})

			return
		}

		c.JSON(400, gin.H{"error": err.Error()})

		return
	}

	product := &entity.Product{
		Type: req.Type,
	}

	product, err := h.receptionService.CreateProduct(product, req.PvzId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})

		return
	}

	c.JSON(201, product.ToResponse())
}

func (h *Handler) PostPvzPvzIdDeleteLastProduct(c *gin.Context, pvzId uuid.UUID) {
	log.SetPrefix("handler.PostPvzPvzIdDeleteLastProduct")

	middleware.Auth(entity.EmployeeRole)(c)
	if c.IsAborted() {
		return
	}

	if pvzId == uuid.Nil {
		c.JSON(400, gin.H{"error": InvalidPvzId})

		return
	}

	err := h.receptionService.DeleteLastProduct(pvzId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.Status(200)
}

func (h *Handler) PostPvzPvzIdCloseLastReception(c *gin.Context, pvzId uuid.UUID) {
	log.SetPrefix("handler.PostPvzPvzIdCloseLastReception")

	middleware.Auth(entity.EmployeeRole)(c)
	if c.IsAborted() {
		return
	}

	if pvzId == uuid.Nil {
		c.JSON(400, gin.H{"error": InvalidPvzId})
		return
	}

	reception, err := h.receptionService.CloseLastReception(pvzId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, reception.ToResponse())
}
