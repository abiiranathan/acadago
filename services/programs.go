package services

import (
	"github.com/abiiranathan/acada/models"
	"gorm.io/gorm"
)

type ProgramService interface {
	GetAllPrograms() ([]models.Program, error)
	GetProgram(id uint) (models.Program, error)
	CreateProgram(program *models.Program) error
	UpdateProgram(id uint, p models.Program) (models.Program, error)
	DeleteProgram(id uint) error
}

type programservice struct {
	DB *gorm.DB
}

func NewProgramService(db *gorm.DB) ProgramService {
	return &programservice{DB: db}
}

func (svc *programservice) GetAllPrograms() ([]models.Program, error) {
	programs := []models.Program{}
	err := svc.DB.Preload("CourseUnits.Resources").Preload("CourseUnits.CreatedBy").Find(&programs).Error
	return programs, err
}

func (svc *programservice) GetProgram(id uint) (models.Program, error) {
	program := models.Program{}
	err := svc.DB.Preload("CourseUnits.Resources").Preload("CourseUnits.CreatedBy").First(&program, id).Error
	return program, err
}

func (svc *programservice) CreateProgram(program *models.Program) error {
	return svc.DB.Create(program).Error
}

func (svc *programservice) UpdateProgram(id uint, p models.Program) (models.Program, error) {
	program := models.Program{}
	if err := svc.DB.Preload("CourseUnits.Resources").Preload("CourseUnits.CreatedBy").First(&program, id).Error; err != nil {
		return program, err
	}

	err := svc.DB.Model(&program).Updates(p).Error
	return program, err
}

func (svc *programservice) DeleteProgram(id uint) error {
	return svc.DB.Delete(&models.Program{}, id).Error
}
