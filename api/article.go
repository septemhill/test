package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/septemhill/test/db"
	"github.com/septemhill/test/module"
)

func NewPost(c *gin.Context) {
	art := new(module.Article)

	requestHandler(c, art, func(ctx context.Context, db *db.DB, v interface{}) error {
		art := v.(*module.Article)
		return module.NewPost(c, db, *art)
	})
}

func EditPost(c *gin.Context) {
	art := new(module.Article)

	requestHandler(c, art, func(ctx context.Context, db *db.DB, v interface{}) error {
		art := v.(*module.Article)
		return module.EditPost(c, db, *art)
	})
}

func DeletePost(c *gin.Context) {
	art := new(module.Article)

	requestHandler(c, art, func(ctx context.Context, db *db.DB, v interface{}) error {
		art := v.(*module.Article)
		return module.DeletePost(c, db, *art)
	})
}

func GetPosts(c *gin.Context) {
	pi := paginator{Size: 10, Offset: 0, Ascend: false}
	if err := c.BindQuery(&pi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	db := PostgresDB(c)
	accs, err := module.GetPosts(c, db, pi.Size, pi.Offset, pi.Ascend)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": accs})
}

func GetPost(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func NewComment(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func GetComments(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func DeleteComment(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func ArticleService(router *gin.Engine) *gin.Engine {
	article := router.Group("/article")

	//blog.Use(validateSessionToken)

	article.PUT("/", NewPost)
	article.POST("/", EditPost)
	article.GET("/", GetPosts)
	article.GET("/:postid", GetPost)
	article.DELETE("/", DeletePost)

	article.PUT("/:postid/comment", NewComment)
	article.GET("/:postid/comment", GetComments)
	article.DELETE("/:postid/comment/:commentid", DeleteComment)

	return router
}
