package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/septemhill/test/db"
)

type paginator struct {
	Size   int  `form:"size"`
	Offset int  `form:"offset"`
	Ascend bool `form:"asc"`
}

type reqAction func(ctx context.Context, db *db.DB, v interface{}) error

func requestHandler(c *gin.Context, v interface{}, handle reqAction) {
	if err := c.ShouldBindJSON(v); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	db := PostgresDB(c)

	if err := handle(c, db, v); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successful",
	})
}
