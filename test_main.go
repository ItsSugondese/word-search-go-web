package main

import (
	"fmt"
	"word-meaning-finder/constants/file_path_constants_method"
	file_utils "word-meaning-finder/pkg/utils/file_utils"
)

func testMain() {

	filePathSlices, _ := file_utils.GetFileNamesInPathFromDirectory(file_path_constants_method.GetAllWordPath())

	allWords := file_utils.GetCombinedLinesFromFilesParallel(filePathSlices)
	fmt.Println(len(allWords))
}
