package api

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/septemhill/test/middleware"
	"github.com/septemhill/test/module"
	"github.com/septemhill/test/utils"
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

	db := middleware.PostgresDB(c)
	redis := middleware.RedisDB(c)

	code, err := module.Login(c, db, redis, acc.Email, acc.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
	})
}

func (h *rootHandler) Logout(c *gin.Context) {
	token := c.GetHeader(utils.HEADER_SESSION_TOKEN)
	rdb := middleware.RedisDB(c)
	rdb.Del(token)
}

func (h *rootHandler) Signup(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func (h *rootHandler) VerifyUserRegistration(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "unknown verification code",
		})
		return
	}

	db := middleware.PostgresDB(c)
	redis := middleware.RedisDB(c)

	if err := module.VerifyUserRegistration(c, db, redis, code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "verification successful",
	})
}

func (h *rootHandler) ForgetPassword(c *gin.Context) {
	mail := email{}
	if err := c.BindJSON(&mail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	db := middleware.PostgresDB(c)
	redis := middleware.RedisDB(c)
	mailer := middleware.Mailer(c)

	code, err := module.ForgetPassword(c, db, redis, mail.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	buff := bytes.NewBuffer(nil)
	t, err := template.New("").Parse(utils.ForgetPasswordLetterTemplate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	if err := t.Execute(buff, struct{ Code string }{code}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	//TODO: add channel to receive result
	go func() {
		_ = utils.SendMail(*mailer, utils.MailInfo{
			From:    "septemhill@gmail.com",
			To:      "septemhill@gmail.com",
			Subject: "Reset password email confirm",
			Body:    buff.String(),
		})
	}()

	c.JSON(http.StatusOK, nil)
}

func (h *rootHandler) ResetPassword(c *gin.Context) {
	type resetCode struct {
		Code string `form:"code"`
	}

	code := resetCode{}
	pass := password{}
	if err := c.BindQuery(&code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	if err := c.BindJSON(&pass); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	db := middleware.PostgresDB(c)
	redis := middleware.RedisDB(c)

	if err := module.ResetPassword(c, db, redis, code.Code, pass.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "reset successful",
	})
}

func RootService(r gin.IRouter) gin.IRouter {
	handler := rootHandler{}
	root := r.Group("/")

	root.POST("/login", handler.Login)
	root.GET("/logout", middleware.ValidateSessionToken, handler.Logout)
	root.POST("/signup", handler.Signup)
	root.GET("/verify", handler.VerifyUserRegistration)
	root.POST("/forget", handler.ForgetPassword)
	root.POST("/reset", handler.ResetPassword)

	return r
}
