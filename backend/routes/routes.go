package routes

import (
	"net/http"

	"backend/controllers"
	"backend/dao"
	"backend/domain"
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
	adminEventService := services.NewAdminEventService(eventDAO, ticketDAO)
	adminEventController := controllers.NewAdminEventController(adminEventService)
	adminReportDAO := dao.NewAdminReportDAO(db)
	adminReportService := services.NewAdminReportService(adminReportDAO)
	adminReportController := controllers.NewAdminReportController(adminReportService)
	festivalScheduleDAO := dao.NewFestivalScheduleDAO(db)
	festivalScheduleService := services.NewFestivalScheduleService(eventDAO, festivalScheduleDAO)
	festivalScheduleController := controllers.NewFestivalScheduleController(festivalScheduleService)
	userItineraryDAO := dao.NewUserItineraryDAO(db)
	userItineraryService := services.NewUserItineraryService(eventDAO, festivalScheduleDAO, userItineraryDAO)
	userItineraryController := controllers.NewUserItineraryController(userItineraryService)

	router.GET("/eventos", eventController.GetAll)
	router.GET("/eventos/:id", eventController.GetByID)
	router.GET("/eventos/:id/grilla", festivalScheduleController.List)
	router.GET("/mis-entradas", middlewares.AuthMiddleware(), ticketController.GetMyTickets)

	itineraryRoutes := router.Group("/mis-itinerarios")
	itineraryRoutes.Use(middlewares.AuthMiddleware())
	{
		itineraryRoutes.GET("/:eventoId", userItineraryController.Get)
		itineraryRoutes.POST("/:eventoId/shows", userItineraryController.AddShow)
		itineraryRoutes.POST("/:eventoId/actividades", userItineraryController.AddPersonalActivity)
		itineraryRoutes.DELETE("/:eventoId/items/:itemId", userItineraryController.DeleteItem)
	}

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

	adminRoutes := router.Group("/admin")
	adminRoutes.Use(middlewares.AuthMiddleware())
	adminRoutes.Use(middlewares.RequireRole(domain.UserRoleAdmin))
	{
		adminRoutes.POST("/eventos", adminEventController.Create)
		adminRoutes.PATCH("/eventos/:id", adminEventController.Update)
		adminRoutes.DELETE("/eventos/:id", adminEventController.Cancel)
		adminRoutes.GET("/eventos/:id/reporte", adminEventController.Report)
		adminRoutes.GET("/reportes/resumen", adminReportController.Summary)
		adminRoutes.GET("/reportes/eventos", adminReportController.EventReports)
		adminRoutes.POST("/eventos/:id/grilla", festivalScheduleController.Create)
		adminRoutes.GET("/eventos/:id/grilla", festivalScheduleController.List)
		adminRoutes.DELETE("/grilla/:id", festivalScheduleController.Delete)
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

		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, Cache-Control, Pragma, Expires")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
