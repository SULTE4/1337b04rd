package postgres

import (
	"1337b04rd/internal/domain/post"
	"database/sql"
	"time"
)

type PostgresPostRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) *PostgresPostRepo {
	return &PostgresPostRepo{db: db}
}

func (r *PostgresPostRepo) Insert(p post.Post) (int, error) {
	stmt := `INSERT INTO "Post" (title, content, imageURL, userID, created, expires)
			VALUES ($1, $2, $3, $4, $5, $6) returning id;`
	id := 0

	err := r.db.QueryRow(stmt, p.Title, p.Content, p.ImageURL, p.UserID, p.Created, p.Expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostgresPostRepo) GetByID(id int) (post.Post, error) {
	return post.Post{}, nil
}

func (r *PostgresPostRepo) GetAll() ([]post.Post, error) {
	// stmt := `SELECT `
	return []post.Post{}, nil
}

func (r *PostgresPostRepo) DeleteById(id int) error {
	stmt := `DELETE FROM "Post"
			WHERE id = $1`

	_, err := r.db.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresPostRepo) UpdatePostExpire(id int) error {
	stmt := `UPDATE "Post"
			SET expires = $1
			WHERE id = $2`

	t := time.Now().UTC().Add(15 * time.Minute)

	_, err := r.db.Exec(stmt, t, id)
	if err != nil {
		return err
	}
	return nil
}
