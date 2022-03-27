package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/abiiranathan/acada/models"
	"gorm.io/gorm"
)

var (
	op       = flag.String("op", "", "create admin operation")
	name     = flag.String("name", "", "Admin name")
	email    = flag.String("email", "", "email of admin")
	password = flag.String("password", "", "password of admin")
)

func CreateAdmin(db *gorm.DB, user *models.User) error {
	user.IsAdmin = true
	user.IsActive = true
	user.CreatedAt = time.Now()

	if err := db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func ParseFlags() {
	if *name == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *email == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *password == "" {
		flag.Usage()
		os.Exit(1)
	}

	if len(*password) < 8 {
		log.Fatal("Admin password must be at least 8 characters")
	}

	if len(*password) > 128 {
		log.Fatal("Admin password must be less than 128 characters")
	}

}

func IfCreateAdmin(db *gorm.DB) {
	// if second argument is present, create admin
	if *op == "create-admin" {
		ParseFlags()

		err := CreateAdmin(db, getUserFromFlags())
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Admin created successfully")
		os.Exit(0)
	}
}

func getUserFromFlags() *models.User {
	return &models.User{
		Name:     *name,
		Email:    *email,
		Password: *password,
	}
}
