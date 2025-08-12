package postgres

import (
	"1337b04rd/internal/domain/post"
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

func (r *PostgresPostRepo) Insert(p post.Post) (int, error) {
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

func (r *PostgresPostRepo) GetByID(id int) (post.Post, error) {
	var p post.Post
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
			return post.Post{}, fmt.Errorf("invalid post id: %d", id)
		}
		return post.Post{}, err
	}
	return p, nil
}

func (r *PostgresPostRepo) GetAll() ([]post.Post, error) {
	var posts []post.Post

	stmt := `SELECT * 
			FROM post
			WHERE expires > NOW()
			ORDER BY expires ASC;`

	rows, err := r.db.Query(stmt)
	if err != nil {
		return []post.Post{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var p post.Post
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
			return []post.Post{}, err
		}
		posts = append(posts, p)
	}
	err = rows.Err()
	if err != nil {
		return []post.Post{}, err
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

func (r *PostgresPostRepo) GetExpiredPosts() ([]post.Post, error) {
	var posts []post.Post

	stmt := `SELECT * 
			FROM post
			WHERE expires <= NOW()
			ORDER BY expires ASC;`

	rows, err := r.db.Query(stmt)
	if err != nil {
		return []post.Post{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var p post.Post
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
			return []post.Post{}, err
		}
		posts = append(posts, p)
	}
	err = rows.Err()
	if err != nil {
		return []post.Post{}, err
	}

	return posts, nil
}

func (r *PostgresPostRepo) IsPostExpired(id int) error {
	return nil
}
