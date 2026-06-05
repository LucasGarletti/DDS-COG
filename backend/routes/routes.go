package routes

import (
	"net/http"

	"backend/controllers"
	"backend/dao"
	"backend/middlewares"
	"backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	userDAO := dao.NewUserDAO(db)
	authService := services.NewAuthService(userDAO)
	authController := controllers.NewAuthController(authService)

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/login", authController.Login)
		authRoutes.GET("/me", middlewares.AuthMiddleware(), authController.Me)
	}

	return router
}
