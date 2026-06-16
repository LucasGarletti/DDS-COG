package main

import (
	"log"
	"os"

	"backend/config"
	"backend/domain"
	"backend/routes"
)

func main() {
	config.LoadEnv()
	config.ConnectDatabase()

	if err := config.DB.AutoMigrate(&domain.User{}, &domain.Event{}, &domain.Ticket{}); err != nil {
		log.Fatal("Error running database migrations: ", err)
	}

	router := routes.SetupRouter(config.DB)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
