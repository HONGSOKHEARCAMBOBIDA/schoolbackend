package controllers

import (
	"net/http"
	"schoolbackend/config"
	"schoolbackend/models"

	"github.com/gin-gonic/gin"
)

func GetDistrict(c *gin.Context) {
	provinceID := c.Param("id")
	var district []models.District
	if err := config.DB.Where("province_id = ?", provinceID).Find(&district).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, district)
}

// GET /district-by-id/:id
func GetDistrictByID(c *gin.Context) {
	id := c.Param("id")

	var d models.District
	if err := config.DB.First(&d, "id = ?", id).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	type DistrictResp struct {
		ID         uint   `json:"id"`
		Name       string `json:"name"`
		ProvinceID uint   `json:"province_id"`
	}
	c.JSON(http.StatusOK, DistrictResp{
		ID: d.ID, Name: d.Name, ProvinceID: uint(d.ProvinceID),
	})
}
