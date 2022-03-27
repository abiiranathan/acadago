package services

import (
	"github.com/abiiranathan/acada/models"
	"gorm.io/gorm"
)

type CourseUnitService interface {
	GetAllCourseUnits() ([]models.CourseUnit, error)
	GetCourseUnit(id uint) (models.CourseUnit, error)
	CreateCourseUnit(cu *models.CourseUnit) error
	UpdateCourseUnit(id uint, cu models.CourseUnit) (models.CourseUnit, error)
	DeleteCourseUnit(id uint) error
}

type courseunitservice struct {
	DB *gorm.DB
}

func NewCourseUnitService(db *gorm.DB) CourseUnitService {
	return &courseunitservice{DB: db}
}

func (svc *courseunitservice) GetAllCourseUnits() ([]models.CourseUnit, error) {
	cu := []models.CourseUnit{}
	err := svc.DB.Preload("Resources").Preload("CreatedBy").Find(&cu).Error
	return cu, err
}

func (svc *courseunitservice) GetCourseUnit(id uint) (models.CourseUnit, error) {
	cu := models.CourseUnit{}
	err := svc.DB.Preload("Resources").Preload("CreatedBy").First(&cu, id).Error
	return cu, err
}

func (svc *courseunitservice) CreateCourseUnit(cu *models.CourseUnit) error {
	return svc.DB.Create(cu).Error
}

func (svc *courseunitservice) UpdateCourseUnit(id uint, cu models.CourseUnit) (models.CourseUnit, error) {
	cu1 := models.CourseUnit{}
	if err := svc.DB.Preload("Resources").Preload("CreatedBy").First(&cu1, id).Error; err != nil {
		return cu1, err
	}

	err := svc.DB.Model(&cu1).Updates(cu).Error
	return cu1, err
}

func (svc *courseunitservice) DeleteCourseUnit(id uint) error {
	return svc.DB.Delete(&models.CourseUnit{}, id).Error
}
