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
	"github.com/septemhill/test/utils"
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
	tk := test.NewTestSessionToken(r)

	user := test.NewAccount(ctx, d, false)
	header := map[string]string{
		utils.HEADER_SESSION_TOKEN: tk,
	}

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

			req, err := NewRequestWithTestHeader("POST", ts.URL+"/article/", bytes.NewBuffer(b), header)
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
	ctx := context.Background()
	ts := test.NewTestRouter(gin.Default(), BlogService)
	d, r := test.NewTestDB()
	defer func() {
		d.Close()
		r.Close()
	}()
	tk := test.NewTestSessionToken(r)

	users := []*module.Account{
		test.NewAccount(ctx, d, false),
		test.NewAccount(ctx, d, false),
	}

	posts := []*module.Article{
		test.NewPost(ctx, d, users[0].Username),
		test.NewPost(ctx, d, users[0].Username),
		test.NewPost(ctx, d, users[0].Username),
		test.NewPost(ctx, d, users[0].Username),
		test.NewPost(ctx, d, users[0].Username),
		test.NewPost(ctx, d, users[0].Username),
		test.NewPost(ctx, d, users[1].Username),
		test.NewPost(ctx, d, users[1].Username),
		test.NewPost(ctx, d, users[1].Username),
		test.NewPost(ctx, d, users[1].Username),
	}

	defer func() {
		test.DeletePosts(ctx, d, posts...)
		test.DeleteAccounts(ctx, d, users...)
	}()

	header := map[string]string{
		utils.HEADER_SESSION_TOKEN: tk,
	}

	tests := []struct {
		Description string
		Account     module.Account
		StatusCode  int
		Expected    []module.Article
	}{
		{
			Description: "Get user[0] posts, should have 6",
			Account:     *users[0],
			StatusCode:  http.StatusOK,
			Expected:    []module.Article{*posts[5], *posts[4], *posts[3], *posts[2], *posts[1], *posts[0]},
		},
		{
			Description: "Get user[1] posts, should have 4",
			Account:     *users[1],
			StatusCode:  http.StatusOK,
			Expected:    []module.Article{*posts[9], *posts[8], *posts[7], *posts[6]},
		},
	}

	asserts := assert.New(t)

	for _, test := range tests {
		req, err := NewRequestWithTestHeader("GET", ts.URL+"/blog/"+fmt.Sprint(test.Account.Username)+"/article/", nil, header)
		asserts.NoError(err)

		rsp, err := http.DefaultClient.Do(req)
		asserts.NoError(err)
		defer rsp.Body.Close()

		asserts.Equal(test.StatusCode, rsp.StatusCode)

		body, err := ioutil.ReadAll(rsp.Body)
		asserts.NoError(err)

		arts := make([]module.Article, 0)
		err = json.Unmarshal(body, &arts)
		asserts.NoError(err)

		asserts.Equal(test.Expected, arts)
	}
}

func TestGetPost(t *testing.T) {
	ctx := context.Background()
	ts := test.NewTestRouter(gin.Default(), ArticleService)
	d, r := test.NewTestDB()
	defer func() {
		d.Close()
		r.Close()
	}()
	tk := test.NewTestSessionToken(r)

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
		test.DeleteComments(ctx, d, comments...)
		test.DeletePosts(ctx, d, posts...)
		test.DeleteAccounts(ctx, d, users...)
	}()

	header := map[string]string{
		utils.HEADER_SESSION_TOKEN: tk,
	}

	tests := []struct {
		Description string
		Article     module.Article
		StatusCode  int
		Comments    int
		Expected    module.Article
	}{
		{
			Description: "Get post 1 with 2 comments",
			Article:     *posts[0],
			StatusCode:  http.StatusOK,
			Comments:    2,
			Expected: module.Article{
				ID:       posts[0].ID,
				Author:   posts[0].Author,
				Title:    posts[0].Title,
				Content:  posts[0].Content,
				CreateAt: posts[0].CreateAt,
				UpdateAt: posts[0].UpdateAt,
				Comments: []module.Comment{*comments[0], *comments[3]},
			},
		}, {
			Description: "Get post 2 with 3 comments",
			Article:     *posts[1],
			StatusCode:  http.StatusOK,
			Comments:    3,
			Expected: module.Article{
				ID:       posts[1].ID,
				Author:   posts[1].Author,
				Title:    posts[1].Title,
				Content:  posts[1].Content,
				CreateAt: posts[1].CreateAt,
				UpdateAt: posts[1].UpdateAt,
				Comments: []module.Comment{*comments[1], *comments[2], *comments[4]},
			},
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
		t.Run(test.Description, func(t *testing.T) {
			req, err := NewRequestWithTestHeader("GET", ts.URL+"/article/"+fmt.Sprint(test.Article.ID), nil, header)
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
			asserts.Equal(test.Expected, *article)
		})
	}
}

