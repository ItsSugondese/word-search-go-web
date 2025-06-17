package main

import (
	"log"
	"net/http"
	"word-meaning-finder/config"
	socket_config "word-meaning-finder/config/socket-config"
	"word-meaning-finder/internal/auth/route"
	roleRoute "word-meaning-finder/internal/role/route"
	tempAttachmentRoute "word-meaning-finder/internal/temporary-attachments/route"
	"word-meaning-finder/internal/user"
	baseUserRoute "word-meaning-finder/internal/user/route"

	_ "word-meaning-finder/docs"
	global_gin_context "word-meaning-finder/global/global-gin-context"
	global_validation "word-meaning-finder/global/global-validation"
	"word-meaning-finder/pkg/common/localization"
	"word-meaning-finder/pkg/middleware"
	cors_middleware "word-meaning-finder/pkg/middleware/cors-middleware"
	lang_middleware "word-meaning-finder/pkg/middleware/lang-middleware"
	paseto_token "word-meaning-finder/pkg/utils/paseto-token"

	// "cloud.google.com/go/storage"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const ( // FILL IN WITH YOURS
	bucketName = "blackpearlbucket" // FILL IN WITH YOURS
)

// init method runs before the main method so that the environment variables are loaded before the application starts
func init() {
	log.Println("Loading environment variables and database connection")
	// load .env
	config.LoadEnvVariables()

	// load database connection
	//database.ConnectToDB()

	// load oAuth2 server
	//oauth2_setup.SetUpOAuth2()

	// paseto setup
	setupPaseto()

	// global gin hanler setup
	global_gin_context.NewGlobalGinContext()

	// lang a
	localization.InitLocalizationManager()

	// google json locaiton
	//os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./continual-mind-432410-g5-df1fc7f32718.json") // FILL IN WITH YOUR FILE PATH

	// Register the audit log callbacks and perform migrations
	//errVal := audit_middleware.RegisterCallbacks(database.DB)
	//if errVal != nil {
	//	panic("failed to register audit log callbacks")
	//}
	//
	//database.DB.AutoMigrate(&baseUserModel.BaseUser{}, &roleModel.Role{}, &model.TemporaryAttachments{})
	//
	//database.InitializeValuesInDb()

}
func main() {
	hub := socket_config.NewHub()
	go hub.Run()

	log.Println("Starting the application")
	r := gin.Default()
	validate := validator.New()

	// Initialize the Google Cloud Storage client using Gin's context.
	//r.Use(func(c *gin.Context) {
	//	client, err := storage.NewClient(c.Request.Context())
	//	if err != nil {
	//		log.Fatalf("Failed to create storage client: %v", err)
	//	}
	//
	//	// Store the client in Gin's context for use in handlers.
	//	c.Set("storageClient", client)
	//
	//	utils.Uploader = &utils.ClientUploader{
	//		Cl:         client,
	//		BucketName: bucketName,
	//		//ProjectID:  projectID,
	//		UploadPath: "test-files/",
	//	}
	//	// Make sure to close the client after the request is processed.
	//	defer client.Close()
	//
	//	c.Next()
	//})

	// middlewares
	r.Use(cors_middleware.CorsMiddleware())

	r.Use(middleware.RecoveryMiddleware())
	r.Use(lang_middleware.LocalizationMiddleware(localization.InitBundle()))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/ws", func(c *gin.Context) {
		socket_config.ServeWs(c, hub, c.Writer, c.Request)
	})

	r.Static("/static", "./static")

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// payload validations
	payloadValidations()

	// Registering routes
	route.AuthRoutes(r, validate)
	tempAttachmentRoute.TempAttachmentsRoutes(r, validate)
	baseUserRoute.UserRoutes(r, validate)
	roleRoute.RoleRoutes(r, validate)

	// Run the seed function to ensure default positions and permissions are set
	//seed.SeedDefaultPositions(database.DB)
	log.Println("_____________")
	// Serve static files from the images directory

	r.Run()
}

func setupPaseto() {
	tokenMaker, err := paseto_token.NewPaseto("abcdefghijkl12345678901234567890")
	if err != nil {
		panic("Couldnt open tokenmaker " + err.Error())
	}

	paseto_token.TokenMaker = tokenMaker
}

func payloadValidations() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// for user Type validation
		v.RegisterValidation("validUserType", user.ValidUserType)
		// for when both id and fileId is nil where id is uuid
		v.RegisterValidation("FileValidationIfIdNil", global_validation.RequiredIfIdNil)
		// for when both id and fileId is nil
		v.RegisterValidation("FieldValidationIfIdIsNil", global_validation.RequiredIfIdNilNotUUID)

	}
}
