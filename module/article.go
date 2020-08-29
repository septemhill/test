package module

import (
	"context"

	"github.com/septemhill/test/db"
)

type Article struct {
	Title    string   `db:"title" json:"title"`
	Content  string   `db:"content" json:"content"`
	Comments []string `json:"comments"`
}

func NewPost(ctx context.Context, db *db.DB, art Article) error {
	return nil
}

func EditPost(ctx context.Context, db *db.DB, art Article) error {
	return nil
}

func DeletePost(ctx context.Context, db *db.DB, art Article) error {
	return nil
}

func GetPosts(ctx context.Context, db *db.DB, size, offset int, asc bool) ([]Article, error) {
	return nil, nil
}
