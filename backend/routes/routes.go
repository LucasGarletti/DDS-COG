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
	router.Use(corsMiddleware())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	userDAO := dao.NewUserDAO(db)
	authService := services.NewAuthService(userDAO)
	authController := controllers.NewAuthController(authService)

	eventDAO := dao.NewEventDAO(db)
	eventService := services.NewEventService(eventDAO)
	eventController := controllers.NewEventController(eventService)

	ticketDAO := dao.NewTicketDAO(db)
	ticketService := services.NewTicketService(ticketDAO, userDAO)
	ticketController := controllers.NewTicketController(ticketService)

	router.GET("/eventos", eventController.GetAll)
	router.GET("/eventos/:id", eventController.GetByID)
	router.GET("/mis-entradas", middlewares.AuthMiddleware(), ticketController.GetMyTickets)

	ticketRoutes := router.Group("/entradas")
	ticketRoutes.Use(middlewares.AuthMiddleware())
	{
		ticketRoutes.POST("/comprar/:eventoId", ticketController.Purchase)
		ticketRoutes.PATCH("/:id/cancelar", ticketController.Cancel)
		ticketRoutes.PATCH("/:id/transferir", ticketController.Transfer)
	}

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/login", authController.Login)
		authRoutes.GET("/me", middlewares.AuthMiddleware(), authController.Me)
	}

	return router
}

func corsMiddleware() gin.HandlerFunc {
	allowedOrigins := map[string]bool{
		"http://localhost:5173": true,
		"http://127.0.0.1:5173": true,
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
