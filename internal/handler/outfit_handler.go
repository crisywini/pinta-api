package handler

import (
	"net/http"

	"github.com/crisywini/pinta-api/internal/model"
	"github.com/crisywini/pinta-api/internal/service"
	"github.com/gin-gonic/gin"
)

type OutfitHandler struct {
	service *service.OutfitService
}

func NewOutfitHandler(service *service.OutfitService) *OutfitHandler {
	return &OutfitHandler{
		service: service,
	}
}

func (h *OutfitHandler) PostOutfit(c *gin.Context) {

	var outfit *model.Outfit

	if err := c.ShouldBindBodyWithJSON(&outfit); err == nil {
		created, warns, serviceError := h.service.Create(outfit)

		if serviceError != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error while creating outfit",
				"error":   serviceError.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"outfit":   created.ID.Hex(),
			"warnings": warns,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Malformed outfit body",
			"error":   err.Error(),
		})
	}
}

func (h *OutfitHandler) GetOutfitById(c *gin.Context) {
	id := c.Param("id")
	if outfit, err := h.service.GetById(id); err == nil {
		c.JSON(http.StatusOK, outfit)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error while getting the outfit",
			"error":   err.Error(),
		})
	}
}

func (h *OutfitHandler) GetAllOutfits(c *gin.Context) {
	outfits, err := h.service.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while getting the outfits",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"count": len(outfits),
		"items": outfits,
	})
}

func (h *OutfitHandler) DeleteOutfitById(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteById(id); err == nil {
		c.JSON(http.StatusNoContent, gin.H{})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while deleting outfit",
			"error":   err.Error(),
		})
	}
}

func (h *OutfitHandler) PutOutfit(c *gin.Context) {
	id := c.Param("id")
	var updated *model.Outfit

	if err := c.ShouldBindBodyWithJSON(&updated); err == nil {
		warns, serviceError := h.service.Update(id, updated)

		if serviceError == nil {
			c.JSON(http.StatusCreated, gin.H{
				"outfit":   id,
				"warnings": warns,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error while updating outfit",
				"error":   serviceError.Error(),
			})
		}
	}
}
