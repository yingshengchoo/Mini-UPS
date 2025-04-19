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
		// user api
		userGroup := apiGroup.Group("/user")
		{
			userGroup.POST("/register", controller.Register)
			userGroup.POST("/login", controller.Login)
			userGroup.POST("/logout", controller.Logout)
			userGroup.GET("/info", controller.GetUserInfo)
		}

		// truck api
		truckGroup := apiGroup.Group("/truck")
		{
			truckGroup.POST("/register", controller.RegisterTruck)
			truckGroup.GET("/info", controller.GetTruckInfo)
		}

	}
	return router
}
