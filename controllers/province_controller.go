package controllers

import (
	"net/http"
	"schoolbackend/config"
	"schoolbackend/models"

	"github.com/gin-gonic/gin"
)

func GetProvince(c *gin.Context) {
	var province []models.Province
	config.DB.Find(&province)
	c.JSON(http.StatusOK, province)
}
