package controllers

import (
	"net/http"
	"schoolbackend/config"
	"schoolbackend/models"

	"github.com/gin-gonic/gin"
)

// GetDisabilities returns all disabilities from the database
func GetDisabilities(c *gin.Context) {
	var disabilities []models.DisabilityRes

	// Use Find to get all records
	if err := config.DB.Find(&disabilities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch disabilities: " + err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"data":  disabilities,
		"count": len(disabilities),
	})
}
