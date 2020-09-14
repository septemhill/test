package testing

import (
	"context"
	"strings"
	"time"

	"github.com/septemhill/test/db"
	"github.com/septemhill/test/module"
	"github.com/sethvargo/go-diceware/diceware"
)

func NewPost(ctx context.Context, d *db.DB, acc *module.Account, comments []module.Comment) *module.Article {
	ts, _ := diceware.Generate(6)
	cs, _ := diceware.Generate(20)

	art := &module.Article{
		Author:   acc.Username,
		Title:    strings.Join(ts, " "),
		Content:  strings.Join(cs, " "),
		Comments: comments,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
	_ = module.NewPost(ctx, d, art)

	return art
}

func DeletePosts(ctx context.Context, d *db.DB, arts ...*module.Article) {
	for _, art := range arts {
		_ = module.DeletePost(ctx, d, art)
	}
}

func NewComment(ctx context.Context, d *db.DB, acc *module.Account, art *module.Article) *module.Comment {
	cs, _ := diceware.Generate(10)

	comment := &module.Comment{
		ArticleID: art.ID,
		Author:    acc.Username,
		Content:   strings.Join(cs, " "),
		CreateAt:  time.Now(),
		UpdateAt:  time.Now(),
	}

	_ = module.NewComment(ctx, d, art.ID, comment)

	return comment
}
