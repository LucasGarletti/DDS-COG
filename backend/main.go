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

	if err := normalizeUserRoles(); err != nil {
		log.Fatal("Error normalizing user roles: ", err)
	}

	if err := normalizeEventStatuses(); err != nil {
		log.Fatal("Error normalizing event statuses: ", err)
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

func normalizeUserRoles() error {
	return config.DB.Model(&domain.User{}).
		Where("role IS NULL OR role = ? OR role NOT IN ?", "", []string{domain.UserRoleClient, domain.UserRoleAdmin}).
		Update("role", domain.UserRoleClient).Error
}

func normalizeEventStatuses() error {
	return config.DB.Model(&domain.Event{}).
		Where("status IS NULL OR status = ? OR status NOT IN ?", "", []string{domain.EventStatusActive, domain.EventStatusCancelled}).
		Update("status", domain.EventStatusActive).Error
}
