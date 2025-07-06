package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/swimresults/import-service/dto"
	"github.com/swimresults/import-service/importer"
	"github.com/swimresults/import-service/model"
	"github.com/swimresults/import-service/service"
	"net/http"
)

func importFileController() {
	router.POST("/file", importFile)
	router.POST("/pdf_to_text", readPdfToText)
	router.POST("/certificates", importCertificates)
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

func readPdfToText(c *gin.Context) {
	var request model.ImportFileRequest
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	text, err := importer.GetPdfFileContent(request.Url)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	request.Text = text
	c.IndentedJSON(http.StatusOK, request)
}

func importCertificates(c *gin.Context) {
	var request dto.ImportCertificatesRequestDto
	if err := c.BindJSON(&request); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	stats, err := importer.ImportCertificates(request.Directory, request.Meeting)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, stats)
}
