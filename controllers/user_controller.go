package controllers

import (
	"net/http"
	"roottrack-backend/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	UserRepo repositories.UserRepository
}

func (ctrl *UserController) GetMe(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	user, err := ctrl.UserRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) Update(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.UserRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if input.Name != "" {
		user.Name = input.Name
	}
	if input.Email != "" {
		user.Email = input.Email
	}

	if err := ctrl.UserRepo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) UpdatePushToken(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	var input struct {
		Token string `json:"token"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.UserRepo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.ExpoPushToken = input.Token

	if err := ctrl.UserRepo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update push token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Push token updated successfully"})
}
