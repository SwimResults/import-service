package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/swimresults/import-service/model"
	"github.com/swimresults/import-service/service"
	"net/http"
)

func timingSoftwareController() {
	router.POST("/easywk", easyWkLivetiming)
	router.OPTIONS("/easywk", ok)
}

func easyWkLivetiming(c *gin.Context) {
	var request model.EasyWkActionRequest
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	str, err := service.EasyWkLivetimingRequest(request)
	if err != nil {
		c.String(http.StatusInternalServerError, "ERROR: %s", err.Error())
		return
	}

	c.String(http.StatusOK, str)
}

func ok(c *gin.Context) {
	c.Status(200)
}
