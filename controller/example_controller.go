package controller

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"sr-example/example-service/model"
	"sr-example/example-service/service"
)

func exampleController() {
	router.GET("/example", getExamples)
	router.GET("/example/:id", getExample)
	router.DELETE("/example/:id", removeExample)
	router.POST("/example", addExample)
	router.PUT("/example", updateExample)
}

func getExamples(c *gin.Context) {
	examples, err := service.GetExamples()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, examples)
}

func getExample(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	example, err := service.GetExampleById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, example)
}

func removeExample(c *gin.Context) {
	id, convErr := primitive.ObjectIDFromHex(c.Param("id"))
	if convErr != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "given id was not of type ObjectID"})
		return
	}

	err := service.RemoveExampleById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusNoContent, "")
}

func addExample(c *gin.Context) {
	var example model.Example
	if err := c.BindJSON(&example); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	r, err := service.AddExample(example)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}

func updateExample(c *gin.Context) {
	var example model.Example
	if err := c.BindJSON(&example); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	r, err := service.UpdateExample(example)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, r)
}
