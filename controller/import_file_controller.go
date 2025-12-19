package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/swimresults/import-service/dto"
	"github.com/swimresults/import-service/importer"
	"github.com/swimresults/import-service/model"
	"github.com/swimresults/import-service/service"
	"io"
	"net/http"
)

func importFileController() {
	router.POST("/file", importFile)
	router.GET("/stream/:sessionId", streamProgress)
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

	err := service.ImportCertificates(request)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// streamProgress establishes an SSE connection for real-time progress updates
func streamProgress(c *gin.Context) {
	sessionID := c.Param("sessionId")

	// Create a new session
	session := service.CreateSession(sessionID)

	// Set headers for SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// Send initial connection success message
	service.SendLog(sessionID, "Connected to import stream", "info")

	// Stream events to client
	c.Stream(func(w io.Writer) bool {
		select {
		case <-session.Done:
			// Session closed, stop streaming
			return false
		case msg, ok := <-session.Channel:
			if !ok {
				// Channel closed
				return false
			}
			// Send SSE formatted message
			c.SSEvent("message", msg)
			return true
		}
	})

	// Cleanup when client disconnects
	service.CloseSession(sessionID)
}
