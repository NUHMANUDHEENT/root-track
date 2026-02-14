package controllers

import (
	"fmt"
	"net/http"
	"roottrack-backend/models"
	"roottrack-backend/repositories"
	"roottrack-backend/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoutineController struct {
	Repo repositories.RoutineRepository
}

func (ctrl *RoutineController) Create(c *gin.Context) {
	var input struct {
		Title      string   `json:"title" binding:"required"`
		ProductID  string   `json:"product_id"`
		StartDate  string   `json:"start_date" binding:"required"`
		EndDate    string   `json:"end_date"`
		Recurrence string   `json:"recurrence" binding:"required"` // none, daily, weekly
		DaysOfWeek []string `json:"days_of_week"`
		Completed  bool     `json:"completed"`
		Notes      string   `json:"notes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := c.MustGet("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	startDate, err := time.Parse(time.RFC3339, input.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date"})
		return
	}

	var endDate *time.Time
	if input.EndDate != "" {
		parsedEnd, err := time.Parse(time.RFC3339, input.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date"})
			return
		}
		endDate = &parsedEnd
	}

	// Validate recurrence logic
	if input.Recurrence == "weekly" && len(input.DaysOfWeek) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Days of week required for weekly recurrence"})
		return
	}

	var productID *uuid.UUID
	if input.ProductID != "" {
		pID, _ := uuid.Parse(input.ProductID)
		productID = &pID
	}

	routine := models.Routine{
		UserID:     userID,
		ProductID:  productID,
		Title:      input.Title,
		StartDate:  startDate,
		EndDate:    endDate,
		Recurrence: input.Recurrence,
		DaysOfWeek: input.DaysOfWeek,
		Completed:  input.Completed,
		Notes:      input.Notes,
	}

	if err := ctrl.Repo.Create(&routine); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create routine"})
		return
	}

	c.JSON(http.StatusCreated, routine)
}

func (ctrl *RoutineController) GetAll(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	routines, err := ctrl.Repo.GetAllByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch routines"})
		return
	}

	c.JSON(http.StatusOK, routines)
}

func (ctrl *RoutineController) GetToday(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	today := time.Now().UTC()
	weekday := strings.ToLower(today.Weekday().String())

	routines, err := ctrl.Repo.GetActiveByUser(userID, today)
	if err != nil {
		fmt.Println("Failed to fetch routines", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch routines"})
		return
	}

	todayRoutines := []models.Routine{}

	for _, routine := range routines {

		switch routine.Recurrence {

		case "none":
			if utils.SameDate(routine.StartDate, today) {
				todayRoutines = append(todayRoutines, routine)
			}

		case "daily":
			todayRoutines = append(todayRoutines, routine)

		case "weekly":
			for _, day := range routine.DaysOfWeek {
				if strings.ToLower(day) == weekday {
					todayRoutines = append(todayRoutines, routine)
					break
				}
			}
		}
	}

	c.JSON(http.StatusOK, todayRoutines)
}

func (ctrl *RoutineController) GetByID(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	routine, err := ctrl.Repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Routine not found"})
		return
	}
	c.JSON(http.StatusOK, routine)
}

func (ctrl *RoutineController) Update(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	routine, err := ctrl.Repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Routine not found"})
		return
	}

	var input struct {
		Title     string `json:"title"`
		ProductID string `json:"product_id"`
		Completed *bool  `json:"completed"`
		Notes     string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Title != "" {
		routine.Title = input.Title
	}
	if input.ProductID != "" {
		pID, _ := uuid.Parse(input.ProductID)
		routine.ProductID = &pID
	}
	if input.Completed != nil {
		routine.Completed = *input.Completed
	}
	if input.Notes != "" {
		routine.Notes = input.Notes
	}

	if err := ctrl.Repo.Update(routine); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update routine"})
		return
	}

	c.JSON(http.StatusOK, routine)
}

func (ctrl *RoutineController) Delete(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	if err := ctrl.Repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete routine"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Routine deleted"})
}
