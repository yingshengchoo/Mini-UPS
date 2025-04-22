package controller

import (
	"fmt"
	"log"
	"mini-ups/model"
	"mini-ups/service"
	"mini-ups/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

//ERROR TYPES 文檔有specify 可能要檢查有沒有用對

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
	case "loading_package":
		{
			LoadingPackage(c)
			break
		}
	case "package_loaded":
		{
			LoadedPackage(c) //handle their request
			Deliver(c)       //notify them
			break
		}
	case "query_status":
		{
			CheckStatus(c)
			break
		}
	case "truck_arrived_response":
		{
			HandleTruckArrivedResponse(c)
			break
		}
	case "package_delivered_response":
		{
			HandlePackageDelieveredResponse(c)
			break
		}
	case "query_status_response":
		{
			HandlePackageQueryStatusResponse(c)
			break
		}
	case "redirect_pacakage_response":
		{
			HandleRedirectPackageResponse(c)
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
	var req struct {
		PackageID   string         `json:"package_id" binding:"required"`
		Username    string         `json:"username" binding:"required"`
		Items       datatypes.JSON `json:"items" binding:"required"`
		DestX       int            `json:"destination_x" binding:"required"`
		DestY       int            `json:"destination_y" binding:"required"`
		WarehouseID uint           `json:"warehouse_id" binding:"required"`
		MessageID   string         `json:"message_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"action":         "pickup_response",
			"in_response_to": req.MessageID,
			"status":         "error",
			"message":        "Invalid input",
		})
		return
	}
	packageID, err := service.CreatePackage(req.PackageID, req.Username, req.Items, req.DestX, req.DestY, req.WarehouseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// goroutine send world pickup
	//要是我們沒有truck 這裡因該放在一個queue裡
	// 還是就回個error 等他們重新send pickup_request? <-- 目前是這樣
	truck, err := service.GetFirstIdleTruck()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"action":         "pickup_response",
			"in_response_to": req.MessageID,
			"status":         "error",
			"message":        "No available truck",
		})
		return
	}

	truckID, err := service.GetIDByTruck(truck)
	//因該不會有ERROR
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"action":         "pickup_response",
			"in_response_to": req.MessageID,
			"status":         "error",
			"message":        err.Error(),
		})
		return
	}

	service.LinkTruckToPackage(string(packageID), uint(truckID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"action":          "pickup_response",
			"in_response_to":  req.MessageID,
			"tracking_number": packageID,
			"status":          "error",
			"message":         err.Error(),
		})
		return
	}

	seqnum := util.GenerateSeqNum() //  <------ 還沒implement丟包 所以沒有存 seqnum request pair. 現在只是assign seqnum 而已
	//TODO: Implement seqnum, request 儲存 for 掉包

	err = service.SendWorldRequestToGoPickUp(truckID, req.WarehouseID, seqnum) //Do World Command // 用thread?
	if err != nil {
		log.Println("Error sending GoPickUp command:", err)
	}

	//respond to Amazon
	resp := gin.H{
		"action":          "pickup_response",
		"in_response_to":  req.MessageID,
		"tracking_number": packageID,
		"status":          "success",
		"message":         fmt.Sprintf("Package %s marked as ready", req.PackageID),
	}
	c.JSON(http.StatusOK, resp)
}

// seems unnecessary
// POST /api/ups/package-ready
func RespondPackageReady(c *gin.Context) {
	// TODO implement

	//看不懂這是用來做什麼

}

// POST /api/ups/load
func LoadingPackage(c *gin.Context) {
	// TODO implement
	// just update the status info
	var req struct {
		MessageID   string `json:"message_id" binding:"required"`
		PackageID   string `json:"package_id" binding:"required"`
		TruckID     int    `json:"truck_id" binding:"required"`
		WarehouseID int    `json:"warehouse_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"action":         "loading_package_response",
			"message_id":     uuid.New().String(),
			"in_response_to": req.MessageID,
			"status":         "error",
			"message":        "Invalid request format: " + err.Error(),
		})
		return
	}

	err := service.ChangePackageStatus(req.PackageID, model.PackageStatus(model.StatusLoading))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"action":         "loading_package_response",
			"message_id":     uuid.New().String(),
			"in_response_to": req.MessageID,
			"status":         "error",
			"message":        err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"action":         "loading_package_response",
		"message_id":     uuid.New().String(),
		"in_response_to": req.MessageID,
		"status":         "success",
		"message":        fmt.Sprintf("Package %s is now loaded onto truck %d", req.PackageID, req.TruckID),
	})

}

