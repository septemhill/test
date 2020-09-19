package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	"github.com/septemhill/test/middleware"
	"github.com/septemhill/test/module"
)

type articleHandler struct{}

func (h *articleHandler) NewPost(c *gin.Context) {
	art := new(module.Article)
	if err := c.BindJSON(art); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	d := middleware.PostgresDB(c)
	requestHandler(c, func(ctx context.Context) (interface{}, error) {
		return module.NewPost(c, d, art)
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

func (h *articleHandler) EditPost(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	art := new(module.Article)
	if err := c.BindJSON(art); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	art.ID = id
	d := middleware.PostgresDB(c)

	requestHandler(c, func(ctx context.Context) (interface{}, error) {
		return module.EditPost(ctx, d, art)
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

func (h *articleHandler) DeletePost(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	d := middleware.PostgresDB(c)

	requestHandler(c, func(ctx context.Context) (interface{}, error) {
		return module.DeletePost(c, d, id)
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

func (h *articleHandler) GetPosts(c *gin.Context) {
	user := c.Param("user")
	pi := paginator{Size: 10, Offset: 0, Ascend: false}
	if err := c.BindQuery(&pi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	d := middleware.PostgresDB(c)

	requestHandler(c, func(ctx context.Context) (interface{}, error) {
		return module.GetPosts(c, d, user, pi.Size, pi.Offset, pi.Ascend)
	}, func(c *gin.Context, err error) {
		var pgerr pgx.PgError
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, nil)
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

func (h *articleHandler) GetPost(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	d := middleware.PostgresDB(c)

	requestHandler(c, func(ctx context.Context) (interface{}, error) {
		return module.GetPost(c, d, id)
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

func (h *articleHandler) NewComment(c *gin.Context) {
	id, comment := new(module.Comment), new(module.Comment)
	if err := c.BindUri(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	if err := c.BindJSON(comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	d := middleware.PostgresDB(c)
	comment.ArticleID = id.ArticleID

	requestHandler(c, func(ctx context.Context) (interface{}, error) {
		return module.NewComment(c, d, comment)
	}, func(c *gin.Context, err error) {
		var pgerr pgx.PgError
		if err == sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		if errors.As(err, &pgerr) {
			c.JSON(http.StatusNotFound, gin.H{
				"errMessage": pgerr.Code + ":" + pgerr.Error(),
			})
			return
		}
	})
}

func (h *articleHandler) UpdateComment(c *gin.Context) {
	ids, comment := new(module.Comment), new(module.Comment)
	if err := c.BindUri(&ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	if err := c.BindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	comment.ArticleID = ids.ArticleID
	comment.ID = ids.ID
	d := middleware.PostgresDB(c)

	requestHandler(c, func(ctx context.Context) (interface{}, error) {
		return module.UpdateComment(c, d, comment)
	}, func(c *gin.Context, err error) {
		var pgerr pgx.PgError
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, err.Error())
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

func (h *articleHandler) GetComments(c *gin.Context) {
	art := new(module.Article)
	if err := c.BindUri(art); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	pi := paginator{Size: 10, Offset: 0, Ascend: false}
	if err := c.BindQuery(&pi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	d := middleware.PostgresDB(c)

	requestHandler(c, func(ctx context.Context) (interface{}, error) {
		return module.GetComments(c, d, art.ID, pi.Size, pi.Offset)
	}, func(c *gin.Context, err error) {
		var pgerr pgx.PgError
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, nil)
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

func (h *articleHandler) DeleteComment(c *gin.Context) {
	comment := new(module.Comment)
	if err := c.BindUri(comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	d := middleware.PostgresDB(c)

	requestHandler(c, func(ctx context.Context) (interface{}, error) {
		return module.DeleteComment(c, d, comment)
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

func ArticleService(r gin.IRouter) gin.IRouter {
	handler := articleHandler{}
	article := r.Group("/article")

	article.Use(middleware.ValidateSessionToken)

	article.POST("/", handler.NewPost)
	article.PUT("/:id", handler.EditPost)
	article.GET("/", handler.GetPosts)
	article.GET("/:id", handler.GetPost)
	article.DELETE("/:id", handler.DeletePost)

	article.POST("/:id/comment", handler.NewComment)
	article.GET("/:id/comments", handler.GetComments)
	article.DELETE("/:id/comment/:commentid", handler.DeleteComment)
	article.PUT("/:id/comment/:commentid", handler.UpdateComment)

	return r
}
