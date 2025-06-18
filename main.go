package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"word-meaning-finder/config"
	"word-meaning-finder/constants/file_path_constants_method"
	file_utils "word-meaning-finder/pkg/utils/file-utils"
)

// init method runs before the main method so that the environment variables are loaded before the application starts
func init() {
	log.Println("Loading environment variables and database connection")
	// load .env
	config.LoadEnvVariables()

}
func main() {

	log.Println("Starting the application")
	r := gin.Default()

	r.Static("/static", "./static")

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.GET("/words", func(c *gin.Context) {
		filePathSlices, _ := file_utils.GetFileNamesInPathFromDirectory(file_path_constants_method.GetAllWordPath())
		allWords := file_utils.GetCombinedLinesFromFilesParallel(filePathSlices)

		c.JSON(http.StatusOK, allWords)
	})

	log.Println("_____________")

	projectPath := os.Getenv("JAVA_PROJECT_RESOURCE_PATH")
	files, err := os.ReadDir(projectPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	var fileList []string
	for _, file := range files {
		if !file.IsDir() {
			fileList = append(fileList, file.Name())
		}
	}

	fmt.Println("Files in directory:")
	for _, name := range fileList {
		fmt.Println(name)
	}
	r.Run()
}