// TODO test
// POST /api/ups/status
func CheckStatus(c *gin.Context) {
	var req struct {
		Action    string `json:"action" binding:"required"`
		Timestamp string `json:"timestamp" binding:"required"`
		MessageID string `json:"message_id" binding:"required"`
		PackageID string `json:"package_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"action":         "query_status_response",
			"timestamp":      time.Now().UTC().Format(time.RFC3339),
			"message_id":     uuid.New().String(),
			"in_response_to": req.MessageID,
			"status":         "error",
			"message":        "Invalid request format: " + err.Error(),
		})
		return
	}

	// Query package from DB
	pkg, err := service.GetPackageInfo(req.PackageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"action":         "query_status_response",
			"timestamp":      time.Now().UTC().Format(time.RFC3339),
			"message_id":     uuid.New().String(),
			"in_response_to": req.MessageID,
			"status":         "error",
			"message":        "Package not found",
		})
		return
	}

	// Default truck info
	var truckID int
	var truckStatus string
	var truckLocation = gin.H{"x": nil, "y": nil}

	if pkg.TruckID != nil {
		truck, err := service.GetTruckByID(*pkg.TruckID)
		if err == nil {
			truckID = int(truck.ID)
			truckStatus = string(truck.Status)
			truckLocation = gin.H{
				"x": truck.Coord.X,
				"y": truck.Coord.Y,
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"action":         "query_status_response",
		"timestamp":      time.Now().UTC().Format(time.RFC3339),
		"message_id":     uuid.New().String(),
		"in_response_to": req.MessageID,
		"status":         "success",
		"package_status": string(pkg.Status),
		"truck_id":       truckID,
		"truck_status":   truckStatus,
		"truck_location": truckLocation,
		"message":        fmt.Sprintf("Status for package %s retrieved", req.PackageID),
	})
}

func LoadedPackage(c *gin.Context) {
	var req struct {
		MessageID string `json:"message_id" binding:"required"`
		PackageID string `json:"package_id" binding:"required"`
		TruckID   int    `json:"truck_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"action":         "delivery_started", // 不知道這裡是不是該這樣寫
			"message_id":     uuid.New().String(),
			"in_response_to": req.MessageID,
			"status":         "error",
			"message":        "Invalid request format: " + err.Error(),
		})
		return
	}

	err := service.ChangeTruckStatus(req.TruckID, model.TruckStatus.LOADED)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	err = service.ChangePackageStatus(req.PackageID, model.PackageStatus(model.StatusLoaded))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"action":         "loading_package_response",
			"message_id":     uuid.New().String(),
			"in_response_to": req.MessageID,
			"status":         "error",
			"message":        err.Error(),
		})
		return
	}

	resp := gin.H{
		"action":         "package_delivered_response",
		"message_id":     uuid.New().String(),
		"in_response_to": req.MessageID,
		"status":         "success",
		"message":        fmt.Sprintf("Package %s marked as ready", req.PackageID),
	}
	c.JSON(http.StatusOK, resp)
}

// deliever packages (a truck)
// POST /api/ups/deliver
func Deliver(c *gin.Context) {
	var req struct {
		MessageID string `json:"message_id" binding:"required"`
		PackageID string `json:"package_id" binding:"required"`
		TruckID   int    `json:"truck_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"action":         "delivery_started",
			"message_id":     uuid.New().String(),
			"in_response_to": req.MessageID,
			"status":         "error",
			"message":        "Invalid request format: " + err.Error(),
		})
		return
	}

	err := service.ChangeTruckStatus(req.TruckID, model.TruckStatus.DELIVERING)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	err = service.ChangePackageStatus(req.PackageID, model.PackageStatus(model.StatusOutForDelivery))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"action":         "loading_package_response",
			"message_id":     uuid.New().String(),
			"in_response_to": req.MessageID,
			"status":         "error",
			"message":        err.Error(),
		})
		return
	}
	seqnum := util.GenerateSeqNum() //  <------ 還沒implement丟包 所以沒有存 seqnum request pair. 現在只是assign seqnum 而已
	service.SendWorldDeliveryRequest(req.PackageID, seqnum)
	//when world responds with UFinish, notify Amazon <-- happens in the ParseWorldResponse(controller - world.go)
}

// Simply checks that the response is in the right format.
func HandleActionResponse(c *gin.Context) {
	var req struct {
		MessageID    string `json:"message_id" binding:"required"`
		Status       string `json:"status" binding:"required"`
		InResponseTo string `json:"in_response_to" binding:"required"`
		Message      string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON format"})
		return
	}
}

func HandleTruckArrivedResponse(c *gin.Context) {
	HandleActionResponse(c)
	//SUCCESS 需要做什麼嗎

	//如果Error 重寄?
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if req.Status == "error" {
		//Missing truckID + warehouseID
		//service.NotifyAmazonTruckArrived() <-- 如果重寄我們需要 存我們send 的JSON?
	}
}

func HandlePackageDelieveredResponse(c *gin.Context) {
	HandleActionResponse(c)
	//跟上面一樣
}
func HandleRedirectPackageResponse(c *gin.Context) {
	HandleActionResponse(c)
	//跟上面一樣
}

// 沒有用到
func HandlePackageQueryStatusResponse(c *gin.Context) {
	var req struct {
		Action        string `json:"action" binding:"required"`
		Timestamp     string `json:"timestamp" binding:"required"`
		MessageID     string `json:"message_id" binding:"required"`
		InResponseTo  string `json:"in_response_to" binding:"required"`
		Status        string `json:"status" binding:"required"`
		PackageStatus string `json:"package_status" binding:"required"`
		TruckID       int    `json:"truck_id" binding:"required"`
		TruckStatus   string `json:"truck_status" binding:"required"`
		TruckLocation struct {
			X int `json:"x" binding:"required"`
			Y int `json:"y" binding:"required"`
		} `json:"truck_location" binding:"required"`
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON format"})
		return
	}
	//跟上面一樣
}
