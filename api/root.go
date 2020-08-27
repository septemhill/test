package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/septemhill/test/module"
	"github.com/sirupsen/logrus"
)

type rootHandler struct {
}

func sessionTokenGenerate() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (h *rootHandler) Login(c *gin.Context) {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}

	m := make(map[string]string)
	if err := json.Unmarshal(b, &m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}

	logger := Logger(c)
	rdb := PostgresDB(c)
	mdb := RedisDB(c)

	logger.WithFields(logrus.Fields{
		"username": m["username"],
		"password": m["password"],
	}).Debugln("User Info")

	row := rdb.QueryRowxContext(c, `SELECT COUNT(*) FROM account_private WHERE username = $1 AND password = $2`, m["username"], m["password"])
	n := 0
	if err := row.Scan(&n); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMessage": err.Error(),
		})
		return
	}

	if n != 1 {
		c.JSON(http.StatusOK, gin.H{
			"errorMessage": "Invalid username or password",
		})
		return
	}

	token := sessionTokenGenerate()
	mdb.Set(token, "", time.Hour*1)

	c.JSON(http.StatusOK, gin.H{
		"session_token": token,
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

func RootService(r *gin.Engine) *gin.Engine {
	handler := rootHandler{}
	root := r.Group("/")

	root.POST("/login", handler.Login)
	root.GET("/logout", validateSessionToken, handler.Logout)
	root.POST("/signup", handler.Signup)
	root.GET("/verify", handler.VerifyUserRegistration)

	return r
}
