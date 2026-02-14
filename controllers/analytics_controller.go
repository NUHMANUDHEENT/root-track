package controllers

import (
	"net/http"
	"roottrack-backend/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AnalyticsController struct {
	RoutineRepo  repositories.RoutineRepository
	SheddingRepo repositories.SheddingRepository
	ProductRepo  repositories.ProductRepository
	PhotoRepo    repositories.PhotoRepository
}

func (ctrl *AnalyticsController) GetDashboard(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	routines, _ := ctrl.RoutineRepo.GetAllByUser(userID)
	totalRoutines := len(routines)

	avgShedding, _ := ctrl.SheddingRepo.GetWeeklyAverage(userID)

	products, _ := ctrl.ProductRepo.GetAllByUser(userID)
	activeProducts := 0
	for _, p := range products {
		if p.EndDate == nil {
			activeProducts++
		}
	}

	latestPhoto, _ := ctrl.PhotoRepo.GetLatestByUser(userID)

	c.JSON(http.StatusOK, gin.H{
		"total_routines":   totalRoutines,
		"average_shedding": avgShedding,
		"active_products":  activeProducts,
		"latest_photo":     latestPhoto,
	})
}
