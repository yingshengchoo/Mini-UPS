package controller

import (
	"fmt"
	"mini-ups/model"
	"mini-ups/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /api/package/user/:username
func GetPackagesForUser(c *gin.Context) {
	username := c.Param("username")
	packages, err := service.GetPackagesForUser(username)
	if err != nil {
		fmt.Println("GetPackagesForUser error:", err) //debug

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, packages)
}

// GET /api/package/info/:packageID
func GetPackageInfo(c *gin.Context) {
	packageID := c.Param("packageID")
	pack, err := service.GetPackageInfo(packageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Package not found"})
		return
	}
	c.JSON(http.StatusOK, pack)
}

// PUT /api/package/destination
func ChangePackageDestination(c *gin.Context) {
	var req struct {
		PackageID string           `json:"package_id"`
		NewCoord  model.Coordinate `json:"coord"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	msg, err := service.ChangePackageDestination(req.PackageID, req.NewCoord)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

// POST /api/package/create
func CreatePackage(c *gin.Context) {
	var req struct {
		PackageID   string `json:"package_id" binding:"required"`
		Username    string `json:"username" binding:"required"`
		Items       string `json:"items" binding:"required"`
		DestX       int    `json:"destination_x" binding:"required"`
		DestY       int    `json:"destination_y" binding:"required"`
		WarehouseID uint   `json:"warehouse_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	err := service.CreatePackage(req.PackageID, req.Username, req.Items, req.DestX, req.DestY, req.WarehouseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Package created successfully"}) //Status 201 for request succeed and a new resource has been created
}

// PUT /api/package/assign-truck
func LinkTruckToPackage(c *gin.Context) {
	var req struct {
		PackageID string `json:"package_id" binding:"required"`
		TruckID   uint   `json:"truck_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	err := service.LinkTruckToPackage(req.PackageID, req.TruckID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Truck assigned to package"})
}

// PUT /api/package/status
func ChangePackageStatus(c *gin.Context) {
	var req struct {
		PackageID string `json:"package_id"`
		Status    string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	err := service.ChangePackageStatus(req.PackageID, model.PackageStatus(req.Status))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Package status updated"})
}

// GET /api/package/warehouse/:packageID
func GetWarehouseID(c *gin.Context) {
	packageID := c.Param("packageID")
	id, err := service.GetWarehouseID(packageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"warehouse_id": id})
}

// Helper function to parse uint from path param
func parseUintParam(c *gin.Context, name string) (uint, error) {
	valStr := c.Param(name)
	var val uint
	_, err := fmt.Sscanf(valStr, "%d", &val)
	return val, err
}
