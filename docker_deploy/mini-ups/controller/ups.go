package controller

import (
	"mini-ups/model"
	"mini-ups/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// POST /api/ups
func ParseAction(c *gin.Context) {

	// json with action
	var req struct {
		Action string `json:"action" binding:"required"`
	}

	// parse action
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no action in json"})
		return
	}

	// forward by action
	switch req.Action {
	case "request_pickup":
		{
			PickUp(c)
			break
		}
	case "package_ready":
		{
			RespondPackageReady(c)
			break
		}
	case "load_package":
		{
			LoadPackage(c)
			break
		}
	case "query_status":
		{
			CheckStatus(c)
			break
		}
	default:
		{
			c.JSON(http.StatusBadRequest, gin.H{"error": "unknown action <" + req.Action + "> in json"})
			break
		}
	}
}

// TODO test
// POST /api/ups/pickup
func PickUp(c *gin.Context) {
	CreatePackage(c)
}

// POST /api/ups/package-ready
func RespondPackageReady(c *gin.Context) {
	// TODO implement
}

// POST /api/ups/load
func LoadPackage(c *gin.Context) {
	// TODO implement
}

// TODO test
// POST /api/ups/status
func CheckStatus(c *gin.Context) {
	// json with action
	var req struct {
		PackageID string `json:"package_id" binding:"required"`
	}

	// get package info
	pack, err := service.GetPackageInfo(req.PackageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Package not found"})
		return
	}

	// construct response
	response := struct {
		model.Package
		Action string `json:"action"`
	}{
		Package: *pack,
		Action:  "query_status_response",
	}

	c.JSON(http.StatusOK, response)
}
