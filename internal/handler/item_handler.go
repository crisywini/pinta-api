package handler

import (
	"net/http"

	"github.com/crisywini/pinta-api/internal/model"
	"github.com/crisywini/pinta-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ItemHandler struct {
	service *service.ItemService
}

func NewItemHandler(service *service.ItemService) *ItemHandler {
	return &ItemHandler{
		service: service,
	}
}

// Create a new item

func (h *ItemHandler) PostItem(c *gin.Context) {
	var item *model.Item

	if err := c.ShouldBindBodyWith(&item, binding.JSON); err == nil {
		response, serviceError := h.service.Create(item)

		if serviceError != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": serviceError.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id": response.ID.Hex(),
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Malformed item",
			"error":   err.Error(),
		})
	}
}

func (h *ItemHandler) GetItemByID(c *gin.Context) {

	itemId := c.Param("id")

	if len(itemId) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing item id",
		})
	}

	item, err := h.service.GetById(itemId)

	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	}

	c.JSON(http.StatusOK, item)
}

func (h *ItemHandler) GetAllItems(c *gin.Context) {
	allItems, _ := h.service.GetAll()

	c.JSON(http.StatusOK, gin.H{
		"count": len(allItems),
		"items": allItems,
	})
}

func (h *ItemHandler) PutItem(c *gin.Context) {

	id := c.Param("id")
	var updatedItem *model.Item
	if err := c.ShouldBindBodyWith(&updatedItem, binding.JSON); err == nil {
		serviceError := h.service.Update(id, updatedItem)
		if serviceError != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error while updating item",
				"error":   serviceError.Error(),
			})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Malformed request",
			"error":   err.Error(),
		})
	}
}

func (h *ItemHandler) DeleteItem(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteById(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error while deleting item",
			"error":   err.Error(),
		})
	} else {
		c.JSON(http.StatusNoContent, gin.H{})
	}

}
