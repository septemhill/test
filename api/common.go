package api

import (
	"bytes"
	"context"
	"crypto/sha512"
	"encoding/gob"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
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

type reqAction func(ctx context.Context) (interface{}, error)

type errHandler func(c *gin.Context, err error)

func requestHandler(c *gin.Context, handle reqAction, errHandle errHandler) {
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
	c.JSON(http.StatusOK, v)
}
