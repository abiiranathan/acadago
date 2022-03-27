package handlers

import (
	"github.com/abiiranathan/acada/models"
	"github.com/abiiranathan/acada/services"
	"github.com/gin-gonic/gin"
)

type ProgramHandler interface {
	GetAllPrograms() gin.HandlerFunc
	GetProgram() gin.HandlerFunc
	CreateProgram() gin.HandlerFunc
	UpdateProgram() gin.HandlerFunc
	DeleteProgram() gin.HandlerFunc
}

type programHandler struct {
	ProgramService services.ProgramService
}

func NewProgramHandler(programService services.ProgramService) ProgramHandler {
	return &programHandler{ProgramService: programService}
}

func (h *programHandler) GetAllPrograms() gin.HandlerFunc {
	return func(c *gin.Context) {
		programs, err := h.ProgramService.GetAllPrograms()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, programs)
	}
}

func (h *programHandler) GetProgram() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := GetParam(c, "id")
		if err != nil {
			return
		}

		program, err := h.ProgramService.GetProgram(id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, program)
	}
}

func (h *programHandler) CreateProgram() gin.HandlerFunc {
	return func(c *gin.Context) {
		program := models.Program{}
		if err := c.ShouldBindJSON(&program); err != nil {
			RespondWithError(c, err)
			return
		}

		if err := h.ProgramService.CreateProgram(&program); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		program.CourseUnits = []models.CourseUnit{}
		c.JSON(200, program)
	}
}

func (h *programHandler) UpdateProgram() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := GetParam(c, "id")
		if err != nil {
			return
		}

		program := models.Program{}
		if err := c.ShouldBindJSON(&program); err != nil {
			RespondWithError(c, err)
			return
		}

		program, err = h.ProgramService.UpdateProgram(id, program)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, program)
	}
}

func (h *programHandler) DeleteProgram() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := GetParam(c, "id")
		if err != nil {
			return
		}

		err = h.ProgramService.DeleteProgram(id)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Program deleted"})
	}
}
