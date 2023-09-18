package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/swimresults/import-service/model"
	"github.com/swimresults/import-service/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func settingsController() {
	router.GET("/settings", getImportSettings)
	router.GET("/settings/:id", getImportSetting)
	router.GET("/settings/meet/:meet_id", getImportSettingByMeetId)

	router.DELETE("/settings/:id", removeImportSetting)
	router.POST("/settings", addImportSetting)
	router.PUT("/settings", updateImportSetting)
}

func getImportSettings(c *gin.Context) {
	settings, err := service.GetImportSettings()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, settings)
}

func getImportSetting(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	setting, err := service.GetImportSettingById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, setting)
}

func getImportSettingByMeetId(c *gin.Context) {
	id := c.Param("meet_id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given meet_id is empty"})
		return
	}

	setting, err := service.GetImportSettingByMeeting(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, setting)
}

func removeImportSetting(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	err := service.RemoveImportSettingById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusNoContent, "")
}

func addImportSetting(c *gin.Context) {
	var setting model.ImportSetting
	if err := c.BindJSON(&setting); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	r, err := service.AddImportSetting(setting)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}

func updateImportSetting(c *gin.Context) {
	var setting model.ImportSetting
	if err := c.BindJSON(&setting); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	r, err := service.UpdateImportSetting(setting)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}
