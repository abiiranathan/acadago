package services

import (
	"github.com/abiiranathan/acada/models"
	"gorm.io/gorm"
)

type ResourceService interface {
	GetAllResources() ([]models.Resource, error)
	GetResource(id uint) (models.Resource, error)
	CreateResource(r *models.Resource) error
	DeleteResource(id uint) error
}

type resourceunitservice struct {
	DB *gorm.DB
}

func NewResourceService(db *gorm.DB) ResourceService {
	return &resourceunitservice{DB: db}
}

func (svc *resourceunitservice) GetAllResources() ([]models.Resource, error) {
	resources := []models.Resource{}
	err := svc.DB.Find(&resources).Error
	return resources, err
}

func (svc *resourceunitservice) GetResource(id uint) (models.Resource, error) {
	resource := models.Resource{}
	err := svc.DB.First(&resource, id).Error
	return resource, err
}

func (svc *resourceunitservice) CreateResource(r *models.Resource) error {
	return svc.DB.Create(r).Error
}

func (svc *resourceunitservice) DeleteResource(id uint) error {
	return svc.DB.Delete(&models.Resource{}, id).Error
}
