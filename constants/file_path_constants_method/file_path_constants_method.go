package file_path_constants_method

import (
	"os"
	"path/filepath"
	filepathconstants "word-meaning-finder/constants/file_path_constants"
)

func GetAllWordPath() string {
	return filepath.Join(GetJavaResPath(), filepathconstants.FileSeparator, "german", filepathconstants.FileSeparator, "all-word")

}

func GetJavaResPath() string {
	return os.Getenv("JAVA_PROJECT_RESOURCE_PATH")
}
