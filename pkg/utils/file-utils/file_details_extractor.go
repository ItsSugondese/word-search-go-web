package file_utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	filepathconstants "word-meaning-finder/constants/file_path_constants"
)

func GetFileNamesInPathFromDirectory(dir string) ([]string, error) {
	var fileNames []string

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, dir+filepathconstants.FileSeparator+file.Name())
		}
	}

	return fileNames, nil
}

func GetAllFromFileAsSlices(path string) []string {
	var lines []string

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return []string{}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return []string{}
	}

	return lines
}

func GetCombinedLinesFromFilesParallel(filePaths []string) []string {
	var wg sync.WaitGroup
	lineChan := make(chan []string, len(filePaths))

	for _, path := range filePaths {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			lines := GetAllFromFileAsSlices(p)
			lineChan <- lines
		}(path)
	}

	wg.Wait()
	close(lineChan)

	var allLines []string
	for lines := range lineChan {
		allLines = append(allLines, lines...)
	}

	return allLines
}
