package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swimresults/import-service/service"
	"net/http"
	"os"
)

var router = gin.Default()
var serviceKey string

func Run() {

	port := os.Getenv("SR_IMPORT_PORT")

	if port == "" {
		fmt.Println("no application port given! Please set SR_IMPORT_PORT.")
		return
	}

	serviceKey = os.Getenv("SR_SERVICE_KEY")

	if serviceKey == "" {
		fmt.Println("no security for inter-service communication given! Please set SR_SERVICE_KEY.")
	}

	timingSoftwareController()
	importFileController()
	settingsController()

	router.GET("/actuator", actuator)

	err := router.Run(":" + port)
	if err != nil {
		fmt.Println("Unable to start application on port " + port)
		return
	}
}

func actuator(c *gin.Context) {

	state := "OPERATIONAL"

	if !service.PingDatabase() {
		state = "DATABASE_DISCONNECTED"
	}
	c.String(http.StatusOK, state)
}
