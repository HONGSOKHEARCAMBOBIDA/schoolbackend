package helper

import "github.com/gin-gonic/gin"

func GetUserID(c *gin.Context) (int, bool) {
	userID, ok := c.Get("user_id")
	if !ok {
		return 0, false
	}

	id, ok := userID.(float64)
	if !ok {
		return 0, false
	}

	return int(id), true
}
