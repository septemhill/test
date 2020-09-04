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
	"github.com/septemhill/test/errors"
)

type paginator struct {
	Size   int  `form:"size"`
	Offset int  `form:"offset"`
	Ascend bool `form:"asc"`
}

type reqAction func(ctx context.Context, db *db.DB, redis *redis.Client, v interface{}) error

func httpErrHandler(c *gin.Context, err errors.HttpError) {
	c.JSON(err.Code(), gin.H{
		"errMessage": err.Error(),
	})
}

func requestHandler(c *gin.Context, v interface{}, handle reqAction) {
	if err := c.ShouldBindJSON(v); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	db := PostgresDB(c)
	redis := RedisDB(c)

	if err := handle(c, db, redis, v); err != nil {
		httpErrHandler(c, err.(errors.HttpError))
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

	if c.GetHeader(HEADER_IF_NOT_MATCH) == eTag {
		c.JSON(http.StatusNotModified, nil)
		return
	}

	c.Header(HEADER_ETAG, eTag)
	c.JSON(http.StatusOK, gin.H{
		"data": v,
	})
}
