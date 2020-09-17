package testing

import (
	"context"
	"strings"

	"github.com/septemhill/test/db"
	"github.com/septemhill/test/module"
	"github.com/sethvargo/go-diceware/diceware"
)

func NewPost(ctx context.Context, d *db.DB, author string) *module.Article {
	ts, _ := diceware.Generate(6)
	cs, _ := diceware.Generate(20)

	art := &module.Article{
		Author:  author,
		Title:   strings.Join(ts, " "),
		Content: strings.Join(cs, " "),
	}

	art.ID, _ = module.NewPost(ctx, d, art)
	arti, _ := module.GetPost(ctx, d, art.ID)

	return arti
}

func DeletePosts(ctx context.Context, d *db.DB, arts ...*module.Article) {
	for _, art := range arts {
		_, _ = module.DeletePost(ctx, d, art.ID)
	}
}

func DeletePostsByAccount(ctx context.Context, d *db.DB, accs ...string) {
	for _, acc := range accs {
		expr := `DELETE FROM articles WHERE author = $1`
		tx, _ := d.BeginTxx(ctx, nil)
		if err := tx.GetContext(ctx, expr, acc); err != nil {
			_ = tx.Rollback()
		}
		_ = tx.Commit()
	}
}

func NewComment(ctx context.Context, d *db.DB, author string, artID int) *module.Comment {
	cs, _ := diceware.Generate(10)

	comment := &module.Comment{
		ArticleID: artID,
		Author:    author,
		Content:   strings.Join(cs, " "),
	}

	comment.ID, _ = module.NewComment(ctx, d, comment)
	comm, _ := module.GetComment(ctx, d, artID, comment.ID)

	return comm
}

func DeleteComments(ctx context.Context, d *db.DB, comments ...*module.Comment) {
	for _, comment := range comments {
		_, _ = module.DeleteComment(ctx, d, comment)
	}
}
