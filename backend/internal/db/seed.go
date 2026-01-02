package db

import (
	"github.com/sweetfish329/sabakan/backend/internal/auth"
	"github.com/sweetfish329/sabakan/backend/internal/models"
	"gorm.io/gorm"
)

// Seed populates the database with default roles and permissions.
func Seed() error {
	// Seed permissions
	for _, p := range models.DefaultPermissions() {
		result := DB.Where(models.Permission{
			Resource: p.Resource,
			Action:   p.Action,
		}).FirstOrCreate(&p)
		if result.Error != nil {
			return result.Error
		}
	}

	// Seed roles
	for _, r := range models.DefaultRoles() {
		result := DB.Where(models.Role{Name: r.Name}).FirstOrCreate(&r)
		if result.Error != nil {
			return result.Error
		}
	}

	// Assign permissions to roles
	if err := assignRolePermissions(); err != nil {
		return err
	}

	// Seed default admin user
	if err := seedDefaultAdmin(); err != nil {
		return err
	}

	return nil
}

// assignRolePermissions assigns default permissions to roles.
func assignRolePermissions() error {
	// Get the admin role and assign system:admin permission
	var adminRole models.Role
	if err := DB.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	var systemAdminPerm models.Permission
	if err := DB.Where("resource = ? AND action = ?", "system", "admin").First(&systemAdminPerm).Error; err != nil {
		return err
	}

	if err := DB.Model(&adminRole).Association("Permissions").Append(&systemAdminPerm); err != nil {
		return err
	}

	// Get the moderator role and assign game_server:* and mod:* permissions
	var moderatorRole models.Role
	if err := DB.Where("name = ?", "moderator").First(&moderatorRole).Error; err != nil {
		return err
	}

	var modPerms []models.Permission
	if err := DB.Where("resource IN ?", []string{"game_server", "mod"}).Find(&modPerms).Error; err != nil {
		return err
	}

	var userReadPerm models.Permission
	if err := DB.Where("resource = ? AND action = ?", "user", "read").First(&userReadPerm).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
	} else {
		modPerms = append(modPerms, userReadPerm)
	}

	if err := DB.Model(&moderatorRole).Association("Permissions").Replace(&modPerms); err != nil {
		return err
	}

	// Get the user role and assign read permissions
	var userRole models.Role
	if err := DB.Where("name = ?", "user").First(&userRole).Error; err != nil {
		return err
	}

	var readPerms []models.Permission
	if err := DB.Where("resource IN ? AND action = ?", []string{"game_server", "mod"}, "read").Find(&readPerms).Error; err != nil {
		return err
	}

	if err := DB.Model(&userRole).Association("Permissions").Replace(&readPerms); err != nil {
		return err
	}

	// Guest role has no permissions by default

	return nil
}

// seedDefaultAdmin creates the initial admin user if it doesn't exist.
func seedDefaultAdmin() error {
	// Check if admin user already exists
	var existingAdmin models.User
	if err := DB.Where("username = ?", "admin").First(&existingAdmin).Error; err == nil {
		// Admin already exists
		return nil
	}

	// Get admin role
	var adminRole models.Role
	if err := DB.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	// Hash default password
	// WARNING: Change this password on first login!
	passwordHash, err := auth.HashPassword("admin")
	if err != nil {
		return err
	}

	// Create admin user
	adminUser := models.User{
		Username:     "admin",
		PasswordHash: passwordHash,
		RoleID:       adminRole.ID,
		IsActive:     true,
	}

	if err := DB.Create(&adminUser).Error; err != nil {
		return err
	}

	return nil
}
