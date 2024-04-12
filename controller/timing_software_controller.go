package controller

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/swimresults/import-service/model"
	"github.com/swimresults/import-service/service"
	"io"
	"net/http"
)

func timingSoftwareController() {
	router.POST("/easywk/livework.php", easyWkLivetiming)
	router.GET("/easywk/livework.php", easyWkLivetiming)
	router.OPTIONS("/easywk/livework.php", ok)

	router.POST("/easywk/v2/livework.php", easyWkLivetimingV2)
	router.GET("/easywk/v2/livework.php", easyWkLivetimingV2)
	router.OPTIONS("/easywk/v2/livework.php", ok)

	router.POST("/easywk/v3", easyWkLivetimingV3)
	router.GET("/easywk/v3", easyWkLivetimingV3)
	router.OPTIONS("/easywk/v3", ok)
}

func easyWkLivetiming(c *gin.Context) {
	var request model.EasyWkActionRequest

	body, _ := io.ReadAll(c.Request.Body)
	println(string(body))

	c.Request.URL.RawQuery = string(body)

	if err := c.BindQuery(&request); err != nil {
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

func easyWkLivetimingV3(c *gin.Context) {
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

func easyWkLivetimingV2(c *gin.Context) {
	var request []model.EasyWkAction

	body, _ := io.ReadAll(c.Request.Body)
	println(string(body))

	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	str, err := service.EasyWkLivetimingRequestV2(request)
	if err != nil {
		c.String(http.StatusInternalServerError, "ERROR: %s", err.Error())
		return
	}

	c.String(http.StatusOK, str)
}

func ok(c *gin.Context) {
	c.Status(200)
}