func TestEditPost(t *testing.T) {
	ctx := context.Background()
	ts := test.NewTestRouter(gin.Default(), ArticleService)
	d, r := test.NewTestDB()
	defer func() {
		d.Close()
		r.Close()
	}()
	tk := test.NewTestSessionToken(r)

	users := []test.LoggedAccount{
		test.NewAccountWithLogin(ctx, d, r),
		test.NewAccountWithLogin(ctx, d, r),
	}

	posts := []*module.Article{
		test.NewPost(ctx, d, users[0].Username),
		test.NewPost(ctx, d, users[0].Username),
		test.NewPost(ctx, d, users[1].Username),
	}

	header := map[string]string{
		utils.HEADER_SESSION_TOKEN: tk,
	}

	defer func() {
		test.DeletePosts(ctx, d, posts...)
		test.DeleteLoggedAccounts(ctx, d, users...)
	}()

	tests := []struct {
		Description    string
		Article        module.Article
		Token          string
		EditStatusCode int
		Expected       module.Article
		GetStatusCode  int
	}{
		{
			Description: "Edit existed article content with corresponding author",
			Article: module.Article{
				ID:      posts[0].ID,
				Author:  posts[0].Author,
				Title:   posts[0].Title,
				Content: "This article is updated",
			},
			Token:          users[0].Token,
			EditStatusCode: http.StatusOK,
			Expected: module.Article{
				ID:      posts[0].ID,
				Author:  posts[0].Author,
				Title:   posts[0].Title,
				Content: "This article is updated",
			},
			GetStatusCode: http.StatusOK,
		}, {
			Description: "Edit existed article title with corresponding author",
			Article: module.Article{
				ID:      posts[1].ID,
				Author:  posts[1].Author,
				Title:   posts[1].Title,
				Content: "This article is updated",
			},
			Token:          users[0].Token,
			EditStatusCode: http.StatusOK,
			Expected: module.Article{
				ID:      posts[1].ID,
				Author:  posts[1].Author,
				Title:   posts[1].Title,
				Content: "This article is updated",
			},
			GetStatusCode: http.StatusOK,
		},
		{
			Description: "Edit non-exist article",
			Article: module.Article{
				ID:      92673,
				Author:  users[1].Username,
				Title:   "This is new title",
				Content: posts[1].Content,
			},
			Token:          users[0].Token,
			EditStatusCode: http.StatusNotFound,
			GetStatusCode:  http.StatusNotFound,
		}, {
			Description: "Edit existed article title with wrong author",
			Article: module.Article{
				ID:      posts[1].ID,
				Author:  posts[1].Author,
				Title:   posts[1].Title,
				Content: "This article is updated",
			},
			Token: users[1].Token,
			//TODO: Need to be http.StatusForbidden, need fix EditPost
			EditStatusCode: http.StatusNotFound,
			Expected: module.Article{
				ID:      posts[1].ID,
				Author:  posts[1].Author,
				Title:   posts[1].Title,
				Content: "This article is updated",
			},
			GetStatusCode: http.StatusOK,
		},
	}

	asserts := assert.New(t)

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			b, err := json.Marshal(&test.Article)
			asserts.NoError(err)

			header[utils.HEADER_SESSION_TOKEN] = test.Token
			ereq, err := NewRequestWithTestHeader("PUT", ts.URL+"/article/"+fmt.Sprint(test.Article.ID), bytes.NewBuffer(b), header)
			asserts.NoError(err)

			ersp, err := http.DefaultClient.Do(ereq)
			asserts.NoError(err)
			defer ersp.Body.Close()

			ebody, err := ioutil.ReadAll(ersp.Body)
			asserts.NoError(err)

			asserts.Equal(test.EditStatusCode, ersp.StatusCode, string(ebody))

			// Get article back after updating
			greq, err := NewRequestWithTestHeader("GET", ts.URL+"/article/"+fmt.Sprint(test.Article.ID), nil, header)
			asserts.NoError(err)

			grsp, err := http.DefaultClient.Do(greq)
			asserts.NoError(err)
			defer grsp.Body.Close()

			asserts.Equal(test.GetStatusCode, grsp.StatusCode)

			if test.GetStatusCode == http.StatusNotFound {
				return
			}

			body, err := ioutil.ReadAll(grsp.Body)
			asserts.NoError(err)

			art := new(module.Article)
			err = json.Unmarshal(body, &art)
			asserts.NoError(err)

			asserts.Equal(test.Expected.Title, art.Title)
			asserts.Equal(test.Expected.Content, art.Content)
		})
	}
}

