package controller

import (
	"mini-ups/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

package controller

import (
	"mini-ups/dao"
	"mini-ups/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /api/user/:userID/packages
func GetPackagesByUser(c *gin.Context) {
	userIDParam := c.Param("userID")
	var userID uint
	if _, err := fmt.Sscanf(userIDParam, "%d", &userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	packages, err := dao.GetPackagesByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"packages": packages})
}

// GET /api/package/:packageID
func GetPackageByID(c *gin.Context) {
	packageID := c.Param("packageID")
	packageInfo, err := dao.GetPackagesByPackageID(packageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"package": packageInfo})
}

// GET /api/track/:packageID
func TrackPackage(c *gin.Context) {
	packageID := c.Param("packageID")
	pkg, err := dao.GetPackagesByPackageID(packageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Package not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         pkg.ID,
		"status":     pkg.Status,
		"location":   pkg.Destination,
		"warehouse":  pkg.WarehouseID,
		"truck_id":   pkg.TruckID,
		"items":      pkg.Items,
	})
}
