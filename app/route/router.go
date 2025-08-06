package route

import (
	"setup-preoject/app/controller"
	"setup-preoject/app/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, cache *redis.Client) {
	api := r.Group("/api")
	v1 := api.Group("/v1")

	// controller
	authController := controller.NewAuthController(db, cache)

	// setup route
	auth := v1.Group("/auth")
	{
		auth.POST("/login", authController.Login)
		auth.POST("/register", authController.Register)
		auth.GET("/me", middleware.AuthMiddleware(cache), authController.Me)
	}
}
