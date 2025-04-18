package router

import (
	"mini-ups/controller"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/users/:username", controller.GetUserByUsername)
	return router
}
