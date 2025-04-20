package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestAction struct {
	Action string `json:"action" binding:"required"`
}

// POST /api/ups
func ParseAction(c *gin.Context) {
	var req RequestAction
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

// POST /api/ups/pickup
func PickUp(c *gin.Context) {
	// TODO implement
}

// POST /api/ups/package-ready
func RespondPackageReady(c *gin.Context) {
	// TODO implement
}

// POST /api/ups/load
func LoadPackage(c *gin.Context) {
	// TODO implement
}

// POST /api/ups/status
func CheckStatus(c *gin.Context) {
	// TODO implement
}
