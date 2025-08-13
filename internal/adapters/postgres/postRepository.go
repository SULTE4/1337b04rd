package postgres

import (
	"1337b04rd/internal/core/domain"
	"database/sql"
	"errors"
	"fmt"
)

type PostgresPostRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) *PostgresPostRepo {
	return &PostgresPostRepo{db: db}
}

func (r *PostgresPostRepo) Insert(p domain.Post) (int, error) {
	create := `create TABLE if not exists Post  (
    id serial primary key not null,
    title varchar(50) not null,
    content text,
    imageURL varchar(255),
    userID int not null,
    created timestamp not null,
    expires timestamp not null
);
`
	stmt := `INSERT INTO post (title, content, imageURL, userID, created, expires)
			VALUES ($1, $2, $3, $4, $5, $6) returning id;`
	id := 0

	_, err := r.db.Exec(create)
	if err != nil {
		return 0, err
	}

	err = r.db.QueryRow(stmt, p.Title, p.Content, p.ImageURL, p.UserID, p.Created, p.Expires).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostgresPostRepo) GetByID(id int) (domain.Post, error) {
	var p domain.Post
	stmt := `SELECT * FROM post
			WHERE id = $1`
	row := r.db.QueryRow(stmt, id)
	err := row.Scan(
		&p.ID,
		&p.Title,
		&p.Content,
		&p.ImageURL,
		&p.UserID,
		&p.Created,
		&p.Expires,
	)

	if err != nil {
		if errors.As(err, sql.ErrNoRows) {
			return domain.Post{}, fmt.Errorf("invalid post id: %d", id)
		}
		return domain.Post{}, err
	}
	return p, nil
}

func (r *PostgresPostRepo) GetAll() ([]domain.Post, error) {
	var posts []domain.Post

	stmt := `SELECT * 
			FROM post
			WHERE expires > NOW()
			ORDER BY expires ASC;`

	rows, err := r.db.Query(stmt)
	if err != nil {
		return []domain.Post{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var p domain.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.ImageURL,
			&p.UserID,
			&p.Created,
			&p.Expires,
		)
		if err != nil {
			return []domain.Post{}, err
		}
		posts = append(posts, p)
	}
	err = rows.Err()
	if err != nil {
		return []domain.Post{}, err
	}

	return posts, nil
}

func (r *PostgresPostRepo) DeleteById(id int) error {
	stmt := `DELETE FROM post
			WHERE id = $1`

	_, err := r.db.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresPostRepo) GetExpiredPosts() ([]domain.Post, error) {
	var posts []domain.Post

	stmt := `SELECT * 
			FROM post
			WHERE expires <= NOW()
			ORDER BY expires ASC;`

	rows, err := r.db.Query(stmt)
	if err != nil {
		return []domain.Post{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var p domain.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.ImageURL,
			&p.UserID,
			&p.Created,
			&p.Expires,
		)
		if err != nil {
			return []domain.Post{}, err
		}
		posts = append(posts, p)
	}
	err = rows.Err()
	if err != nil {
		return []domain.Post{}, err
	}

	return posts, nil
}

func (r *PostgresPostRepo) IsPostExpired(id int) error {
	return nil
}
