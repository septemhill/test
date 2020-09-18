package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/septemhill/test/middleware"
	"github.com/septemhill/test/module"
)

type accountHandler struct{}

func (h *accountHandler) CreateAccount(c *gin.Context) {
	acc := new(module.Account)

	if err := c.BindJSON(acc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	d := middleware.PostgresDB(c)

	requestHandler(c, func(ctx context.Context) (interface{}, error) {
		return module.CreateAccount(ctx, d, *acc)
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
	acc := new(module.Account)
	if err := c.BindUri(acc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	d := middleware.PostgresDB(c)

	requestHandler(c, func(ctx context.Context) (interface{}, error) {
		return module.GetAccountInfo(ctx, d, acc)
	}, func(c *gin.Context, err error) {
		var pgerr pgx.PgError
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
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

func (h *accountHandler) UpdateAccountInfo(c *gin.Context) {
	id, acc := new(module.Account), new(module.Account)
	if err := c.BindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	if err := c.BindJSON(&acc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	acc.ID = id.ID

	d := middleware.PostgresDB(c)

	requestHandler(c, func(ctx context.Context) (interface{}, error) {
		return module.UpdateAccountInfo(ctx, d, *acc)
	}, func(c *gin.Context, err error) {
		var pgerr pgx.PgError
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
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

func (h *accountHandler) DeleteAccount(c *gin.Context) {
	acc := new(module.Account)
	id, acc := new(module.Account), new(module.Account)
	if err := c.BindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	if err := c.BindJSON(&acc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	acc.ID = id.ID

	d := middleware.PostgresDB(c)

	requestHandler(c, func(ctx context.Context) (interface{}, error) {
		return module.DeleteAccount(ctx, d, *acc)
	}, func(c *gin.Context, err error) {
		var pgerr pgx.PgError
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
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

func (h *accountHandler) ChangePassword(c *gin.Context) {
	pass := password{}
	if err := c.BindJSON(&pass); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	d := middleware.PostgresDB(c)
	mail := middleware.UserEmail(c)

	requestHandler(c, func(ctx context.Context) (interface{}, error) {
		return module.ChangePassword(ctx, d, mail, pass.Password)
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

func AccountService(r gin.IRouter) gin.IRouter {
	handler := accountHandler{}
	account := r.Group("/account")

	account.Use(middleware.SetTestPostgresDB(), middleware.SetTestRedisDB())
	account.Use(middleware.ValidateSessionToken)

	account.POST("/", handler.CreateAccount)
	account.PUT("/:user", handler.UpdateAccountInfo)
	account.GET("/:user", handler.GetAccountInfo)
	account.DELETE("/:user", handler.DeleteAccount)
	account.PUT("/:user/chgpasswd", handler.ChangePassword)

	return r
}
