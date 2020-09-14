package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/septemhill/test/module"
	test "github.com/septemhill/test/testing"
	"github.com/stretchr/testify/assert"
)

func TestNewPost(t *testing.T) {
	ts := test.NewTestRouter(gin.Default(), ArticleService)
	d, r := test.NewTestDB()
	defer func() {
		d.Close()
		r.Close()
	}()

	tests := []struct {
		Description string
		Article     module.Article
		StatusCode  int
		Clean       bool
	}{}

	defer func() {
	}()

	asserts := assert.New(t)

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			b, err := json.Marshal(&test.Article)
			asserts.NoError(err)

			req, err := http.NewRequest("DELETE", ts.URL+"/blog/article/", bytes.NewBuffer(b))
			asserts.NoError(err)

			rsp, err := http.DefaultClient.Do(req)
			asserts.NoError(err)
			defer rsp.Body.Close()

			asserts.Equal(test.StatusCode, rsp.StatusCode)
		})
	}
}

func TestGetPosts(t *testing.T) {}

func TestGetPost(t *testing.T) {}

func TestEditPost(t *testing.T) {}

func TestDeletePost(t *testing.T) {}
