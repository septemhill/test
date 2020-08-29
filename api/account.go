package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/septemhill/test/db"
	"github.com/septemhill/test/module"
)

func CreateAccount(c *gin.Context) {
	acc := new(module.Account)

	requestHandler(c, acc, func(ctx context.Context, db *db.DB, v interface{}) error {
		acc := v.(*module.Account)
		return module.CreateAccount(c, db, *acc)
	})
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
	acc := new(module.Account)

	requestHandler(c, acc, func(ctx context.Context, db *db.DB, v interface{}) error {
		acc := v.(*module.Account)
		return module.UpdateAccountInfo(c, db, *acc)
	})
}

func DeleteAccount(c *gin.Context) {
	acc := new(module.Account)

	requestHandler(c, acc, func(ctx context.Context, db *db.DB, v interface{}) error {
		acc := v.(*module.Account)
		return module.DeleteAccount(c, db, *acc)
	})
}

func AccountService(router *gin.Engine) *gin.Engine {
	account := router.Group("/account")

	account.Use(validateSessionToken)
	account.PUT("/", CreateAccount)
	account.POST("/", UpdateAccountInfo)
	account.GET("/:user", GetAccountInfo)
	account.DELETE("/", DeleteAccount)

	return router
}
