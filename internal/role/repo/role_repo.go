package repo

import (
	"gorm.io/gorm"
	"word-meaning-finder/internal/role/model"
)

func SaveRoleRepo(tx *gorm.DB, role *model.Role) (*model.Role, error) {
	result := tx.Create(&role)
	return role, result.Error
}
