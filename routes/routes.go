package routes

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/opshive/go-td/controllers"
	_ "github.com/opshive/go-td/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "8282"
	}
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler)) //.WrapHandler(swaggerfiles.Handler))
	router.GET("/healthcheck", controllers.HealthCheck)

	v1 := router.Group("/api/v1")

	v1.POST("/tenantCreate", controllers.TenantDeploy)
	v1.GET("/tenantGet", controllers.TenantGet)
	v1.POST("/tenantDelete", controllers.TenantDelete)

	v1.POST("/appCreate", controllers.HelmlDeploy)
	v1.GET("/appGet", controllers.GetChart)
	v1.POST("/appDelete", controllers.ChartDelete)

	router.Run(fmt.Sprintf("%s:%s", host, port))
	gin.SetMode(gin.DebugMode)

}
