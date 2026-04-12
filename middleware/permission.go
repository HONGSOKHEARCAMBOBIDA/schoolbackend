package middleware

import (
	"schoolbackend/config" // import database config
	"schoolbackend/models"

	// import data model
	"net/http"            // HTTP status code
	"schoolbackend/utils" // import util function
	"strings"             // manipulate string

	"github.com/dgrijalva/jwt-go" // for JWT parsing and validation
	"github.com/gin-gonic/gin"    // Gin web framework
	"gorm.io/gorm"                // GORM ORM for database
)

func PermissionMiddleware(requiredPermissions ...string) gin.HandlerFunc {

	// requirdpermisssion...string mean that accept multiple permission name
	// return gin.HandlerFunc used by Gin framework to handle HTTP request
	// start middleware logic
	return func(c *gin.Context) {
		// Get token from header
		authHeader := c.GetHeader("Authorization")
		// get and validate JWT Token
		// get token from Authoriczation header
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
			c.Abort()
			//  Stops the request and sends a response with a status code.
			return
		}
		// if header is empty return 401

		parts := strings.Split(authHeader, " ")
		// split the string authHeader into substring
		// example parts:= string.split(authHeader," ")
		// parts will be []string{"Bearer","mytokent"}
		// part[0]->"Bearer"
		// part[1]->"mytokent"
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
			c.Abort()
			return
		}
		// expected format Bearer <token>
		// if it doesn't match return 401

		tokenString := parts[1]
		// extract the actual token value
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return utils.JwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
			c.Abort()
			return
		}
		// if have error token is valid return 401

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
			c.Abort()
			return
		}
		// token.claim this is payload of the JWT
		//
		// Converts the claims to a map so you can access them (e.g., user ID, role ID).
		// expect the token type jwt.MapClaims
		// try to convert the token.claim to map (jwt.MapClaims)

		roleIDFloat, ok := claims["role_id"].(float64)
		// accesss value of key from JWT claim expext this value to be float64
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User role not found"})
			c.Abort()
			return
		}
		roleID := uint(roleIDFloat)
		c.Set("userID", roleID)
		// Query DB: Load role with permissions
		var role models.Role
		err = config.DB.Preload("Permissions").First(&role, roleID).Error
		// find first record the roles table where id = roleID and store in role
		// .Preload("Permissions") "Also load related data from the permissions table linked to this role.
		// First($role,roleID) get the role from roles table
		// Preload("Permissions") → fetches related permissions by joining through the pivot table.
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusForbidden, gin.H{"message": "Role not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
			}
			c.Abort()
			return
		}

		// Check if role has any of the required permissions
		hasPermission := false
		for _, perm := range role.Permissions {
			// role.Permission list permission the role has
			for _, required := range requiredPermissions {
				//requiredPermissions comes from how the middleware is used.
				if perm.Name == required {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"message": "You don't have permission to access this resource"})
			c.Abort()
			return
		}

		// User has permission, continue
		c.Next()
	}
}
