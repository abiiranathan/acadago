package database

import (
	"log"
	"os"

	"github.com/abiiranathan/acada/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var tables = []interface{}{
	&models.User{},
	&models.Program{},
	&models.CourseUnit{},
	&models.Resource{},
}

func MustConnect(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       os.Getenv("DSN"),
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("connected to mysql database")

	return db
}

func MigrateAll(db *gorm.DB) {
	if os.Getenv("MIGRATE") != "1" {
		return
	}

	err := db.AutoMigrate(tables...)
	if err != nil {
		log.Fatalf("migrations failed: %s", err)
	}

	log.Println("Migrations done! All tables created.")
}
