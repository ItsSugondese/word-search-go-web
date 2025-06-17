package globaldto

import "word-meaning-finder/constants/file_type_constants"

type FileDetails struct {
	FilePath string
	Size     int64
	FileType file_type_constants.FileType
}
