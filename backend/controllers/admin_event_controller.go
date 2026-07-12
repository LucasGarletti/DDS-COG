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

type adminEventService interface {
	CreateEvent(input services.CreateEventInput) (*domain.Event, error)
	UpdateEvent(input services.UpdateEventInput) (*domain.Event, error)
	CancelEvent(id uint) (*domain.Event, error)
	GetEventReport(id uint) (*services.EventReport, error)
}

type AdminEventController struct {
	adminEventService adminEventService
}

type createEventRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	Location    string  `json:"location"`
	Capacity    int     `json:"capacity"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
}

type updateEventRequest struct {
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	Date        *string  `json:"date"`
	Location    *string  `json:"location"`
	Capacity    *int     `json:"capacity"`
	Price       *float64 `json:"price"`
	ImageURL    *string  `json:"image_url"`
}

func NewAdminEventController(adminEventService adminEventService) *AdminEventController {
	return &AdminEventController{adminEventService: adminEventService}
}

func (controller *AdminEventController) Create(c *gin.Context) {
	var request createEventRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid request body",
		})
		return
	}

	date, err := parseRequiredDate(request.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid event data",
		})
		return
	}

	event, err := controller.adminEventService.CreateEvent(services.CreateEventInput{
		Title:       request.Title,
		Description: request.Description,
		Date:        date,
		Location:    request.Location,
		Capacity:    request.Capacity,
		Price:       request.Price,
		ImageURL:    request.ImageURL,
	})
	if err != nil {
		controller.handleAdminEventError(c, err, "could not create event")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "event created successfully",
		"data":    event,
	})
}

func (controller *AdminEventController) Update(c *gin.Context) {
	id, ok := parseEventID(c)
	if !ok {
		return
	}

	var request updateEventRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid request body",
		})
		return
	}

	var date *time.Time
	if request.Date != nil {
		parsedDate, err := parseRequiredDate(*request.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "invalid event data",
			})
			return
		}
		date = &parsedDate
	}

	input := services.UpdateEventInput{
		ID:          id,
		Title:       request.Title,
		Description: request.Description,
		Location:    request.Location,
		Capacity:    request.Capacity,
		Price:       request.Price,
		ImageURL:    request.ImageURL,
	}
	if date != nil {
		input.Date = date
	}

	event, err := controller.adminEventService.UpdateEvent(input)
	if err != nil {
		controller.handleAdminEventError(c, err, "could not update event")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "event updated successfully",
		"data":    event,
	})
}

func (controller *AdminEventController) Cancel(c *gin.Context) {
	id, ok := parseEventID(c)
	if !ok {
		return
	}

	event, err := controller.adminEventService.CancelEvent(id)
	if err != nil {
		controller.handleAdminEventError(c, err, "could not cancel event")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "event cancelled successfully",
		"data":    event,
	})
}

func (controller *AdminEventController) Report(c *gin.Context) {
	id, ok := parseEventID(c)
	if !ok {
		return
	}

	report, err := controller.adminEventService.GetEventReport(id)
	if err != nil {
		controller.handleAdminEventError(c, err, "could not get event report")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "event report retrieved successfully",
		"data":    report,
	})
}

func (controller *AdminEventController) handleAdminEventError(c *gin.Context, err error, fallback string) {
	if errors.Is(err, services.ErrInvalidEventData) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid event data",
		})
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "event not found",
		})
		return
	}

	if errors.Is(err, services.ErrEventCancelled) {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   "event is cancelled",
		})
		return
	}

	if errors.Is(err, services.ErrEventCapacityConflict) {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"error":   "event capacity cannot be lower than issued tickets",
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"error":   fallback,
	})
}

func parseEventID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid event id",
		})
		return 0, false
	}

	return uint(id), true
}

func parseRequiredDate(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}
