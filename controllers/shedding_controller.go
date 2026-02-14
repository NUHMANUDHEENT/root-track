package controllers

import (
	"net/http"
	"roottrack-backend/models"
	"roottrack-backend/repositories"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SheddingController struct {
	Repo repositories.SheddingRepository
}

func (ctrl *SheddingController) Create(c *gin.Context) {
	var input struct {
		Count int    `json:"shedding_count" binding:"required"`
		Date  string `json:"date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := c.MustGet("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	date, _ := time.Parse(time.RFC3339, input.Date)

	log := models.SheddingLog{
		UserID:        userID,
		SheddingCount: input.Count,
		Date:          date,
	}

	if err := ctrl.Repo.Create(&log); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log shedding"})
		return
	}

	c.JSON(http.StatusCreated, log)
}

func (ctrl *SheddingController) GetAll(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	logs, err := ctrl.Repo.GetAllByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shedding logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}

func (ctrl *SheddingController) GetSummary(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	avg, err := ctrl.Repo.GetWeeklyAverage(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch summary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"weekly_average": avg})
}
