package main

import (
	"log"
	"mini-ups/db"
	"mini-ups/router"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("initing service...")
	db.InitDB()
	r := router.InitRouter()
	log.Println("start service...")
	r.Run(":8080")
}
