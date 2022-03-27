package services

import (
	"github.com/abiiranathan/acada/models"
	"gorm.io/gorm"
)

type UserService interface {
	GetAllUsers() ([]models.User, error)
	CreateUser(user *models.User) error
	GetUser(id uint) (models.User, error)
	UpdateUser(id uint, u models.User) (models.User, error)
	DeleteUser(id uint) error
	GetByEmail(email string) (models.User, error)
}

type userservice struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userservice{DB: db}
}

func (svc *userservice) GetAllUsers() ([]models.User, error) {
	users := []models.User{}
	err := svc.DB.Find(&users).Error
	return users, err
}

func (svc *userservice) GetUser(id uint) (models.User, error) {
	user := models.User{}
	err := svc.DB.First(&user, id).Error
	return user, err
}

func (svc *userservice) UpdateUser(id uint, u models.User) (models.User, error) {
	user := models.User{}
	if err := svc.DB.First(&user, id).Error; err != nil {
		return user, err
	}

	err := svc.DB.Model(&user).Updates(u).Error
	return user, err
}

func (svc *userservice) DeleteUser(id uint) error {
	return svc.DB.Delete(&models.User{}, id).Error
}

func (svc *userservice) CreateUser(user *models.User) error {
	return svc.DB.Create(user).Error
}

func (svc *userservice) GetByEmail(email string) (models.User, error) {
	user := models.User{}
	err := svc.DB.Where("email = ?", email).First(&user).Error
	return user, err
}
