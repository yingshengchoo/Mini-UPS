package router

import (
	"mini-ups/controller"
	"mini-ups/util"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	store := cookie.NewStore([]byte("mini-ups-secret"))

	router.Use(sessions.Sessions("mini-ups-session", store))
	store.Options(sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // 如果没有启用 HTTPS，可以设置为 false
		SameSite: http.SameSiteLaxMode,
		Domain:   "vcm-46755.vm.duke.edu",
	})

	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://vcm-46755.vm.duke.edu"}, // 允许的前端地址
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	// 	AllowCredentials: true, // 允许携带凭证（cookies）
	// }))

	router.Static("/static", "./frontend")
	router.Static("/home", "./frontend/home")
	router.Static("/login", "./frontend/login")
	router.Static("/register", "./frontend/register")
	router.Static("/devtool", "./frontend/devtool")
	router.LoadHTMLGlob("frontend/home/home.html")

	router.GET("/users/:username", controller.GetUserByUsername)
	router.GET("/share/:packageID", controller.GetShareInfo)
	router.GET("/share/upshost", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"upshost": util.UPS_HOST})
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", nil)
	})

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

		//package api
		packageGroup := apiGroup.Group("/package")
		{
			packageGroup.GET("/user/:username", controller.GetPackagesForUser)
			packageGroup.GET("/info/:packageID", controller.GetPackageInfo)
			packageGroup.PUT("/destination", controller.ChangePackageDestination)
			packageGroup.POST("/create", controller.CreatePackage)
			packageGroup.PUT("/assign-truck", controller.LinkTruckToPackage)
			packageGroup.PUT("/status", controller.ChangePackageStatus)
			packageGroup.GET("/warehouse/:packageID", controller.GetWarehouseID)
			packageGroup.POST("/redirect", controller.RedirectPackage)
			packageGroup.POST("/prioritize/:packageID", controller.PrioritizePackage)
		}

		// ups side api for amazon
		// parse json's action to different further api
		amazonGroup := apiGroup.Group("/amazon")
		{
			amazonGroup.POST("/", controller.ParseAction)
			amazonGroup.POST("/pickup", controller.PickUp)
			amazonGroup.POST("/package-ready", controller.RespondPackageReady)
			amazonGroup.POST("/load", controller.LoadingPackage)
			amazonGroup.POST("/deliver", controller.Deliver)
			amazonGroup.GET("/status", controller.CheckStatus)
		}
	}
	return router
}