func TestDeletePost(t *testing.T) {
	ctx := context.Background()
	ts := test.NewTestRouter(gin.Default(), ArticleService)
	d, r := test.NewTestDB()
	defer func() {
		d.Close()
		r.Close()
	}()
	tk := test.NewTestSessionToken(r)

	users := []test.LoggedAccount{
		test.NewAccountWithLogin(ctx, d, r),
		test.NewAccountWithLogin(ctx, d, r),
		test.NewAccountWithLogin(ctx, d, r),
	}

	posts := []*module.Article{
		test.NewPost(ctx, d, users[0].Username),
		test.NewPost(ctx, d, users[1].Username),
		test.NewPost(ctx, d, users[1].Username),
		test.NewPost(ctx, d, users[2].Username),
	}

	comments := []*module.Comment{
		test.NewComment(ctx, d, users[0].Username, posts[1].ID),
		test.NewComment(ctx, d, users[0].Username, posts[1].ID),
		test.NewComment(ctx, d, users[0].Username, posts[2].ID),
		test.NewComment(ctx, d, users[1].Username, posts[2].ID),
		test.NewComment(ctx, d, users[1].Username, posts[2].ID),
		test.NewComment(ctx, d, users[0].Username, posts[2].ID),
	}

	defer func() {
		test.DeleteComments(ctx, d, comments...)
		test.DeletePosts(ctx, d, posts...)
		test.DeleteLoggedAccounts(ctx, d, users...)
	}()

	header := map[string]string{
		utils.HEADER_SESSION_TOKEN: tk,
	}

	tests := []struct {
		Description      string
		Article          module.Article
		Token            string
		DeleteStatusCode int
		GetStatusCode    int
	}{
		{
			Description:      "Delete post without comments with corresponding user",
			Article:          *posts[0],
			Token:            users[0].Token,
			DeleteStatusCode: http.StatusOK,
			GetStatusCode:    http.StatusNotFound,
		}, {
			Description:      "Delete post with comments with corresponding user",
			Article:          *posts[1],
			Token:            users[1].Token,
			DeleteStatusCode: http.StatusOK,
			GetStatusCode:    http.StatusNotFound,
		}, {
			Description:      "Delete post without comments with wrong user",
			Article:          *posts[3],
			Token:            users[0].Token,
			DeleteStatusCode: http.StatusNotFound,
			GetStatusCode:    http.StatusOK,
		}, {
			Description:      "Delete post with comments with wrong user",
			Article:          *posts[2],
			Token:            users[0].Token,
			DeleteStatusCode: http.StatusNotFound,
			GetStatusCode:    http.StatusOK,
		}, {
			Description: "Delete non-exist post",
			Article: module.Article{
				ID: 94518,
			},
			DeleteStatusCode: http.StatusNotFound,
			GetStatusCode:    http.StatusNotFound,
		},
	}

	asserts := assert.New(t)

	for _, test := range tests {
		header[utils.HEADER_SESSION_TOKEN] = test.Token
		dreq, err := NewRequestWithTestHeader("DELETE", ts.URL+"/article/"+fmt.Sprint(test.Article.ID), nil, header)
		asserts.NoError(err)

		drsp, err := http.DefaultClient.Do(dreq)
		asserts.NoError(err)
		defer drsp.Body.Close()

		dbody, err := ioutil.ReadAll(drsp.Body)
		asserts.NoError(err)

		asserts.Equal(test.DeleteStatusCode, drsp.StatusCode, string(dbody))

		greq, err := NewRequestWithTestHeader("GET", ts.URL+"/article/"+fmt.Sprint(test.Article.ID), nil, header)
		asserts.NoError(err)

		grsp, err := http.DefaultClient.Do(greq)
		asserts.NoError(err)
		defer grsp.Body.Close()

		asserts.Equal(test.GetStatusCode, grsp.StatusCode)
	}
}

