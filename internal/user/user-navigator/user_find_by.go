package user_navigator

import (
	"github.com/google/uuid"
	localization_enums "word-meaning-finder/enums/struct-enums/localization-enums"
	"word-meaning-finder/enums/struct-enums/project_module"
	generic_repo "word-meaning-finder/generics/generic-repo"
	"word-meaning-finder/internal/user/model"
	"word-meaning-finder/internal/user/repo"
	"word-meaning-finder/pkg/common/localization"
)

func FindUserByIdService(id uuid.UUID) model.BaseUser {
	userDetails, err := generic_repo.FindSingleByField[model.BaseUser]("id", id)

	if err != nil {
		panic(err)
	}

	if userDetails == nil {
		panic(localization.GetLocalizedMessage(localization_enums.MessageCodeEnums.COLUMN_NOT_EXISTS, map[string]interface{}{
			"First":  project_module.ModuleNameEnums.BASE_USER,
			"Second": "Id",
		}))
	}
	return *userDetails
}

func FindUserByEmailService(email string) *model.BaseUser {
	userDetails, err := repo.FindUserByColumnRepo(email, "email")
	if err != nil {
		panic(err)
	}
	if userDetails == nil {
		panic(localization.GetLocalizedMessage(localization_enums.MessageCodeEnums.COLUMN_NOT_EXISTS, map[string]interface{}{
			"First":  project_module.ModuleNameEnums.BASE_USER,
			"Second": "email",
		}))
	}
	return userDetails
}

func CheckUserByEmailExistValidationService(email string) bool {
	userDetails, err := repo.FindUserByColumnRepo(email, "email")
	if err != nil {
		panic(err)
	}
	if userDetails == nil {
		return false
	}
	return true
}
