package service

import (
	generic_repo "word-meaning-finder/generics/generic-repo"
	"word-meaning-finder/internal/role/dto"
	"word-meaning-finder/internal/role/model"
	role_navigator "word-meaning-finder/internal/role/role-navigator"
	"word-meaning-finder/pkg/common/database"
)

func CreateRoleService(dto *dto.RoleRequest) *model.Role {
	tx := database.DB.Begin()

	exists := role_navigator.CheckRoleExistValidationService(dto.Name)

	if exists {
		panic("Role already exists")
	}

	savedRole, saveRoleError := generic_repo.SaveRepo(tx, &model.Role{ID: dto.Name})
	if saveRoleError != nil {
		panic(saveRoleError)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		panic(err)
	}
	return savedRole
}
