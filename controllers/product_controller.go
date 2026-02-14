package controllers

import (
	"fmt"
	"net/http"
	"roottrack-backend/models"
	"roottrack-backend/repositories"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductController struct {
	Repo repositories.ProductRepository
}

func (ctrl *ProductController) Create(c *gin.Context) {
	var input struct {
		Name      string `json:"name" binding:"required"`
		Brand     string `json:"brand" binding:"required"`
		IsActive  bool   `json:"is_active"`
		StartDate string `json:"start_date" binding:"omitempty"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println("Error binding JSON: ", err.Error())
		return
	}

	userIDStr := c.MustGet("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	startDate := time.Now()
	if input.StartDate != "" {
		if t, err := time.Parse(time.RFC3339, input.StartDate); err == nil {
			startDate = t
		}
	}

	product := models.Product{
		UserID:    userID,
		Name:      input.Name,
		Brand:     input.Brand,
		IsActive:  input.IsActive,
		StartDate: startDate,
	}

	if err := ctrl.Repo.Create(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (ctrl *ProductController) GetAll(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	products, err := ctrl.Repo.GetAllByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (ctrl *ProductController) Update(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	product, err := ctrl.Repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var input struct {
		Name     string `json:"name"`
		Brand    string `json:"brand"`
		IsActive *bool  `json:"is_active"`
		EndDate  string `json:"end_date"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != "" {
		product.Name = input.Name
	}
	if input.Brand != "" {
		product.Brand = input.Brand
	}
	if input.IsActive != nil {
		product.IsActive = *input.IsActive
	}
	if input.EndDate != "" {
		endDate, _ := time.Parse(time.RFC3339, input.EndDate)
		product.EndDate = &endDate
	}

	if err := ctrl.Repo.Update(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (ctrl *ProductController) Delete(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	if err := ctrl.Repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
