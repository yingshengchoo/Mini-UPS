package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const amazonURL = "http://localhost:8000/api/amazon" //是這個ＵＲＬ嗎

// Send Request to Amazon, notifying that the Delivery was complete
func NotifyAmazonDeliveryComplete(packageID string, truckID, x, y int) error {
	msg := gin.H{
		"action":     "package_delivered",
		"message_id": uuid.New().String(),
		"package_id": packageID,
		"truck_id":   truckID,
		"delivery_x": x,
		"delivery_y": y,
	}

	return sendAmazonPost(msg)
}

// Send notification to Amazon notifying that the truck has arrived at warehouse
func NotifyAmazonTruckArrived(truckID, warehouseID int) error {
	msg := gin.H{
		"action":       "truck_arrived",
		"message_id":   uuid.New().String(),
		"truck_id":     truckID,
		"warehouse_id": warehouseID,
	}
	return sendAmazonPost(msg)
}

// SEnd Notification to Amazon notifying that a package has been redirected
func NotifyAmazonRedirectPacakge(packageID string, newX, newY, userID int) error {
	msg := gin.H{
		"action":            "redirect_package",
		"message_id":        uuid.New().String(),
		"package_id":        packageID,
		"new_destination_x": newX,
		"new_destination_y": newY,
		"user_id":           userID,
	}
	return sendAmazonPost(msg)
}

// 好像我們沒有用到但文當有寫我就寫了
// Request amazon to return query info on package of given packageID
func SendQueryStatusToAmazon(packageID string) error {
	msg := gin.H{
		"action":     "query_status",
		"message_id": uuid.New().String(),
		"package_id": packageID,
	}
	return sendAmazonPost(msg)
}

// helper function to send Post Request
func sendAmazonPost(payload map[string]interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(amazonURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("HTTP POST failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response from Amazon: %d", resp.StatusCode)
	}

	fmt.Printf("Sent to Amazon: %s\n", payload["action"])
	return nil
}
