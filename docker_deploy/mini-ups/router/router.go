package router

import (
	"github.com/gin-gonic/gin"
	"mini-ups/controller"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.Static("/static", "./frontend")
	router.GET("/users/:username", controller.GetUserByUsername)
	return router
}
