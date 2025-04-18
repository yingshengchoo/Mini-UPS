package main

import (
	"mini-ups/db"
	"mini-ups/router"
)

func main() {
	db.InitDB()
	r := router.InitRouter()
	r.Run(":8080")
}
