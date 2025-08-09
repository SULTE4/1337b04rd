package postgres

import (
	"1337b04rd/internal/domain/post"
	"database/sql"
)

type PostgresPostRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) *PostgresPostRepo {
	return &PostgresPostRepo{db: db}
}

func (r *PostgresPostRepo) Insert(p post.Post) error {
	stmt := `INSERT INTO Form (title, content, userID, created)
			VALUES ($1, $2, $3, $4)`

	_, err := r.db.Exec(stmt, p.Title, p.Content, p.UserID, p.Created)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresPostRepo) GetByID(id int) (post.Post, error) {
	return post.Post{}, nil
}

func (r *PostgresPostRepo) GetAll() ([]post.Post, error) {
	return []post.Post{}, nil
}

func (r *PostgresPostRepo) DeleteById(id int) (int, error) {
	return 0, nil
}
