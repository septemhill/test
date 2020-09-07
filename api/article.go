package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/septemhill/test/db"
	"github.com/septemhill/test/middleware"
	"github.com/septemhill/test/module"
)

type articleHandler struct{}

func (h *articleHandler) NewPost(c *gin.Context) {
	art := new(module.Article)

	requestHandler(c, art, func(ctx context.Context, db *db.DB, redis *redis.Client, v interface{}) error {
		art := v.(*module.Article)
		return module.NewPost(c, db, *art)
	})
}

func (h *articleHandler) EditPost(c *gin.Context) {
	art := new(module.Article)

	requestHandler(c, art, func(ctx context.Context, db *db.DB, redis *redis.Client, v interface{}) error {
		art := v.(*module.Article)
		return module.EditPost(c, db, *art)
	})
}

func (h *articleHandler) DeletePost(c *gin.Context) {
	art := new(module.Article)

	requestHandler(c, art, func(ctx context.Context, db *db.DB, redis *redis.Client, v interface{}) error {
		art := v.(*module.Article)
		return module.DeletePost(c, db, *art)
	})
}

func (h *articleHandler) GetPosts(c *gin.Context) {
	pi := paginator{Size: 10, Offset: 0, Ascend: false}
	if err := c.BindQuery(&pi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	db := middleware.PostgresDB(c)
	arts, err := module.GetPosts(c, db, pi.Size, pi.Offset, pi.Ascend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	sendResponse(c, arts)
}

func (h *articleHandler) GetPost(c *gin.Context) {
	art := module.Article{}
	if err := c.BindUri(&art); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	db := middleware.PostgresDB(c)
	arti, err := module.GetPost(c, db, art.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	sendResponse(c, arti)
}

func (h *articleHandler) NewComment(c *gin.Context) {
	comment := new(module.Comment)
	if err := c.BindUri(comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	requestHandler(c, comment, func(ctx context.Context, db *db.DB, redis *redis.Client, v interface{}) error {
		comment := v.(*module.Comment)
		return module.NewComment(c, db, comment.ArticleID, *comment)
	})
}

func (h *articleHandler) UpdateComment(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func (h *articleHandler) GetComments(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func (h *articleHandler) DeleteComment(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func ArticleService(r gin.IRouter) gin.IRouter {
	handler := articleHandler{}
	article := r.Group("/article")

	//article.Use(validateSessionToken)

	article.POST("/", handler.NewPost)
	article.PUT("/", handler.EditPost)
	article.GET("/", handler.GetPosts)
	article.GET("/:id", handler.GetPost)
	article.DELETE("/", handler.DeletePost)

	article.POST("/:id/comment", handler.NewComment)
	article.GET("/:id/comment", handler.GetComments)
	article.DELETE("/:id/comment/:commentid", handler.DeleteComment)
	article.PUT("/:id/comment/:commentid", handler.UpdateComment)

	return r
}
