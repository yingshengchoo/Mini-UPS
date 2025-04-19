package controller

import (
	"mini-ups/model"
	"mini-ups/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterTruckRequest struct {
	ID model.TruckID `json:"id" binding:"required"`
	X  int           `json:"x" binding:"required"`
	Y  int           `json:"y" binding:"required"`
}

// register a truck
func RegisterTruck(c *gin.Context) {
	var request RegisterTruckRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := service.RegisterTruck(request.ID, request.X, request.Y)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "register success"})
}

// get truck info
func GetTruckInfo(c *gin.Context) {
	var request RegisterTruckRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	truck, err := service.GetTruckByID(request.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"truck": truck})
}
