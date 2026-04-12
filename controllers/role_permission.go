package controllers

import (
	"net/http"
	"schoolbackend/config"
	"schoolbackend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateRolePermissionsInput struct {
	RoleID        int   `json:"role_id" binding:"required"`
	PermissionIDs []int `json:"permission_id" binding:"required"`
}

func CreateRolePermissions(c *gin.Context) {
	var input CreateRolePermissionsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var rolePermissions []models.RolePermission
	for _, pid := range input.PermissionIDs {
		rolePermissions = append(rolePermissions, models.RolePermission{
			RoleID:       uint(input.RoleID),
			PermissionID: uint(pid),
		})
	}

	if err := config.DB.Create(&rolePermissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rolePermissions)
}

type DeleteRolePermissionsInput struct {
	RoleID        int   `json:"role_id" binding:"required"`
	PermissionIDs []int `json:"permission_id" binding:"required"`
}

func DeleteRolePermissions(c *gin.Context) {
	var input DeleteRolePermissionsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Where("role_id = ? AND permission_id IN ?", input.RoleID, input.PermissionIDs).
		Delete(&models.RolePermission{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permissions deleted successfully"})
}

type PermissionWithAssigned struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Assigned    bool   `json:"assigned"`
}

func GetRolePermissions(c *gin.Context) {
	roleIDParam := c.Param("id")
	roleID, err := strconv.Atoi(roleIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	var permissions []PermissionWithAssigned

	err = config.DB.
		Table("permissions").
		Select(`permissions.id, permissions.display_name,permissions.name, 
		        CASE WHEN role_has_permissions.permission_id IS NULL THEN false ELSE true END AS assigned`).
		Joins(`LEFT JOIN role_has_permissions 
		       ON permissions.id = role_has_permissions.permission_id 
		       AND role_has_permissions.role_id = ?`, roleID).
		Scan(&permissions).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}
