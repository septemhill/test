package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jackc/pgx"
	"github.com/septemhill/test/db"
	"github.com/septemhill/test/middleware"
	"github.com/septemhill/test/module"
)

type accountHandler struct{}

func (h *accountHandler) CreateAccount(c *gin.Context) {
	acc := new(module.Account)

	requestHandler(c, acc, func(ctx context.Context, db *db.DB, redis *redis.Client, v interface{}) error {
		acc := v.(*module.Account)
		return module.CreateAccount(c, db, *acc)
	}, func(c *gin.Context, err error) {
		var pgerr pgx.PgError
		if err == sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errMessage": err.Error(),
			})
			return
		}
		if errors.As(err, &pgerr) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errMessage": pgerr.Code + ":" + pgerr.Error(),
			})
			return
		}
	})
}

func (h *accountHandler) GetAccountInfo(c *gin.Context) {
	username := c.Param("user")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": errors.New("unknown user"),
		})
		return
	}

	var acc module.Account
	acc.Username = username

	db := middleware.PostgresDB(c)
	if err := module.GetAccountInfo(c, db, &acc); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, nil)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, acc)
}

func (h *accountHandler) UpdateAccountInfo(c *gin.Context) {
	acc := new(module.Account)

	requestHandler(c, acc, func(ctx context.Context, db *db.DB, redis *redis.Client, v interface{}) error {
		acc := v.(*module.Account)
		return module.UpdateAccountInfo(c, db, *acc)
	}, func(c *gin.Context, err error) {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, nil)
		}
	})
}

func (h *accountHandler) DeleteAccount(c *gin.Context) {
	acc := new(module.Account)

	requestHandler(c, acc, func(ctx context.Context, db *db.DB, redis *redis.Client, v interface{}) error {
		acc := v.(*module.Account)
		return module.DeleteAccount(c, db, *acc)
	}, func(c *gin.Context, err error) {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, nil)
		}
	})
}

func (h *accountHandler) ChangePassword(c *gin.Context) {
	pass := password{}
	if err := c.BindJSON(&pass); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}
}

func AccountService(r gin.IRouter) gin.IRouter {
	handler := accountHandler{}
	account := r.Group("/account")

	account.Use(middleware.ValidateSessionToken)

	account.POST("/", handler.CreateAccount)
	account.PUT("/:user", handler.UpdateAccountInfo)
	account.GET("/:user", handler.GetAccountInfo)
	account.DELETE("/:user", handler.DeleteAccount)
	account.PUT("/:user/chgpasswd", handler.ChangePassword)

	return r
}
