package controllers

import (
	"net/http"
	"schoolbackend/config"
	"schoolbackend/models"

	"github.com/gin-gonic/gin"
)

func GetVillage(c *gin.Context) {
	communeId := c.Param("id")

	var village []models.Village

	if err := config.DB.Where("commune_id =?", communeId).Find(&village).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, village)
}
func GetVillageByID(c *gin.Context) {
	id := c.Param("id")

	var v models.Village
	if err := config.DB.First(&v, "id = ?", id).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	// Return exactly what frontend expects
	type VillageResp struct {
		ID        uint   `json:"id"`
		Name      string `json:"name"`
		CommuneID uint   `json:"commune_id"`
	}
	c.JSON(http.StatusOK, VillageResp{
		ID: v.ID, Name: v.Name, CommuneID: uint(v.CommuneID),
	})
}
