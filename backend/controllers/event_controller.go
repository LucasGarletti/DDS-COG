package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EventController struct {
	eventService *services.EventService
}

func NewEventController(eventService *services.EventService) *EventController {
	return &EventController{eventService: eventService}
}

func (controller *EventController) GetAll(c *gin.Context) {
	events, err := controller.eventService.ListEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "could not get events",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "events retrieved successfully",
		"data":    events,
	})
}

func (controller *EventController) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid event id",
		})
		return
	}

	event, err := controller.eventService.GetEventByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "event not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "could not get event",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "event retrieved successfully",
		"data":    event,
	})
}
