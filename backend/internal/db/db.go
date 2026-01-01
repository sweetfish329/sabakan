package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// DB is the global database connection.
var DB *gorm.DB

// Init initializes the database connection.
func Init(dataSourceName string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

// GetDB returns the database connection.
func GetDB() *gorm.DB {
	return DB
}

// AutoMigrate runs GORM auto-migration for all provided models.
func AutoMigrate(models ...any) error {
	return DB.AutoMigrate(models...)
}
