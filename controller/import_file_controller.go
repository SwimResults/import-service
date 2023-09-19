package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/swimresults/import-service/model"
	"github.com/swimresults/import-service/service"
	"net/http"
)

func importFileController() {
	router.POST("/file", importFile)
}

func importFile(c *gin.Context) {
	var request model.ImportFileRequest
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	err := service.ImportFile(request)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
