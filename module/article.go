package module

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/septemhill/test/db"
)

type Article struct {
	ID       int      `db:"id" json:"id"`
	Author   string   `db:"author" json:"author"`
	Title    string   `db:"title" json:"title"`
	Content  string   `db:"content" json:"content"`
	Comments []string `json:"comments"`
}

type Comment struct {
	ID      int    `db:"id" json:"id"`
	Author  string `db:"author" json:"author"`
	Content string `db:"content" json:"content"`
}

func NewPost(ctx context.Context, db *db.DB, art Article) error {
	expr := `INSERT INTO articles VALUES (DEFAULT, $1, $2, $3)`
	return txAction(ctx, db, func(tx *sqlx.Tx) error {
		if _, err := tx.ExecContext(ctx, expr, art.Author, art.Title, art.Content); err != nil {
			return err
		}
		return nil
	})
}

func EditPost(ctx context.Context, db *db.DB, art Article) error {
	expr := `UPDATE articles SET title = $1 AND content = $2 WHERE id = $3`
	return txAction(ctx, db, func(tx *sqlx.Tx) error {
		if _, err := tx.ExecContext(ctx, expr, art.Title, art.Content, art.ID); err != nil {
			return err
		}
		return nil
	})
}

func DeletePost(ctx context.Context, db *db.DB, art Article) error {
	expr := `DELETE FROM articles WHERE id = $1`
	return txAction(ctx, db, func(tx *sqlx.Tx) error {
		if _, err := tx.ExecContext(ctx, expr, art.ID); err != nil {
			return err
		}
		return nil
	})
}

func GetPosts(ctx context.Context, db *db.DB, size, offset int, asc bool) ([]Article, error) {
	return nil, nil
}

func GetPost(ctx context.Context, db *db.DB, art Article) (*Article, error) {
	return nil, nil
}

func NewComment(ctx context.Context, db *db.DB, comment Comment) error {
	return nil
}

func UpdateComment(ctx context.Context, db *db.DB, comment Comment) error {
	return nil
}

func GetComments(ctx context.Context, db *db.DB, size, offset int, asc bool) ([]Comment, error) {
	return nil, nil
}

func DeleteComment(ctx context.Context, db *db.DB, comment Comment) error {
	//expr := `DELETE FROM comments WHERE id = comment`
	return nil
}
