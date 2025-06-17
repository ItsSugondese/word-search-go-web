package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	localization_enums "word-meaning-finder/enums/struct-enums/localization-enums"
	"word-meaning-finder/enums/struct-enums/project_module"
	generic_repo "word-meaning-finder/generics/generic-repo"
	"word-meaning-finder/global/global_var"
	model2 "word-meaning-finder/internal/role/model"
	role_navigator "word-meaning-finder/internal/role/role-navigator"
	temporary_attachments_navigator "word-meaning-finder/internal/temporary-attachments/temporary-attachments-navigator"
	"word-meaning-finder/internal/user/dto"
	"word-meaning-finder/internal/user/model"
	user_navigator "word-meaning-finder/internal/user/user-navigator"
	"word-meaning-finder/pkg/common/database"
	"word-meaning-finder/pkg/common/localization"
	"word-meaning-finder/pkg/utils"
	dto_utils "word-meaning-finder/pkg/utils/dto-utils"
)

func RegisterBaseUserService(ctx *gin.Context, dto dto.UserRequest) /* *model.BaseUser */ any {
	tx := database.DB.Begin()
	tx.WithContext(ctx)

	exists := user_navigator.CheckUserByEmailExistValidationService(dto.Email)

	if exists {
		panic(localization.GetLocalizedMessage(localization_enums.MessageCodeEnums.COLUMN_ALREADY_EXISTS, map[string]interface{}{
			"First":  project_module.ModuleNameEnums.BASE_USER,
			"Second": "email",
		}))
	}
	var userDetails model.BaseUser
	var err error

	jsonData, _ := json.Marshal(dto)
	jsonUnmarshalError := json.Unmarshal(jsonData, &userDetails)

	if jsonUnmarshalError != nil {
		panic(jsonUnmarshalError)
	}

	if dto.Password != nil {
		hashedPassword, err := utils.HashPassword(*dto.Password)
		if err != nil {
			panic(err)
		}
		userDetails.Password = hashedPassword

	} else {
		panic("Passsword is must")
	}

	if dto.FileId != uuid.Nil {
		fileAttachment := temporary_attachments_navigator.FindByIdService(dto.FileId)
		userDetails.ProfilePath = utils.Ptr(utils.CopyFileToServer(fileAttachment.Location, project_module.ModuleNameEnums.BASE_USER, global_var.ForBucket))
	}

	getRole := role_navigator.FindRoleByIdService(dto.UserType)

	roles := []model2.Role{getRole}
	userDetails.Roles = roles

	userDetails, err = generic_repo.SaveRepo(tx, userDetails)
	if err != nil {
		panic(err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		panic(err)
	}

	return nil
}

func UpdateBaseUserService(ctx *gin.Context, dto dto.UserRequest) model.BaseUser {
	tx := database.DB.Begin()
	tx.WithContext(ctx)

	baseUserDetails := user_navigator.FindUserByIdService(dto.ID)

	// Convert the incoming DTO to Base user model
	dto_utils.DtoConvertErrorHandled(dto, &baseUserDetails)

	// Handle Base user image if provided
	if dto.FileId != uuid.Nil {
		baseUserPic := temporary_attachments_navigator.FindByIdService(dto.FileId)
		baseUserFilePath := utils.CopyFileToServer(baseUserPic.Location, project_module.ModuleNameEnums.BASE_USER, global_var.ForBucket)
		baseUserDetails.ProfilePath = &baseUserFilePath
	}

	if dto.Password != nil {
		hashedPassword, err := utils.HashPassword(*dto.Password)
		if err != nil {
			panic(err)
		}
		baseUserDetails.Password = hashedPassword
	}
	savedBaseUser, saveBaseUserError := generic_repo.SaveRepo(tx, baseUserDetails)

	// Handle any save or update errors
	if saveBaseUserError != nil {
		tx.Rollback()
		panic(saveBaseUserError)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		panic(err)
	}

	return savedBaseUser // Successfully return the saved Base user
}

func GetUserImageService(id uuid.UUID, w http.ResponseWriter) {
	userDetails := user_navigator.FindUserByIdService(id)

	if userDetails.ProfilePath != nil {

		utils.GetFileFromFilePath(*userDetails.ProfilePath, w, global_var.ForBucket)
	} else {
		panic("No user profile found")
	}
}
