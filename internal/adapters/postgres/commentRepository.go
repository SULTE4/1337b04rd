package postgres

import (
	"1337b04rd/internal/core/domain"
	"database/sql"
	"log/slog"
	"time"
)

type PostgresCommentRepo struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewCommentRepo(db *sql.DB, logger *slog.Logger) *PostgresCommentRepo {
	return &PostgresCommentRepo{db: db, logger: logger}
}

func (r *PostgresCommentRepo) Insert(com domain.Comment) error {
	stmt := `INSERT INTO comment (userid, content, imageurl, created, postid)
			VALUES($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(stmt, com.UserID, com.Content, com.ImageURL, com.Created, com.PostID)
	if err != nil {
		r.logger.Error(err.Error())
		return err
	}

	return nil
}

func (r *PostgresCommentRepo) GetCommentsByPost(id int) ([]domain.Comment, error) {
	stmt := `SELECT * FROM comment 
			WHERE postid = $1`

	rows, err := r.db.Query(stmt, id)
	if err != nil {
		r.logger.Error(err.Error())
		return []domain.Comment{}, err
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var c domain.Comment
		err = rows.Scan(
			&c.CommentID,
			&c.Content,
			&c.ImageURL,
			&c.Created,
			&c.UserID,
			&c.PostID,
		)
		if err != nil {
			r.logger.Error(err.Error())
			return []domain.Comment{}, err
		}
		comments = append(comments, c)
	}
	err = rows.Err()
	if err != nil {
		r.logger.Error(err.Error())
		return []domain.Comment{}, err
	}
	return comments, nil
}

func (r *PostgresCommentRepo) UpdatePostExpire(id int) error {
	stmt := `UPDATE post
			SET expires = $1
			WHERE id = $2`

	t := time.Now().UTC().Add(15 * time.Minute)

	_, err := r.db.Exec(stmt, t, id)
	if err != nil {
		r.logger.Error(err.Error())
		return err
	}
	return nil
}

// if post expired then true, else false
func (r *PostgresCommentRepo) IsPostExpired(id int) (bool, error) {
	stmt := `SELECT (expires < NOW()) FROM post
		WHERE  id = $1`

	var is bool

	_ = r.db.QueryRow(stmt, id).Scan(&is)

	return is, nil
}
