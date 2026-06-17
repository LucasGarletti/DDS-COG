package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TicketController struct {
	ticketService *services.TicketService
}

func NewTicketController(ticketService *services.TicketService) *TicketController {
	return &TicketController{ticketService: ticketService}
}

func (controller *TicketController) GetMyTickets(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "user not authenticated",
		})
		return
	}

	authenticatedUserID, ok := userID.(uint)
	if !ok || authenticatedUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "user not authenticated",
		})
		return
	}

	tickets, err := controller.ticketService.ListTicketsByUser(authenticatedUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "could not get tickets",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "tickets retrieved successfully",
		"data":    tickets,
	})
}

func (controller *TicketController) Cancel(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "user not authenticated",
		})
		return
	}

	authenticatedUserID, ok := userID.(uint)
	if !ok || authenticatedUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "user not authenticated",
		})
		return
	}

	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || ticketID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid ticket id",
		})
		return
	}

	ticket, err := controller.ticketService.CancelTicket(uint(ticketID), authenticatedUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "ticket not found",
			})
			return
		}

		if errors.Is(err, services.ErrTicketNotOwned) {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "ticket does not belong to authenticated user",
			})
			return
		}

		if errors.Is(err, services.ErrTicketAlreadyCancelled) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "ticket already cancelled",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "could not cancel ticket",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "ticket cancelled successfully",
		"data":    ticket,
	})
}

func (controller *TicketController) Purchase(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "user not authenticated",
		})
		return
	}

	authenticatedUserID, ok := userID.(uint)
	if !ok || authenticatedUserID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "user not authenticated",
		})
		return
	}

	eventID, err := strconv.ParseUint(c.Param("eventoId"), 10, 64)
	if err != nil || eventID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "invalid event id",
		})
		return
	}

	ticket, err := controller.ticketService.PurchaseTicket(services.PurchaseTicketInput{
		UserID:  authenticatedUserID,
		EventID: uint(eventID),
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "event not found",
			})
			return
		}

		if errors.Is(err, services.ErrEventSoldOut) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "event capacity exhausted",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "could not purchase ticket",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "ticket purchased successfully",
		"data": gin.H{
			"id":       ticket.ID,
			"code":     ticket.Code,
			"status":   ticket.Status,
			"event_id": ticket.EventID,
			"user_id":  ticket.UserID,
		},
	})
}
