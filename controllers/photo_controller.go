package controllers

import (
	"net/http"
	"roottrack-backend/models"
	"roottrack-backend/repositories"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PhotoController struct {
	Repo repositories.PhotoRepository
}

func (ctrl *PhotoController) Create(c *gin.Context) {
	var input struct {
		ImageURL string `json:"image_url" binding:"required"`
		TakenAt  string `json:"taken_at" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := c.MustGet("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	takenAt, _ := time.Parse(time.RFC3339, input.TakenAt)

	photo := models.ProgressPhoto{
		UserID:   userID,
		ImageURL: input.ImageURL,
		TakenAt:  takenAt,
	}

	if err := ctrl.Repo.Create(&photo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save photo"})
		return
	}

	c.JSON(http.StatusCreated, photo)
}

func (ctrl *PhotoController) GetAll(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	photos, err := ctrl.Repo.GetAllByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch photos"})
		return
	}

	c.JSON(http.StatusOK, photos)
}

func (ctrl *PhotoController) Delete(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	if err := ctrl.Repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete photo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Photo deleted"})
}
