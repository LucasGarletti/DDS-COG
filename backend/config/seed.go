package config

import (
	"backend/domain"
	"backend/services"
	"time"

	"gorm.io/gorm"
)

const (
	SeedAdminEmail     = "admin@tickgo.com"
	SeedAdminPassword  = "Admin123!"
	SeedClientEmail    = "client@tickgo.com"
	SeedClientPassword = "Client123!"
)

func SeedDatabase(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := seedUsers(tx); err != nil {
			return err
		}

		cosquinRock, err := seedEvents(tx)
		if err != nil {
			return err
		}

		return seedCosquinRockSchedule(tx, cosquinRock.ID)
	})
}

func seedUsers(db *gorm.DB) error {
	users := []struct {
		name     string
		email    string
		password string
		role     string
	}{
		{
			name:     "Administrador TickGo",
			email:    SeedAdminEmail,
			password: SeedAdminPassword,
			role:     domain.UserRoleAdmin,
		},
		{
			name:     "Cliente de Prueba",
			email:    SeedClientEmail,
			password: SeedClientPassword,
			role:     domain.UserRoleClient,
		},
	}

	for _, item := range users {
		hashedPassword, err := services.HashPassword(item.password)
		if err != nil {
			return err
		}

		user := domain.User{
			Name:     item.name,
			Email:    item.email,
			Password: hashedPassword,
			Role:     item.role,
		}

		if err := db.Where("email = ?", item.email).FirstOrCreate(&user).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedEvents(db *gorm.DB) (*domain.Event, error) {
	festivalDate := time.Date(2027, time.February, 15, 16, 0, 0, 0, time.Local)

	events := []domain.Event{
		{
			Title:             "Festival Cosquín Rock",
			Description:       "Festival de rock con grilla oficial de artistas y escenarios.",
			Date:              festivalDate,
			Location:          "Aeródromo Santa María de Punilla, Córdoba",
			Capacity:          60000,
			AvailableCapacity: 60000,
			Price:             85000,
			ImageURL:          "",
			Status:            domain.EventStatusActive,
			IsFestival:        true,
		},
		{
			Title:             "Los Pumas en el Estadio Mario Alberto Kempes",
			Description:       "Partido de rugby de Los Pumas en Córdoba.",
			Date:              time.Date(2027, time.March, 20, 21, 0, 0, 0, time.Local),
			Location:          "Estadio Mario Alberto Kempes, Córdoba",
			Capacity:          57000,
			AvailableCapacity: 57000,
			Price:             45000,
			ImageURL:          "",
			Status:            domain.EventStatusActive,
			IsFestival:        false,
		},
		{
			Title:             "Las Pastillas del Abuelo",
			Description:       "Recital de Las Pastillas del Abuelo en Córdoba.",
			Date:              time.Date(2027, time.April, 12, 20, 30, 0, 0, time.Local),
			Location:          "Plaza de la Música, Córdoba",
			Capacity:          8000,
			AvailableCapacity: 8000,
			Price:             30000,
			ImageURL:          "",
			Status:            domain.EventStatusActive,
			IsFestival:        false,
		},
	}

	var cosquinRock domain.Event
	for _, seedEvent := range events {
		event := seedEvent
		if err := db.Where("title = ?", seedEvent.Title).FirstOrCreate(&event).Error; err != nil {
			return nil, err
		}

		if seedEvent.Title == "Festival Cosquín Rock" {
			cosquinRock = event
		}
	}

	return &cosquinRock, nil
}

func seedCosquinRockSchedule(db *gorm.DB, eventID uint) error {
	schedules := []domain.FestivalSchedule{
		{
			EventID:   eventID,
			Artist:    "Airbag",
			Stage:     "Escenario Norte",
			StartTime: time.Date(2027, time.February, 15, 18, 0, 0, 0, time.Local),
			EndTime:   time.Date(2027, time.February, 15, 19, 0, 0, 0, time.Local),
			ImageURL:  "",
		},
		{
			EventID:   eventID,
			Artist:    "Dillom",
			Stage:     "Escenario Sur",
			StartTime: time.Date(2027, time.February, 15, 19, 30, 0, 0, time.Local),
			EndTime:   time.Date(2027, time.February, 15, 20, 30, 0, 0, time.Local),
			ImageURL:  "",
		},
		{
			EventID:   eventID,
			Artist:    "Guasones",
			Stage:     "Escenario Norte",
			StartTime: time.Date(2027, time.February, 15, 21, 0, 0, 0, time.Local),
			EndTime:   time.Date(2027, time.February, 15, 22, 0, 0, 0, time.Local),
			ImageURL:  "",
		},
	}

	for _, seedSchedule := range schedules {
		schedule := seedSchedule
		if err := db.Where(
			"event_id = ? AND artist = ? AND start_time = ?",
			seedSchedule.EventID,
			seedSchedule.Artist,
			seedSchedule.StartTime,
		).FirstOrCreate(&schedule).Error; err != nil {
			return err
		}
	}

	return nil
}
