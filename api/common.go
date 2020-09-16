package api

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/gob"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/septemhill/test/db"
	"github.com/septemhill/test/middleware"
	"github.com/septemhill/test/utils"
)

type password struct {
	Password string `json:"new_password"`
}

type paginator struct {
	Size   int  `form:"size"`
	Offset int  `form:"offset"`
	Ascend bool `form:"asc"`
}

type email struct {
	Email string `json:"email"`
}

type reqAction func(ctx context.Context, db *db.DB, redis *redis.Client, v interface{}) error
type reqAction2 func(ctx context.Context) (interface{}, error)

type errHandler func(c *gin.Context, err error)

func requestHandler2(c *gin.Context, handle reqAction2, errHandle errHandler) {
	v, err := handle(c)
	if err != nil {
		errHandle(c, err)
		return
	}

	if v == nil {
		c.JSON(http.StatusOK, nil)
		return
	}

	sendResponse(c, v)
}

func requestHandler(c *gin.Context, v interface{}, handle reqAction, errHandle errHandler) {
	if err := c.BindJSON(v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	d := middleware.PostgresDB(c)
	r := middleware.RedisDB(c)

	if err := handle(c, d, r, v); err != nil {
		errHandle(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "successful",
	})
}

func eTagCompute(v interface{}) (string, error) {
	buff := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buff)
	if err := enc.Encode(v); err != nil {
		return "", err
	}

	b := sha512.Sum512(buff.Bytes())
	return hex.EncodeToString(b[:]), nil
}

func sendResponse(c *gin.Context, v interface{}) {
	eTag, err := eTagCompute(v)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	if c.GetHeader(utils.HEADER_IF_NOT_MATCH) == eTag {
		c.JSON(http.StatusNotModified, nil)
		return
	}

	c.Header(utils.HEADER_ETAG, eTag)
	c.JSON(http.StatusOK, gin.H{
		"data": v,
	})
}
