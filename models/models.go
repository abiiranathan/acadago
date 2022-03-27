// Package models declares structs for all database tables
package models

import (
	"fmt"
	"time"

	"github.com/abiiranathan/acada/auth"
	"gorm.io/gorm"
)

type ResourceType uint

const (
	PastPaper ResourceType = iota + 1
	LectureNotes
	Assignment
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey;not null"`
	Name      string    `json:"name" gorm:"not null" binding:"required"`
	Email     string    `json:"email" gorm:"not null;unique" binding:"email,lowercase,required"`
	Password  string    `json:"-" gorm:"not null" binding:"required,min=8"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	IsAdmin   bool      `json:"is_admin" gorm:"default:false"`
}

type Login struct {
	Email    string `json:"email" gorm:"not null" binding:"required,email"`
	Password string `json:"password" gorm:"not null" binding:"required"`
}

type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

func hash_password(u *User, password string) error {
	hash, err := auth.HashPassword(password)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}

	u.Password = hash
	return nil
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	return hash_password(u, u.Password)
}

// Model for a single course unit
type CourseUnit struct {
	ID            uint   `json:"id" gorm:"primaryKey;not null"`
	Name          string `json:"name" gorm:"not null;index;type:varchar(200)" binding:"required"`
	Code          string `json:"code" gorm:"not null;unique;type:varchar(200)" binding:"required"`
	Description   string `json:"description" gorm:"not null;type:text" binding:"required"`
	Semester      string `json:"semester" gorm:"not null" binding:"required"`
	Year          string `json:"year" gorm:"not null;index;type:varchar(25)" binding:"required"`
	Instructor    string `json:"instructor" gorm:"not null" binding:"required"`
	EnrollmentKey string `json:"enrollment_key" gorm:"not null;type:varchar(100);" binding:"required"`

	// FK for Program
	ProgramID uint `json:"program_id" gorm:"not_null;"`

	// FK for User
	UserID uint `json:"user_id" gorm:"not null;index" binding:"required"`

	CreatedAt time.Time  `json:"created_at" gorm:"not null;type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP;"`
	CreatedBy *User      `json:"created_by" gorm:"not null;foreignKey:UserID;"`
	Resources []Resource `json:"resources" gorm:"courseunit_resources"`
}

// Model for a single course unit
type CourseUnitUpdates struct {
	Name          string `json:"name" gorm:"not null;index;type:varchar(200)"`
	Code          string `json:"code" gorm:"not null;unique;type:varchar(200)"`
	Description   string `json:"description" gorm:"not null;type:text"`
	Semester      string `json:"semester" gorm:"not null"`
	Year          string `json:"year" gorm:"not null;index;type:varchar(25)"`
	Instructor    string `json:"instructor" gorm:"not null"`
	EnrollmentKey string `json:"enrollment_key" gorm:"not null;type:varchar(100);"`
}

type Resource struct {
	ID           uint         `json:"id" gorm:"primaryKey;not null"`
	Name         string       `json:"name" gorm:"not null" binding:"required"`
	ResourceType ResourceType `json:"resource_type" gorm:"not null" binding:"required"`
	File         string       `json:"file" gorm:"not null" binding:"required,base64url"`
	Ext          string       `json:"ext" gorm:"not null" binding:"required"`
	CourseUnitID uint         `json:"course_unit_id" gorm:"not_null;"`
}

type Program struct {
	ID          uint         `json:"id" gorm:"primaryKey;not null"`
	Name        string       `json:"name" gorm:"not null;uniqueIndex" binding:"required"`
	CourseUnits []CourseUnit `json:"course_units" gorm:"program_course_units"`
}
