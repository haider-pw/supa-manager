package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (a *Api) getPropsProjectJwtSecretUpdateStatus(c *gin.Context) {
	_, err := a.GetAccountFromRequest(c)
	if err != nil {
		println(err.Error())
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	// projectRef := c.Param("ref")

	// Return success status - JWT secret is already set
	c.JSON(http.StatusOK, gin.H{
		"jwtSecretUpdateStatus": gin.H{
			"status":   "success",
			"progress": 100,
		},
	})
}
