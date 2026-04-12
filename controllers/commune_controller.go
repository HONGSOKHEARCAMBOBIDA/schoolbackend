package controllers

import (
	"net/http"
	"schoolbackend/config"
	"schoolbackend/models"

	"github.com/gin-gonic/gin"
)

func GetCommune(c *gin.Context) {
	districtID := c.Param("id")

	var commune []models.Commune
	if err := config.DB.Where("district_id = ?", districtID).Find(&commune).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, commune)
}
func GetCommuneByID(c *gin.Context) {
	id := c.Param("id")

	var cm models.Commune
	if err := config.DB.First(&cm, "id = ?", id).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	type CommuneResp struct {
		ID         uint   `json:"id"`
		Name       string `json:"name"`
		DistrictID uint   `json:"district_id"`
	}
	c.JSON(http.StatusOK, CommuneResp{
		ID: cm.ID, Name: cm.Name, DistrictID: uint(cm.DistrictId),
	})
}
