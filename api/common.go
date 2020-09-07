package api

import (
	"bytes"
	"context"
	"crypto/sha512"
	"database/sql"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/jackc/pgconn"
	"github.com/septemhill/test/db"
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

func httpErrHandler(c *gin.Context, err error) {
	// errors from pgx driver
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errMessage": pgErr.Code + ":" + pgErr.Error(),
		})
		return
	}

	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{
			"message": "successful",
		})
		return
	}

	// Unknown error type
	c.JSON(http.StatusInternalServerError, gin.H{
		"errMessage": err.Error(),
	})
}

func requestHandler(c *gin.Context, v interface{}, handle reqAction) {
	if err := c.BindJSON(v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errMessage": err.Error(),
		})
		return
	}

	db := PostgresDB(c)
	redis := RedisDB(c)

	if err := handle(c, db, redis, v); err != nil {
		httpErrHandler(c, err)
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