func TestNewAndGetComments(t *testing.T) {
	ctx := context.Background()
	ts, d, r, hdr := test.NewTestEntities(gin.Default(), ArticleService)
	defer func() {
		d.Close()
		r.Close()
	}()

	users := []*module.Account{
		test.NewAccount(ctx, d, false),
		test.NewAccount(ctx, d, false),
	}

	posts := []*module.Article{
		test.NewPost(ctx, d, users[0].Username),
	}

	comments := []*module.Comment{
		{
			Author:  users[1].Username,
			Content: "This is 1st comment for post[0]",
		}, {
			Author:  users[1].Username,
			Content: "This is 2nd comment for post[0]",
		}, {
			Author:  users[1].Username,
			Content: "This is 1st comment for non-existed",
		}, {
			Author:  users[1].Username,
			Content: "This is 2nd comment for non-existed",
		},
	}

	defer func() {
		test.DeleteComments(ctx, d, comments...)
		test.DeletePosts(ctx, d, posts...)
		test.DeleteAccounts(ctx, d, users...)
	}()

	tests := []struct {
		Description        string
		Article            module.Article
		Comment            []*module.Comment
		PostStatusCode     int
		GetStatusCode      int
		ExpectedCommentLen int
		Expected           []*module.Comment
	}{
		{
			Description:        "New comments to exist post",
			Article:            *posts[0],
			Comment:            comments[:2],
			PostStatusCode:     http.StatusOK,
			GetStatusCode:      http.StatusOK,
			ExpectedCommentLen: 2,
			Expected:           comments[:2],
		}, {
			Description: "New comments to non-existed post",
			Article: module.Article{
				ID: 56482,
			},
			Comment:            comments[2:4],
			PostStatusCode:     http.StatusNotFound,
			GetStatusCode:      http.StatusOK,
			ExpectedCommentLen: 0,
			Expected:           nil,
		},
	}

	asserts := assert.New(t)

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			for _, comment := range test.Comment {
				b, err := json.Marshal(comment)
				asserts.NoError(err)

				nreq, err := NewRequestWithTestHeader("POST", ts.URL+"/article/"+fmt.Sprint(test.Article.ID)+"/comment", bytes.NewBuffer(b), hdr)
				asserts.NoError(err)

				nrsp, err := http.DefaultClient.Do(nreq)
				asserts.NoError(err)
				defer nrsp.Body.Close()

				nbody, err := ioutil.ReadAll(nrsp.Body)
				asserts.NoError(err)

				asserts.Equal(test.PostStatusCode, nrsp.StatusCode, string(nbody))
			}

			greq, err := NewRequestWithTestHeader("GET", ts.URL+"/article/"+fmt.Sprint(test.Article.ID)+"/comments", nil, hdr)
			asserts.NoError(err)

			grsp, err := http.DefaultClient.Do(greq)
			asserts.NoError(err)
			defer grsp.Body.Close()

			asserts.Equal(test.GetStatusCode, grsp.StatusCode)

			body, err := ioutil.ReadAll(grsp.Body)
			asserts.NoError(err)

			comments := []*module.Comment{}
			err = json.Unmarshal(body, &comments)
			asserts.NoError(err)

			asserts.Equal(test.ExpectedCommentLen, len(comments))

			for i := 0; i < len(test.Expected); i++ {
				asserts.Equal(test.Expected[i].Author, comments[i].Author)
				asserts.Equal(test.Expected[i].Content, comments[i].Content)
			}
		})
	}
}

