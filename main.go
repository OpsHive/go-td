package main

import (
	"github.com/joho/godotenv"
	_ "github.com/opshive/go-td/docs"
	"github.com/opshive/go-td/routes"
)

// @title Tenant Deploy API
// @version 1.0
// @description this is Tenant deployment api which help to deploy user tenants when user signup.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email qasim@opshive.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8282
// @BasePath /api/v1
// @schemes http
func main() {
	godotenv.Load(".env")
	routes.Router()
	// router := gin.Default()

	// //router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, url))
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// router.GET("/healthcheck", controllers.HealthCheck)
	// // routes.Router()

}
