package main

import (
	"log"
	"mini-ups/db"
	"mini-ups/router"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	db.InitDB()
	r := router.InitRouter()
	r.Run(":8080")
}