func TestDeleteComment(t *testing.T) {
	ctx := context.Background()
	ts, d, r, hdr := test.NewTestEntities(gin.Default(), ArticleService)
	defer func() {
		d.Close()
		r.Close()
	}()

	users := []test.LoggedAccount{
		test.NewAccountWithLogin(ctx, d, r),
		test.NewAccountWithLogin(ctx, d, r),
		test.NewAccountWithLogin(ctx, d, r),
	}

	posts := []*module.Article{
		test.NewPost(ctx, d, users[2].Username),
	}

	comments := []*module.Comment{
		test.NewComment(ctx, d, users[0].Username, posts[0].ID),
		test.NewComment(ctx, d, users[2].Username, posts[0].ID),
		test.NewComment(ctx, d, users[1].Username, posts[0].ID),
		test.NewComment(ctx, d, users[2].Username, posts[0].ID),
	}

	defer func() {
		test.DeleteComments(ctx, d, comments...)
		test.DeletePosts(ctx, d, posts...)
		test.DeleteLoggedAccounts(ctx, d, users...)
	}()

	tests := []struct {
		Description        string
		Token              string
		Article            module.Article
		Comment            module.Comment
		DeleteStatusCode   int
		GetStatusCode      int
		ExpectedCommentLen int
		Expected           []module.Comment
	}{
		{
			Description:        "Delete comment from exist post and exist comment with comment author",
			Token:              users[0].Token,
			Article:            *posts[0],
			Comment:            *comments[0],
			DeleteStatusCode:   http.StatusOK,
			GetStatusCode:      http.StatusOK,
			ExpectedCommentLen: 3,
			Expected:           []module.Comment{*comments[1], *comments[2], *comments[3]},
		}, {
			Description: "Delete comment from exist post and non-exist comment",
			Token:       users[0].Token,
			Article:     *posts[0],
			Comment: module.Comment{
				ID: 292823,
			},
			DeleteStatusCode:   http.StatusNotFound,
			GetStatusCode:      http.StatusOK,
			ExpectedCommentLen: 3,
			Expected:           []module.Comment{*comments[1], *comments[2], *comments[3]},
		}, {
			Description: "Delete comment from non-exist post",
			Token:       users[0].Token,
			Article: module.Article{
				ID: 23914,
			},
			Comment: module.Comment{
				ID: 292823,
			},
			DeleteStatusCode:   http.StatusNotFound,
			GetStatusCode:      http.StatusOK,
			ExpectedCommentLen: 0,
			Expected:           []module.Comment{},
		}, {
			Description:        "Delete comment from exist post and exist comment with post author",
			Token:              users[2].Token,
			Article:            *posts[0],
			Comment:            *comments[1],
			DeleteStatusCode:   http.StatusOK,
			GetStatusCode:      http.StatusOK,
			ExpectedCommentLen: 2,
			Expected:           []module.Comment{*comments[2], *comments[3]},
		},
	}

	asserts := assert.New(t)

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			hdr[utils.HEADER_SESSION_TOKEN] = test.Token
			dreq, err := NewRequestWithTestHeader("DELETE",
				ts.URL+"/article/"+fmt.Sprint(test.Article.ID)+"/comment/"+fmt.Sprint(test.Comment.ID), nil, hdr)
			asserts.NoError(err)

			drsp, err := http.DefaultClient.Do(dreq)
			asserts.NoError(err)
			defer drsp.Body.Close()

			dbody, err := ioutil.ReadAll(drsp.Body)
			asserts.NoError(err)

			asserts.Equal(test.DeleteStatusCode, drsp.StatusCode, string(dbody))

			greq, err := NewRequestWithTestHeader("GET", ts.URL+"/article/"+fmt.Sprint(test.Article.ID)+"/comments", nil, hdr)
			asserts.NoError(err)

			grsp, err := http.DefaultClient.Do(greq)
			asserts.NoError(err)
			defer grsp.Body.Close()

			asserts.Equal(test.GetStatusCode, grsp.StatusCode)

			body, err := ioutil.ReadAll(grsp.Body)
			asserts.NoError(err)

			comments := []*module.Comment{}
			err = json.Unmarshal(body, &comments)
			asserts.NoError(err)

			asserts.Equal(test.ExpectedCommentLen, len(comments))

			for i := 0; i < len(test.Expected); i++ {
				asserts.Equal(test.Expected[i].Author, comments[i].Author)
				asserts.Equal(test.Expected[i].Content, comments[i].Content)
			}
		})
	}
}

