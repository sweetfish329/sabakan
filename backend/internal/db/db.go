package db

import (
	"github.com/glebarez/sqlite"
	"github.com/sweetfish329/sabakan/backend/internal/models"
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

// Migrate runs GORM auto-migration for all models.
func Migrate() error {
	return DB.AutoMigrate(
		&models.Role{},
		&models.Permission{},
		&models.User{},
		&models.OAuthAccount{},
		&models.APIToken{},
		&models.RefreshToken{},
		&models.GameServer{},
		&models.GameServerPort{},
		&models.GameServerEnv{},
		&models.Mod{},
		&models.GameServerMod{},
		&models.AuditLog{},
	)
}
