package controller

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swimresults/import-service/model"
	"github.com/swimresults/import-service/service"
	"io"
	"net/http"
)

func timingSoftwareController() {
	router.POST("/easywk/livework.php", easyWkLivetiming) // used by relay server at https://weiss-konrad.de/swimresults/livetiming/livework.php
	router.GET("/easywk/get/livework.php", easyWkLivetimingGet)
	router.Any("/easywk/any/livework.php", easyWkLivetimingGet)
	router.OPTIONS("/easywk/livework.php", ok)
	router.OPTIONS("/easywk/get/livework.php", ok)

	router.POST("/easywk/v2/livework.php", easyWkLivetimingV2)
	router.GET("/easywk/v2/livework.php", easyWkLivetimingV2)
	router.OPTIONS("/easywk/v2/livework.php", ok)

	router.POST("/easywk/v3", easyWkLivetimingV3) // used from admin backend
	router.GET("/easywk/v3", easyWkLivetimingV3)
	router.OPTIONS("/easywk/v3", ok)

	router.POST("/alge", algeLivetimingData)
}

func easyWkLivetimingGet(c *gin.Context) {
	var request model.EasyWkActionRequest

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		println(err)
	}
	println(string(body))

	c.Request.URL.RawQuery = string(body)

	paramPairs := c.Request.URL.Query()
	for key, values := range paramPairs {
		fmt.Printf("key = %v, value(s) = %v\n", key, values)
	}

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
	var request model.EasyWkActionV3Request

	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	println(request.Password)

	str, err := service.EasyWkLivetimingRequest(model.EasyWkActionRequest(request))
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

func algeLivetimingData(c *gin.Context) {
	var request model.AlgeActionRequest

	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	str, err := service.AlgeLivetimingRequest(request)
	if err != nil {
		c.String(http.StatusInternalServerError, "ERROR: %s", err.Error())
		return
	}

	c.String(http.StatusOK, str)
}

func ok(c *gin.Context) {
	c.Status(200)
}
