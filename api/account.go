package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
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

	db := PostgresDB(c)
	if err := module.CreateAccount(c, db, acc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "create successful"})
}

func GetAccountInfo(c *gin.Context) {
	username := c.Param("user")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": errors.New("unknown user"),
		})
		return
	}

	var acc module.Account
	acc.Username = username

	db := PostgresDB(c)
	if err := module.GetAccountInfo(c, db, &acc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, acc)
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

	db := PostgresDB(c)
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

	db := PostgresDB(c)
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
