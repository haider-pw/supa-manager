package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getPlatformNotificationsSummary(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		println(err.Error())
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// Return empty notification summary
	// TODO: Implement actual notification tracking in future phases
	c.JSON(http.StatusOK, gin.H{
		"total":  0,
		"unread": 0,
		"new":    0,
	})
}