func TestUpdateComment(t *testing.T) {
	ctx := context.Background()
	ts, d, r, hdr := test.NewTestEntities(gin.Default(), ArticleService)
	defer func() {
		d.Close()
		r.Close()
	}()

	users := []test.LoggedAccount{
		test.NewAccountWithLogin(ctx, d, r),
		test.NewAccountWithLogin(ctx, d, r),
		test.NewAccountWithLogin(ctx, d, r),
	}

	posts := []*module.Article{
		test.NewPost(ctx, d, users[2].Username),
	}

	comments := []*module.Comment{
		test.NewComment(ctx, d, users[0].Username, posts[0].ID),
		test.NewComment(ctx, d, users[2].Username, posts[0].ID),
		test.NewComment(ctx, d, users[1].Username, posts[0].ID),
		test.NewComment(ctx, d, users[2].Username, posts[0].ID),
	}

	defer func() {
		test.DeleteComments(ctx, d, comments...)
		test.DeletePosts(ctx, d, posts...)
		test.DeleteLoggedAccounts(ctx, d, users...)
	}()

	newComment := func(id, art_id int, author, content string) *module.Comment {
		return &module.Comment{
			ID:        id,
			ArticleID: art_id,
			Author:    author,
			Content:   content,
		}
	}

	tests := []struct {
		Description        string
		Token              string
		Article            module.Article
		Comment            module.Comment
		UpdateStatusCode   int
		GetStatusCode      int
		ExpectedCommentLen int
		Expected           []module.Comment
	}{
		{
			Description:        "Update comment from exist post and exist comment with comment author",
			Token:              users[0].Token,
			Article:            *posts[0],
			Comment:            *newComment(comments[0].ID, posts[0].ID, comments[0].Author, "update content"),
			UpdateStatusCode:   http.StatusOK,
			GetStatusCode:      http.StatusOK,
			ExpectedCommentLen: 4,
			Expected: []module.Comment{*newComment(comments[0].ID, posts[0].ID, comments[0].Author, "update content"),
				*comments[1], *comments[2], *comments[3]},
		}, {
			Description:        "Update comment from exist post and exist comment with other author",
			Token:              users[2].Token,
			Article:            *posts[0],
			Comment:            *newComment(comments[0].ID, posts[0].ID, comments[0].Author, "never happen"),
			UpdateStatusCode:   http.StatusNotFound, //TODO: should be http.StatusForbidden, need fix UpdateComment
			GetStatusCode:      http.StatusOK,
			ExpectedCommentLen: 4,
			Expected: []module.Comment{*newComment(comments[0].ID, posts[0].ID, comments[0].Author, "update content"),
				*comments[1], *comments[2], *comments[3]},
		}, {
			Description: "Update comment from exist post and non-exist comment",
			Article:     *posts[0],
			Comment: module.Comment{
				ID: 292823,
			},
			UpdateStatusCode:   http.StatusNotFound,
			GetStatusCode:      http.StatusOK,
			ExpectedCommentLen: 4,
			Expected: []module.Comment{*newComment(comments[0].ID, posts[0].ID, comments[0].Author, "update content"),
				*comments[1], *comments[2], *comments[3]},
		}, {
			Description: "Update comment from non-exist post",
			Article: module.Article{
				ID: 23914,
			},
			Comment: module.Comment{
				ID: 292823,
			},
			UpdateStatusCode:   http.StatusNotFound,
			GetStatusCode:      http.StatusOK,
			ExpectedCommentLen: 0,
			Expected:           []module.Comment{},
		}, {
			Description:        "Update comment from exist post and exist comment with post author, it should fail",
			Token:              users[2].Token,
			Article:            *posts[0],
			Comment:            *newComment(comments[0].ID, posts[0].ID, comments[0].Author, "never happen"),
			UpdateStatusCode:   http.StatusNotFound,
			GetStatusCode:      http.StatusOK,
			ExpectedCommentLen: 4,
			Expected: []module.Comment{*newComment(comments[0].ID, posts[0].ID, comments[0].Author, "update content"),
				*comments[1], *comments[2], *comments[3]},
		},
	}

	asserts := assert.New(t)

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			hdr[utils.HEADER_SESSION_TOKEN] = test.Token
			b, err := json.Marshal(&test.Comment)
			asserts.NoError(err)

			dreq, err := NewRequestWithTestHeader("PUT",
				ts.URL+"/article/"+fmt.Sprint(test.Article.ID)+"/comment/"+fmt.Sprint(test.Comment.ID), bytes.NewBuffer(b), hdr)
			asserts.NoError(err)

			drsp, err := http.DefaultClient.Do(dreq)
			asserts.NoError(err)
			defer drsp.Body.Close()

			dbody, err := ioutil.ReadAll(drsp.Body)
			asserts.NoError(err)

			asserts.Equal(test.UpdateStatusCode, drsp.StatusCode, string(dbody))

			greq, err := NewRequestWithTestHeader("GET", ts.URL+"/article/"+fmt.Sprint(test.Article.ID)+"/comments", nil, hdr)
			asserts.NoError(err)

			grsp, err := http.DefaultClient.Do(greq)
			asserts.NoError(err)
			defer grsp.Body.Close()

			asserts.Equal(test.GetStatusCode, grsp.StatusCode)

			body, err := ioutil.ReadAll(grsp.Body)
			asserts.NoError(err)

			comments := []*module.Comment{}
			err = json.Unmarshal(body, &comments)
			asserts.NoError(err)

			asserts.Equal(test.ExpectedCommentLen, len(comments))

			for i := 0; i < len(test.Expected); i++ {
				asserts.Equal(test.Expected[i].Author, comments[i].Author)
				asserts.Equal(test.Expected[i].Content, comments[i].Content)
			}
		})
	}
}
