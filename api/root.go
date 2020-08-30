package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/septemhill/test/module"
	"github.com/sirupsen/logrus"
)

type rootHandler struct{}

func (h *rootHandler) Login(c *gin.Context) {
	var acc module.Account

	if err := c.BindJSON(&acc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	db := PostgresDB(c)
	redis := RedisDB(c)

	code, err := module.Login(c, db, redis, acc.Username, acc.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": gin.H{
			"code": code,
		},
	})
}

func (h *rootHandler) Logout(c *gin.Context) {
	token := c.GetHeader(HEADER_SESSION_TOKEN)
	rdb := RedisDB(c)
	rdb.Del(token)
}

func (h *rootHandler) Signup(c *gin.Context) {
	f := logrus.Fields{}

	c.JSON(http.StatusOK, f)
}

func (h *rootHandler) VerifyUserRegistration(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "unknown verification code",
		})
		return
	}

	db := PostgresDB(c)
	redis := RedisDB(c)

	if err := module.VerifyUserRegistration(c, db, redis, code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "verification successful",
	})
}

func RootService(r gin.IRouter) gin.IRouter {
	handler := rootHandler{}
	root := r.Group("/")

	root.POST("/login", handler.Login)
	root.GET("/logout", validateSessionToken, handler.Logout)
	root.POST("/signup", handler.Signup)
	root.GET("/verify", handler.VerifyUserRegistration)

	return r
}
