package controllers

import (
	"net/http"
	"schoolbackend/config"
	"schoolbackend/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleInput struct {
	Name string `json:"name"`
}

func CreateRole(c *gin.Context) {
	var input RoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx := config.DB.Begin()
	role := models.Role{
		Name:     input.Name,
		IsActive: 1,
	}
	if err := tx.Create(&role).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, http.StatusOK)

}
func GetRole(c *gin.Context) {
	var roles []models.Role

	err := config.DB.
		Table("roles").
		Select("roles.*, COUNT(users.id) as count_user").
		Joins("LEFT JOIN users ON users.role_id = roles.id").
		Group("roles.id").
		Scan(&roles).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, roles)
}
func UpdateRole(c *gin.Context) {
	id := c.Param("id")

	// 1. Find role
	var role models.Role
	if err := config.DB.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	// 2. Bind JSON
	var input RoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Update fields
	updates := map[string]interface{}{
		"name": input.Name,
	}

	if err := config.DB.Model(&role).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}
func ChangeStatusRole(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Toggle status
	if err := config.DB.Model(&models.Role{}).
		Where("id = ?", id).
		Update("is_active", gorm.Expr("CASE WHEN is_active = 1 THEN 0 ELSE 1 END")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch updated role
	var role models.Role
	if err := config.DB.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, role)
}
