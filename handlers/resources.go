package handlers

import (
	"github.com/abiiranathan/acada/models"
	"github.com/abiiranathan/acada/services"
	"github.com/gin-gonic/gin"
)

type ResourceHandler interface {
	GetAllResources() gin.HandlerFunc
	GetResource() gin.HandlerFunc
	CreateResource() gin.HandlerFunc
	DeleteResource() gin.HandlerFunc
}

type resourceHandler struct {
	ResourceService services.ResourceService
}

func NewResourceHandler(resourceService services.ResourceService) ResourceHandler {
	return &resourceHandler{ResourceService: resourceService}
}

func (h *resourceHandler) GetAllResources() gin.HandlerFunc {
	return func(c *gin.Context) {
		resources, err := h.ResourceService.GetAllResources()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, resources)
	}
}

func (h *resourceHandler) GetResource() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := GetParam(c, "id")
		if err != nil {
			return
		}

		resource, err := h.ResourceService.GetResource(id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, resource)
	}
}

func (h *resourceHandler) CreateResource() gin.HandlerFunc {
	return func(c *gin.Context) {
		resource := models.Resource{}
		if err := c.ShouldBindJSON(&resource); err != nil {
			RespondWithError(c, err)
			return
		}

		if err := h.ResourceService.CreateResource(&resource); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, resource)
	}
}

func (h *resourceHandler) DeleteResource() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := GetParam(c, "id")
		if err != nil {
			return
		}

		err = h.ResourceService.DeleteResource(id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{})
	}
}
