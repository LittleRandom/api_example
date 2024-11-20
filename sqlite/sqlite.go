package sqlite

import (
	"log"
	"path/filepath"
	"plainrandom/models"

	"gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

func OpenDatabase(path string) (*gorm.DB, error) {

	// Open DB command
	DB, err := gorm.Open(sqlite.Open(path), &gorm.Config{})

	if err != nil {
		return DB, err
	}

	// Run Automigration
	DB.AutoMigrate(&models.Item{})

	return DB, err
}

func NewDB() *gorm.DB {

	// Use stdlib to open a connection to postgres db.
	path := filepath.Join(filepath.Dir("./"), "api.db")
	DB, err := OpenDatabase(path)
	if err != nil {
		log.Printf("error when connecting to database %v", err)
		panic(err)
	}
	return DB
}
