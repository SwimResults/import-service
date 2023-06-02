package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func timingSoftwareController() {
	router.POST("/easywk", easyWkLivetiming)
}

func easyWkLivetiming(c *gin.Context) {
	c.String(http.StatusNotImplemented, "")
}
