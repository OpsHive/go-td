package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags healthcheck
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /healthcheck [get]
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "sucess",
	})
}
