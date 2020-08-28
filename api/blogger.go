package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewPost(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func EditPost(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func DeletePost(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func GetPosts(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func GetPost(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func NewComment(c *gin.Context) {
}

func CommentList(c *gin.Context) {
}

func BloggerService(router *gin.Engine) *gin.Engine {
	account := router.Group("/blog")

	account.Use(validateSessionToken)

	account.PUT("/", NewPost)
	account.POST("/", EditPost)
	account.GET("/", GetPosts)
	account.GET("/:postid", GetPost)
	account.DELETE("/delete", DeletePost)

	return router
}
