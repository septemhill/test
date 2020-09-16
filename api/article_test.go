package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/septemhill/test/module"
	test "github.com/septemhill/test/testing"
	"github.com/stretchr/testify/assert"
)

func TestNewPost(t *testing.T) {
	ctx := context.Background()
	ts := test.NewTestRouter(gin.Default(), ArticleService)
	d, r := test.NewTestDB()
	defer func() {
		d.Close()
		r.Close()
	}()

	user := test.NewAccount(ctx, d, false)

	tests := []struct {
		Description string
		Article     module.Article
		StatusCode  int
	}{
		{
			Description: "Create new post",
			Article: module.Article{
				Author:  user.Username,
				Title:   "This is title",
				Content: "This is content",
			},
			StatusCode: http.StatusOK,
		}, {
			Description: "Create another post with the same title and content",
			Article: module.Article{
				Author:  user.Username,
				Title:   "This is title",
				Content: "This is content",
			},
			StatusCode: http.StatusOK,
		}, {
			Description: "Create new post with different title and content",
			Article: module.Article{
				Author:  user.Username,
				Title:   "Unknown title",
				Content: "Unknown content",
			},
			StatusCode: http.StatusOK,
		},
	}

	defer func() {
		test.DeleteAccounts(ctx, d, user)
		test.DeletePostsByAccount(ctx, d, user.Username)
	}()

	asserts := assert.New(t)

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			b, err := json.Marshal(&test.Article)
			asserts.NoError(err)

			req, err := http.NewRequest("POST", ts.URL+"/article/", bytes.NewBuffer(b))
			asserts.NoError(err)

			rsp, err := http.DefaultClient.Do(req)
			asserts.NoError(err)
			defer rsp.Body.Close()

			body, err := ioutil.ReadAll(rsp.Body)
			asserts.NoError(err)

			asserts.Equal(test.StatusCode, rsp.StatusCode, string(body))
		})
	}
}

func TestGetPosts(t *testing.T) {
}

func TestGetPost(t *testing.T) {
	ctx := context.Background()
	ts := test.NewTestRouter(gin.Default(), ArticleService)
	d, r := test.NewTestDB()
	defer func() {
		d.Close()
		r.Close()
	}()

	users := []*module.Account{
		test.NewAccount(ctx, d, false),
		test.NewAccount(ctx, d, false),
		test.NewAccount(ctx, d, false),
	}

	posts := []*module.Article{
		test.NewPost(ctx, d, users[0].Username),
		test.NewPost(ctx, d, users[1].Username),
	}

	comments := []*module.Comment{
		test.NewComment(ctx, d, users[0].Username, posts[0].ID),
		test.NewComment(ctx, d, users[0].Username, posts[1].ID),
		test.NewComment(ctx, d, users[1].Username, posts[1].ID),
		test.NewComment(ctx, d, users[2].Username, posts[0].ID),
		test.NewComment(ctx, d, users[2].Username, posts[1].ID),
	}

	defer func() {
		test.DeleteAccounts(ctx, d, users...)
		test.DeleteComments(ctx, d, comments...)
		test.DeletePosts(ctx, d, posts...)
	}()

	tests := []struct {
		Description string
		Article     module.Article
		StatusCode  int
		Comments    int
	}{
		{
			Description: "Get post 1 with 2 comments",
			Article:     *posts[0],
			StatusCode:  http.StatusOK,
			Comments:    2,
		}, {
			Description: "Get post 2 with 3 comments",
			Article:     *posts[1],
			StatusCode:  http.StatusOK,
			Comments:    3,
		}, {
			Description: "Get non-exist post",
			Article: module.Article{
				ID: 5633,
			},
			StatusCode: http.StatusNotFound,
		},
	}

	asserts := assert.New(t)

	for _, test := range tests {
		req, err := http.NewRequest("GET", ts.URL+"/article/"+fmt.Sprint(test.Article.ID), nil)
		asserts.NoError(err)

		rsp, err := http.DefaultClient.Do(req)
		asserts.NoError(err)
		defer rsp.Body.Close()

		asserts.Equal(test.StatusCode, rsp.StatusCode)

		if test.StatusCode == http.StatusNotFound {
			return
		}

		body, err := ioutil.ReadAll(rsp.Body)
		asserts.NoError(err)

		article := new(module.Article)
		err = json.Unmarshal(body, &article)
		asserts.NoError(err)

		asserts.Equal(test.Comments, len(article.Comments))
	}
}

func TestEditPost(t *testing.T) {}

func TestDeletePost(t *testing.T) {}
