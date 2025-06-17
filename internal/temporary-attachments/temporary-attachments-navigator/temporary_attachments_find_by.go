package temporary_attachments_navigator

import (
	"github.com/google/uuid"
	generic_repo "word-meaning-finder/generics/generic-repo"
	"word-meaning-finder/internal/temporary-attachments/model"
)

func FindByIdService(id uuid.UUID) model.TemporaryAttachments {
	attachment, err := generic_repo.FindSingleByField[model.TemporaryAttachments]("id", id)
	if err != nil {
		panic("Didn't find attachment with that id")
	}

	return *attachment
}
