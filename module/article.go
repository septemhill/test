package module

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/septemhill/test/db"
)

type Article struct {
	ID       int       `db:"id" json:"id" uri:"id"`
	Author   string    `db:"author" json:"author"`
	Title    string    `db:"title" json:"title"`
	Content  string    `db:"content" json:"content"`
	CreateAt time.Time `db:"create_at" json:"createAt"`
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	Comments []Comment `json:"comments"`
}

type Comment struct {
	ID        int       `db:"id" json:"id" uri:"commentid"`
	ArticleID int       `db:"art_id" json:"art_id" uri:"id"`
	Author    string    `db:"author" json:"author"`
	Content   string    `db:"content" json:"content"`
	CreateAt  time.Time `db:"create_at" json:"createAt"`
	UpdateAt  time.Time `db:"update_at" json:"updateAt"`
}

func NewPost(ctx context.Context, d *db.DB, art *Article) (int, error) {
	var id int
	expr := `INSERT INTO articles VALUES (DEFAULT, $1, $2, $3, $4, $5) RETURNING id`
	err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, &id, expr, art.Author, art.Title, art.Content, time.Now(), time.Now()); err != nil {
			return err
		}
		return nil
	})

	return id, err
}

func EditPost(ctx context.Context, d *db.DB, art *Article) (int, error) {
	var id int
	expr := `UPDATE articles SET title = $1, content = $2, update_at = $3 WHERE id = $4 RETURNING id`
	err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, &id, expr, art.Title, art.Content, time.Now(), art.ID); err != nil {
			return err
		}
		return nil
	})

	return id, err
}

func DeletePost(ctx context.Context, d *db.DB, postID int) (int, error) {
	var id int
	expr := `DELETE FROM articles WHERE id = $1 RETURNING id`
	err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, &id, expr, postID); err != nil {
			return err
		}
		return nil
	})

	return id, err
}

func GetPosts(ctx context.Context, d *db.DB, size, offset int, asc bool) ([]Article, error) {
	expr := `SELECT * FROM articles ORDER BY create_at %s LIMIT $1 OFFSET $2`
	if !asc {
		expr = fmt.Sprintf(expr, "DESC")
	} else {
		expr = fmt.Sprintf(expr, "ASC")
	}

	arts := []Article{}
	err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		if err := tx.SelectContext(ctx, &arts, expr, size, offset); err != nil {
			return err
		}
		return nil
	})

	return arts, err
}

func GetPost(ctx context.Context, d *db.DB, postID int) (*Article, error) {
	expr := `SELECT * FROM articles WHERE id = $1`

	art := new(Article)

	if err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, art, expr, postID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	comments, err := GetComments(ctx, d, postID, 10, 0)
	if err != nil {
		return nil, err
	}

	art.Comments = comments

	return art, nil
}

func NewComment(ctx context.Context, d *db.DB, comment *Comment) (int, error) {
	var id int
	expr := `INSERT INTO comments VALUES (DEFAULT, $1, $2, $3, $4, $5) RETURNING id`

	err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, &id, expr, comment.ArticleID, comment.Author, comment.Content, time.Now(), time.Now()); err != nil {
			return err
		}
		return nil
	})

	return id, err
}

func UpdateComment(ctx context.Context, d *db.DB, comment *Comment) (int, error) {
	var id int
	expr := `UPDATE comments SET content = $1 WHERE id = $2`

	err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, &id, expr, comment.ID); err != nil {
			return err
		}
		return nil
	})

	return id, err
}

func GetComments(ctx context.Context, d *db.DB, postID, size, offset int) ([]Comment, error) {
	expr := `SELECT * FROM comments WHERE art_id = $1 ORDER BY create_at ASC LIMIT $2 OFFSET $3`

	comments := []Comment{}

	err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		if err := tx.SelectContext(ctx, &comments, expr, postID, size, offset); err != nil {
			return err
		}
		return nil
	})

	return comments, err
}

func DeleteComment(ctx context.Context, d *db.DB, comment *Comment) (int, error) {
	var id int
	expr := `DELETE FROM comments WHERE id = $1 RETURNING id`

	err := txAction(ctx, d, func(tx *sqlx.Tx) error {
		if err := tx.GetContext(ctx, &id, expr, comment.ID); err != nil {
			return err
		}
		return nil
	})

	return id, err
}
