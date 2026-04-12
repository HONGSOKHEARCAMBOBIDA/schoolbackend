package controllers

import (
	"net/http"
	"schoolbackend/config"
	"schoolbackend/models"

	"github.com/gin-gonic/gin"
)

func GetTypeExam(c *gin.Context) {
	var typeaxe []models.TypeExam
	if err := config.DB.Find(&typeaxe).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, typeaxe)
}
