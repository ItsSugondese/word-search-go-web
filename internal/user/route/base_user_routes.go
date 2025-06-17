package route

import (
	"word-meaning-finder/internal/user/controller"
	authentication_middleware "word-meaning-finder/pkg/middleware/authentication-middleware"
	paseto_token "word-meaning-finder/pkg/utils/paseto-token"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func UserRoutes(r *gin.Engine, validate *validator.Validate) {
	users := r.Group("/user/")
	{
		users.POST("", func(c *gin.Context) {
			controller.RegisterUser(c, validate)
		})
		users.Use(authentication_middleware.PasetoAuthMiddleware(*paseto_token.TokenMaker))
		users.PUT("", func(c *gin.Context) {
			controller.RegisterUser(c, validate)
		})
		users.GET("doc/:id", controller.GetUserImage)
	}
}
