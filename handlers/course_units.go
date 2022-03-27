package handlers

import (
	"time"

	"github.com/abiiranathan/acada/models"
	"github.com/abiiranathan/acada/services"
	"github.com/gin-gonic/gin"
)

type CourseUnitHandler interface {
	GetAllCourseUnits() gin.HandlerFunc
	GetCourseUnit() gin.HandlerFunc
	CreateCourseUnit() gin.HandlerFunc
	UpdateCourseUnit() gin.HandlerFunc
	DeleteCourseUnit() gin.HandlerFunc
}

type courseunitHandler struct {
	CourseUnitService services.CourseUnitService
}

func NewCourseUnitHandler(courseUnitService services.CourseUnitService) CourseUnitHandler {
	return &courseunitHandler{CourseUnitService: courseUnitService}
}

func (h *courseunitHandler) GetAllCourseUnits() gin.HandlerFunc {
	return func(c *gin.Context) {
		cu, err := h.CourseUnitService.GetAllCourseUnits()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, cu)
	}
}

func (h *courseunitHandler) GetCourseUnit() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := GetParam(c, "id")
		if err != nil {
			return
		}

		cu, err := h.CourseUnitService.GetCourseUnit(id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, cu)
	}
}

func (h *courseunitHandler) CreateCourseUnit() gin.HandlerFunc {
	return func(c *gin.Context) {
		cu := models.CourseUnit{}
		if err := c.ShouldBindJSON(&cu); err != nil {
			RespondWithError(c, err)
			return
		}

		// Add the current user as the creator of the course unit
		user := c.MustGet("user").(models.User)
		cu.CreatedBy = &user
		cu.CreatedAt = time.Now()

		if err := h.CourseUnitService.CreateCourseUnit(&cu); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		cu.Resources = []models.Resource{}

		c.JSON(200, cu)
	}
}

func (h *courseunitHandler) UpdateCourseUnit() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := GetParam(c, "id")
		if err != nil {
			return
		}

		response := models.CourseUnitUpdates{}
		if err := c.ShouldBindJSON(&response); err != nil {
			RespondWithError(c, err)
			return
		}

		updated := models.CourseUnit{
			Name:          response.Name,
			Code:          response.Code,
			Description:   response.Description,
			Semester:      response.Semester,
			Year:          response.Year,
			Instructor:    response.Instructor,
			EnrollmentKey: response.EnrollmentKey,
		}

		cu, err := h.CourseUnitService.UpdateCourseUnit(id, updated)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, cu)
	}
}

func (h *courseunitHandler) DeleteCourseUnit() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := GetParam(c, "id")
		if err != nil {
			return
		}

		err = h.CourseUnitService.DeleteCourseUnit(id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, "course unit deleted successfully")
	}
}
