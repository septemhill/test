package module

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/septemhill/test/db"
)

type Article struct {
	ID       int       `db:"id" json:"id"`
	Author   string    `db:"author" json:"author"`
	Title    string    `db:"title" json:"title"`
	Content  string    `db:"content" json:"content"`
	CreateAt time.Time `db:"create_at" json:"createAt"`
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	Comments []Comment `json:"comments"`
}

type Comment struct {
	ID        int       `db:"id" json:"id"`
	ArticleID int       `db:"art_id" json:"art_id"`
	Author    string    `db:"author" json:"author"`
	Content   string    `db:"content" json:"content"`
	CreateAt  time.Time `db:"create_at" json:"createAt"`
	UpdateAt  time.Time `db:"update_at" json:"updateAt"`
}

func NewPost(ctx context.Context, db *db.DB, art Article) error {
	expr := `INSERT INTO articles VALUES (DEFAULT, $1, $2, $3, $4, $5)`
	return txAction(ctx, db, func(tx *sqlx.Tx) error {
		if _, err := tx.ExecContext(ctx, expr, art.Author, art.Title, art.Content, time.Now(), time.Now()); err != nil {
			return err
		}
		return nil
	})
}

func EditPost(ctx context.Context, db *db.DB, art Article) error {
	expr := `UPDATE articles SET title = $1, content = $2, update_at = $3 WHERE id = $4`
	return txAction(ctx, db, func(tx *sqlx.Tx) error {
		if _, err := tx.ExecContext(ctx, expr, art.Title, art.Content, time.Now(), art.ID); err != nil {
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
	expr := `SELECT * FROM articles ORDER BY create_at %s LIMIT $1 OFFSET $2`
	if !asc {
		expr = fmt.Sprintf(expr, "DESC")
	} else {
		expr = fmt.Sprintf(expr, "ASC")
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryxContext(ctx, expr, size, offset)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	arts := make([]Article, 0)

	for rows.Next() {
		art := Article{}
		if err := rows.StructScan(&art); err != nil {
			return nil, err
		}

		arts = append(arts, art)
	}

	return arts, nil
}

func GetPost(ctx context.Context, db *db.DB, art Article) (*Article, error) {
	return nil, nil
}

func NewComment(ctx context.Context, db *db.DB, postID string, comment Comment) error {
	expr := `INSERT INTO comments VALUES (DEFAULT, $1, $2, $3, $4, $5)`

	return txAction(ctx, db, func(tx *sqlx.Tx) error {
		if _, err := tx.ExecContext(ctx, expr, postID, comment.Author, comment.Content, time.Now(), time.Now()); err != nil {
			return err
		}
		return nil
	})
}

func UpdateComment(ctx context.Context, db *db.DB, comment Comment) error {
	expr := `UPDATE comments SET content = $1 WHERE id = $2`

	return txAction(ctx, db, func(tx *sqlx.Tx) error {
		if _, err := tx.ExecContext(ctx, expr, comment.ID); err != nil {
			return err
		}
		return nil
	})
}

func GetComments(ctx context.Context, db *db.DB, postID string, size, offset int) ([]Comment, error) {
	expr := `SELECT * FROM comments WHERE art_id = $1 ORDER BY create_at ASC LIMIT $2 OFFSET $3`

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	comments := make([]Comment, 0)

	rows, err := tx.QueryxContext(ctx, expr, postID, size, offset)
	if err != nil {
		return nil, err
	}
	tx.Commit()

	for rows.Next() {
		comment := Comment{}
		if err := rows.StructScan(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func DeleteComment(ctx context.Context, db *db.DB, comment Comment) error {
	expr := `DELETE FROM comments WHERE id = $1`

	return txAction(ctx, db, func(tx *sqlx.Tx) error {
		if _, err := tx.ExecContext(ctx, expr, comment.ID); err != nil {
			return err
		}
		return nil
	})
}
