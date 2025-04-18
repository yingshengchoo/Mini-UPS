package router

import (
	"github.com/gin-gonic/gin"
	"mini-ups/controller"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.Static("/static", "./frontend")
	router.Static("/home", "./frontend/home")
	router.Static("/login", "./frontend/login")
	router.Static("/register", "./frontend/register")
	router.GET("/users/:username", controller.GetUserByUsername)

	apiGroup := router.Group("/api")
	{
		apiGroup.POST("/register", controller.Register)
	}
	return router
}
