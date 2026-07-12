package controllers

import (
	"net/http"

	"backend/domain"

	"github.com/gin-gonic/gin"
)

type adminReportService interface {
	GetSummary() (*domain.AdminSummaryReport, error)
	GetEventReports() ([]domain.AdminEventReport, error)
}

type AdminReportController struct {
	adminReportService adminReportService
}

func NewAdminReportController(adminReportService adminReportService) *AdminReportController {
	return &AdminReportController{adminReportService: adminReportService}
}

func (controller *AdminReportController) Summary(c *gin.Context) {
	report, err := controller.adminReportService.GetSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "could not get admin summary",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "admin summary retrieved successfully",
		"data":    report,
	})
}

func (controller *AdminReportController) EventReports(c *gin.Context) {
	reports, err := controller.adminReportService.GetEventReports()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "could not get event reports",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "event reports retrieved successfully",
		"data":    reports,
	})
}
