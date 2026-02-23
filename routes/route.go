package routes

import (
	"restfulapi/config"
	"restfulapi/controller"
	"restfulapi/middleware"
	"restfulapi/repository"
	"restfulapi/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	//buat gin router
	router := gin.Default()

	//dependency injection
	userRepo := repository.NewUserRepository(config.DB)
	authService := services.NewAuthService(userRepo)
	authController := controller.NewAuthController(authService)

	auth := router.Group("/api/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
	}

	//route yang diprotected
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/users", authController.GetAllUsers)
		protected.GET("/users/:id", authController.GetUserById)
		protected.PUT("/users/:id", authController.Update)
		protected.DELETE("/users/:id", authController.Delete)
	}
	return router
}
