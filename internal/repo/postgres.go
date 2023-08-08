package repo

import (
	"context"
	"database/sql"
	"dmmak/simple-rest-crud/internal/entity"
	"dmmak/simple-rest-crud/internal/handler"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5"
)

type pgRepo struct {
	db *sql.DB
}

func NewPGRepo(db *sql.DB) handler.RedditPostsRepo {
	return &pgRepo{db}
}

func (pg *pgRepo) Save(ctx context.Context, post *entity.RedditPost) (err error) {
	tx, err := pg.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault, ReadOnly: false})
	if err != nil {
		return fmt.Errorf("can't create tx: %w", err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO posts VALUES ($1,$2,$3)", post.UUID, post.Title, post.Likes)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("can't rollback tx: %w, insert post err: %w", rbErr, err)
		}
		return fmt.Errorf("can't insert new post: %w", err)
	}

	if len(post.Comments) == 0 {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("can't commit inserting post without comments: %w", err)
		}
		return nil
	}

	sqlBuilder := &strings.Builder{}
	values := []any{}
	sqlBuilder.WriteString("INSERT INTO comments VALUES")
	for i, v := range post.Comments {
		sql := fmt.Sprintf("($%v, $%v, $%v, $%v),", 4*i+1, 4*i+2, 4*i+3, 4*i+4)
		sqlBuilder.WriteString(sql)
		values = append(values, v.UUID, post.UUID, v.Body, v.Likes)
	}
	insertSql := sqlBuilder.String()
	insertSql = insertSql[0 : len(insertSql)-1]

	_, err = tx.ExecContext(ctx, insertSql, values...)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("can't rollback tx: %w, insert comments err: %w", rbErr, err)
		}
		return fmt.Errorf("can't insert post's comments: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("can't commit inserting post with comments: %w", err)
	}
	return nil
}

func (pg *pgRepo) Get(ctx context.Context, postID string) (post *entity.RedditPost, found bool, err error) {
	post = &entity.RedditPost{}

	row := pg.db.QueryRowContext(ctx, "SELECT * FROM posts WHERE uuid = $1", postID)
	err = row.Scan(&post.UUID, &post.Title, &post.Likes)
	if err == sql.ErrNoRows {
		return post, false, nil
	} else if err != nil {
		return post, false, fmt.Errorf("can't query 'posts' table: %w", err)
	}

	comments := []*entity.Comment{}
	rows, err := pg.db.QueryContext(ctx, "SELECT uuid, body, likes FROM comments WHERE post_uuid = $1", postID)
	if err != nil {
		return nil, false, fmt.Errorf("can't query 'comments' table: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		comment := &entity.Comment{}
		err = rows.Scan(&comment.UUID, &comment.Body, &comment.Likes)
		if err != nil {
			return nil, false, fmt.Errorf("can't process query result: %w", err)
		}
		comments = append(comments, comment)
	}

	err = rows.Err()
	if err != nil {
		return nil, false, fmt.Errorf("error during query result iteration: %w", err)
	}

	post.Comments = comments
	return post, true, nil
}

func (pg *pgRepo) Delete(ctx context.Context, postID string) (found bool, err error) {
	res, err := pg.db.ExecContext(ctx, "DELETE FROM posts WHERE uuid = $1", postID)
	if err != nil {
		return false, fmt.Errorf("can't delete post: %w", err)
	}
	num, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("can't determine affected rows: %w", err)
	}
	if num == 0 {
		return false, nil
	} else {
		return true, nil
	}
}
