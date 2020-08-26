package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/septemhill/test/db"
	"github.com/septemhill/test/module"
)

func CreateAccount(c *gin.Context) {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMessage": err.Error(),
		})
		return
	}
	defer c.Request.Body.Close()

	var acc module.Account
	if err := json.Unmarshal(b, &acc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	db := c.MustGet(module.RESOURCE_RDB).(*db.DB)
	if err := module.CreateAccount(c, db, acc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "create successful"})
}

func GetAccountInfo(c *gin.Context) {
}

func UpdateAccountInfo(c *gin.Context) {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMessage": err.Error(),
		})
		return
	}
	defer c.Request.Body.Close()

	var acc module.Account
	if err := json.Unmarshal(b, &acc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	db := c.MustGet(module.RESOURCE_RDB).(*db.DB)
	if err := module.UpdateAccountInfo(c, db, acc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "update successful",
	})
}

func DeleteAccount(c *gin.Context) {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMessage": err.Error(),
		})
		return
	}
	defer c.Request.Body.Close()

	var acc module.Account
	if err := json.Unmarshal(b, &acc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	db := c.MustGet(module.RESOURCE_RDB).(*db.DB)
	if err := module.DeleteAccount(c, db, acc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "delete successful"})
}

func AccountService(router *gin.Engine) *gin.Engine {
	account := router.Group("/account")

	account.Use(validateSessionToken)
	account.PUT("/", CreateAccount)
	account.POST("/", UpdateAccountInfo)
	account.GET("/:user", GetAccountInfo)
	account.DELETE("/delete", DeleteAccount)

	return router
}
