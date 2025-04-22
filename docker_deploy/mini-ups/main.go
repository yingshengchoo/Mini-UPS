package main

import (
	"log"
	"mini-ups/controller"
	"mini-ups/db"
	"mini-ups/protocol"
	"mini-ups/protocol/worldupspb"
	"mini-ups/router"
	"mini-ups/util"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("initing service...")
	//initialize db
	db.InitDB()

	//use different connectUPS call based on who connects to world first
	//example trucks
	trucks := []*worldupspb.UInitTruck{
		protocol.CreateTruck(1, 100, 200),
		protocol.CreateTruck(2, 150, 250),
	}
	worldID := int64(1) //<- make this dynamic in the future
	util.UPSConn = protocol.ConnectUPSWithWorldID(worldID, trucks)
	// UPSConn, worldID = protocol.ConnectUPS(trucks)

	//world response listener
	go controller.ParseWorldResponse()

	//start router
	r := router.InitRouter()
	log.Println("start service...")
	r.Run(":8080")
}
