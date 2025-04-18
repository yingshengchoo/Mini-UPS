package controller

import (
	"mini-ups/service"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// register request struct
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// register
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := service.RegisterUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "register success"})
}

// register
func Login(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := service.LoginUser(req.Username, req.Password, c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"login": false,
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"login":   true,
		"message": "Login successfully!",
	})
}

// get user info
func GetUserInfo(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("user")
	if username == nil {
		c.JSON(http.StatusOK, gin.H{"userlogined": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userlogined": true,
		"username":    username,
	})
}

func GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	user, err := service.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
