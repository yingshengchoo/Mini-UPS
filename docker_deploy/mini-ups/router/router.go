package router

import (
	"mini-ups/controller"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	store := cookie.NewStore([]byte("mini-ups-secret"))
	router.Use(sessions.Sessions("mini-ups-session", store))

	router.Static("/static", "./frontend")
	router.Static("/home", "./frontend/home")
	router.Static("/login", "./frontend/login")
	router.Static("/register", "./frontend/register")
	router.GET("/users/:username", controller.GetUserByUsername)

	apiGroup := router.Group("/api")
	{
		apiGroup.POST("/register", controller.Register)
		apiGroup.POST("/login", controller.Login)
		apiGroup.POST("/logout", controller.Logout)
		apiGroup.GET("/user/info", controller.GetUserInfo)
	}
	return router
}
