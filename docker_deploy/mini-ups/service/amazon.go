package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mini-ups/config"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//UPDATE HERE: Make sure to move the listening to response from Service UPS to here!
//We send Amazon POST request -> We listen to their response

var amazonURL = "http://" + config.AppConfig.Amazon.Host + ":" + strconv.Itoa(config.AppConfig.Amazon.Port) + "/api/ups/"

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

func SendWorldIDToAmazon(worldID int) error {
	msg := gin.H{
		"action":     "world_created",
		"message_id": uuid.New().String(),
		"world_id":   worldID,
	}
	return sendAmazonPost(msg)
}

// helper function to send Post Request
func sendAmazonPost(payload map[string]interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal payload: %v", err)
		return err
	}

	log.Printf("Marshalled JSON payload: %s", string(data))

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(amazonURL, "application/json", bytes.NewBuffer(data)) // Change URL
	if err != nil {
		log.Printf("HTTP POST failed: %v", err)
		return err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading body:", err)
	} else {
		log.Println("Response body:", string(bodyBytes))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("received non-OK response from Amazon: %d", resp.StatusCode)
		return errors.New("non-OK response from Amazon")
	}

	log.Printf("Sent to Amazon: %s\n", payload["action"])
	return nil
}
