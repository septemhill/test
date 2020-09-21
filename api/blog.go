package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type blogHandler struct{}

func (h *blogHandler) Nothing(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func BlogService(r gin.IRouter) gin.IRouter {
	handler := blogHandler{}

	blog := r.Group("/blog/:user")

	// Load `Article` API under `Blog`
	ArticleService(blog)

	blog.GET("/", handler.Nothing)

	return r
}
