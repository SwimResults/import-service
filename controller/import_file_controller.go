package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/swimresults/import-service/dto"
	"github.com/swimresults/import-service/importer"
	"github.com/swimresults/import-service/model"
	"github.com/swimresults/import-service/service"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func importFileController() {
	router.POST("/file", importFile)
	router.GET("/stream/:sessionId", streamProgress)
	router.POST("/pdf_to_text", readPdfToText)
	router.POST("/certificates", importCertificates)
}

func importFile(c *gin.Context) {
	request, cleanup, err := parseImportFileRequest(c)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	err = service.ImportFile(request, cleanup)
	if err != nil {
		cleanup()
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// parseImportFileRequest supports both JSON (existing clients) and multipart uploads (new clients).
// For multipart, the file is saved to a temporary path and the request URL is set to that path.
// The caller receives a cleanup function to remove the temporary file.
func parseImportFileRequest(c *gin.Context) (model.ImportFileRequest, func(), error) {
	cleanup := func() {}
	contentType := c.ContentType()

	// Multipart file upload path
	if strings.HasPrefix(contentType, "multipart/") {
		var req model.ImportFileRequest
		if err := c.ShouldBind(&req); err != nil {
			return req, cleanup, err
		}

		fileHeader, err := c.FormFile("file")
		if err != nil {
			return req, cleanup, err
		}

		src, err := fileHeader.Open()
		if err != nil {
			return req, cleanup, err
		}
		defer src.Close()

		tmp, err := os.CreateTemp("", "import-*"+filepath.Ext(fileHeader.Filename))
		if err != nil {
			return req, cleanup, err
		}

		if _, err = io.Copy(tmp, src); err != nil {
			tmp.Close()
			os.Remove(tmp.Name())
			return req, cleanup, err
		}

		if err = tmp.Close(); err != nil {
			os.Remove(tmp.Name())
			return req, cleanup, err
		}

		req.Url = tmp.Name()

		// Infer extension from uploaded filename if not provided
		if req.FileExtension == "" {
			ext := strings.TrimPrefix(strings.ToUpper(filepath.Ext(fileHeader.Filename)), ".")
			req.FileExtension = ext
		}

		cleanup = func() {
			os.Remove(tmp.Name())
		}

		return req, cleanup, nil
	}

	// JSON path (existing behavior)
	var req model.ImportFileRequest
	if err := c.BindJSON(&req); err != nil {
		return req, cleanup, err
	}
	return req, cleanup, nil
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
