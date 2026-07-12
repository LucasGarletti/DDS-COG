package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"backend/domain"
	"backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type festivalScheduleService interface {
	Create(input services.CreateFestivalScheduleInput) (*domain.FestivalSchedule, error)
	List(eventID uint) ([]domain.FestivalSchedule, error)
	Delete(id uint) error
}

type FestivalScheduleController struct {
	scheduleService festivalScheduleService
}

type createFestivalScheduleRequest struct {
	Artist    string `json:"artist"`
	Stage     string `json:"stage"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	ImageURL  string `json:"image_url"`
}

func NewFestivalScheduleController(scheduleService festivalScheduleService) *FestivalScheduleController {
	return &FestivalScheduleController{scheduleService: scheduleService}
}

func (controller *FestivalScheduleController) Create(c *gin.Context) {
	eventID, ok := parseUintParam(c, "id", "invalid event id")
	if !ok {
		return
	}

	var request createFestivalScheduleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid request body"})
		return
	}

	startTime, err := time.Parse(time.RFC3339, request.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid festival schedule"})
		return
	}

	endTime, err := time.Parse(time.RFC3339, request.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid festival schedule"})
		return
	}

	schedule, err := controller.scheduleService.Create(services.CreateFestivalScheduleInput{
		EventID:   eventID,
		Artist:    request.Artist,
		Stage:     request.Stage,
		StartTime: startTime,
		EndTime:   endTime,
		ImageURL:  request.ImageURL,
	})
	if err != nil {
		controller.handleScheduleError(c, err, "could not create festival schedule")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "festival schedule created successfully", "data": schedule})
}

func (controller *FestivalScheduleController) List(c *gin.Context) {
	eventID, ok := parseUintParam(c, "id", "invalid event id")
	if !ok {
		return
	}

	schedules, err := controller.scheduleService.List(eventID)
	if err != nil {
		controller.handleScheduleError(c, err, "could not get festival schedule")
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "festival schedule retrieved successfully", "data": schedules})
}

func (controller *FestivalScheduleController) Delete(c *gin.Context) {
	scheduleID, ok := parseUintParam(c, "id", "invalid schedule id")
	if !ok {
		return
	}

	if err := controller.scheduleService.Delete(scheduleID); err != nil {
		controller.handleScheduleError(c, err, "could not delete festival schedule")
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "festival schedule deleted successfully"})
}

func (controller *FestivalScheduleController) handleScheduleError(c *gin.Context, err error, fallback string) {
	if errors.Is(err, services.ErrInvalidFestivalSchedule) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid festival schedule"})
		return
	}

	if errors.Is(err, services.ErrEventNotFestival) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "event is not festival"})
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "festival schedule not found"})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": fallback})
}

func parseUintParam(c *gin.Context, name string, message string) (uint, bool) {
	id, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": message})
		return 0, false
	}

	return uint(id), true
}
