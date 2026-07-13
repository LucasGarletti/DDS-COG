package controllers

import (
	"backend/domain"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

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
	filters, ok := parseEventFilters(c)
	if !ok {
		return
	}

	events, err := controller.eventService.ListEvents(filters)
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

func parseEventFilters(c *gin.Context) (domain.EventFilters, bool) {
	var filters domain.EventFilters

	filters.Search = strings.TrimSpace(c.Query("search"))
	filters.Location = strings.TrimSpace(c.Query("location"))

	if value := strings.TrimSpace(c.Query("date_from")); value != "" {
		dateFrom, err := parseFilterDate(value, false)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid date_from"})
			return filters, false
		}
		filters.DateFrom = &dateFrom
	}

	if value := strings.TrimSpace(c.Query("date_to")); value != "" {
		dateTo, err := parseFilterDate(value, true)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid date_to"})
			return filters, false
		}
		filters.DateTo = &dateTo
	}

	if value := strings.TrimSpace(c.Query("min_price")); value != "" {
		minPrice, err := strconv.ParseFloat(value, 64)
		if err != nil || minPrice < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid min_price"})
			return filters, false
		}
		filters.MinPrice = &minPrice
	}

	if value := strings.TrimSpace(c.Query("max_price")); value != "" {
		maxPrice, err := strconv.ParseFloat(value, 64)
		if err != nil || maxPrice < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid max_price"})
			return filters, false
		}
		filters.MaxPrice = &maxPrice
	}

	if value := strings.TrimSpace(c.Query("is_festival")); value != "" {
		isFestival, ok := parseStrictBool(value)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid is_festival"})
			return filters, false
		}
		filters.IsFestival = &isFestival
	}

	if value := strings.TrimSpace(c.Query("available")); value != "" {
		available, ok := parseStrictBool(value)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid available"})
			return filters, false
		}
		filters.AvailableOnly = available
	}

	if value := strings.TrimSpace(c.Query("sort")); value != "" {
		if !isValidEventSort(value) {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid sort"})
			return filters, false
		}
		filters.Sort = value
	}

	return filters, true
}

func parseFilterDate(value string, endOfDay bool) (time.Time, error) {
	if date, err := time.Parse("2006-01-02", value); err == nil {
		if endOfDay {
			return date.AddDate(0, 0, 1).Add(-time.Nanosecond), nil
		}
		return date, nil
	}

	return time.Parse(time.RFC3339, value)
}

func parseStrictBool(value string) (bool, bool) {
	switch strings.ToLower(value) {
	case "true":
		return true, true
	case "false":
		return false, true
	default:
		return false, false
	}
}

func isValidEventSort(value string) bool {
	switch value {
	case domain.EventSortDateAsc, domain.EventSortDateDesc, domain.EventSortPriceAsc, domain.EventSortPriceDesc:
		return true
	default:
		return false
	}
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
