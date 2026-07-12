package controllers

import (
	"errors"
	"net/http"
	"time"

	"backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type userItineraryService interface {
	GetItinerary(userID uint, eventID uint) (*services.ItineraryResult, error)
	AddShow(userID uint, eventID uint, scheduleID uint) (*services.AddItineraryItemResult, error)
	AddPersonalActivity(input services.AddPersonalActivityInput) (*services.AddItineraryItemResult, error)
	DeleteItem(userID uint, eventID uint, itemID uint) error
}

type UserItineraryController struct {
	itineraryService userItineraryService
}

type addShowRequest struct {
	FestivalScheduleID uint `json:"festival_schedule_id"`
}

type addPersonalActivityRequest struct {
	Title     string `json:"title"`
	Location  string `json:"location"`
	Notes     string `json:"notes"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

func NewUserItineraryController(itineraryService userItineraryService) *UserItineraryController {
	return &UserItineraryController{itineraryService: itineraryService}
}

func (controller *UserItineraryController) Get(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}
	eventID, ok := parseUintParam(c, "eventoId", "invalid event id")
	if !ok {
		return
	}

	result, err := controller.itineraryService.GetItinerary(userID, eventID)
	if err != nil {
		controller.handleItineraryError(c, err, "could not get itinerary")
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "itinerary retrieved successfully", "data": result})
}

func (controller *UserItineraryController) AddShow(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}
	eventID, ok := parseUintParam(c, "eventoId", "invalid event id")
	if !ok {
		return
	}

	var request addShowRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid request body"})
		return
	}

	result, err := controller.itineraryService.AddShow(userID, eventID, request.FestivalScheduleID)
	if err != nil {
		controller.handleItineraryError(c, err, "could not add itinerary show")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "itinerary show added successfully", "data": result})
}

func (controller *UserItineraryController) AddPersonalActivity(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}
	eventID, ok := parseUintParam(c, "eventoId", "invalid event id")
	if !ok {
		return
	}

	var request addPersonalActivityRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid request body"})
		return
	}

	startTime, err := time.Parse(time.RFC3339, request.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid itinerary item"})
		return
	}
	endTime, err := time.Parse(time.RFC3339, request.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid itinerary item"})
		return
	}

	result, err := controller.itineraryService.AddPersonalActivity(services.AddPersonalActivityInput{
		UserID:    userID,
		EventID:   eventID,
		Title:     request.Title,
		Location:  request.Location,
		Notes:     request.Notes,
		StartTime: startTime,
		EndTime:   endTime,
	})
	if err != nil {
		controller.handleItineraryError(c, err, "could not add itinerary activity")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "itinerary activity added successfully", "data": result})
}

func (controller *UserItineraryController) DeleteItem(c *gin.Context) {
	userID, ok := authenticatedUserID(c)
	if !ok {
		return
	}
	eventID, ok := parseUintParam(c, "eventoId", "invalid event id")
	if !ok {
		return
	}
	itemID, ok := parseUintParam(c, "itemId", "invalid itinerary item id")
	if !ok {
		return
	}

	if err := controller.itineraryService.DeleteItem(userID, eventID, itemID); err != nil {
		controller.handleItineraryError(c, err, "could not delete itinerary item")
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "itinerary item deleted successfully"})
}

func (controller *UserItineraryController) handleItineraryError(c *gin.Context, err error, fallback string) {
	if errors.Is(err, services.ErrInvalidItineraryItem) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid itinerary item"})
		return
	}

	if errors.Is(err, services.ErrEventNotFestival) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "event is not festival"})
		return
	}

	if errors.Is(err, services.ErrItineraryForbidden) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "active ticket required"})
		return
	}

	if errors.Is(err, services.ErrItineraryItemNotOwned) || errors.Is(err, services.ErrScheduleEventMismatch) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "forbidden itinerary operation"})
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "itinerary resource not found"})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": fallback})
}

func authenticatedUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "user not authenticated"})
		return 0, false
	}

	authenticatedUserID, ok := userID.(uint)
	if !ok || authenticatedUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "user not authenticated"})
		return 0, false
	}

	return authenticatedUserID, true
}
