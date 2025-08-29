package postgres

import (
	"1337b04rd/internal/core/domain"
	"database/sql"
	"log/slog"
)

type PostgresPostRepo struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewPostRepo(db *sql.DB, logger *slog.Logger) *PostgresPostRepo {
	return &PostgresPostRepo{db: db, logger: logger}
}

func (r *PostgresPostRepo) Insert(p domain.Post) (int, error) {
	stmt := `INSERT INTO post (title, content, imageURL, author, created, expires, userid)
			VALUES ($1, $2, $3, $4, $5, $6, $7) returning id;`
	id := 0

	r.logger.Info("Inserting post for userID=%d", p.UserID)

	err := r.db.QueryRow(stmt, p.Title, p.Content, p.ImageURL, p.Username, p.Created, p.Expires, p.UserID).Scan(&id)
	if err != nil {
		r.logger.Error(err.Error())
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
		&p.Created,
		&p.Expires,
		&p.Username,
		&p.UserID,
	)
	if err != nil {
		r.logger.Error(err.Error())
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
		r.logger.Error(err.Error())
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
			&p.Created,
			&p.Expires,
			&p.Username,
			&p.UserID,
		)
		if err != nil {
			r.logger.Error(err.Error())
			return []domain.Post{}, err
		}
		posts = append(posts, p)
	}
	err = rows.Err()
	if err != nil {
		r.logger.Error(err.Error())
		return []domain.Post{}, err
	}

	return posts, nil
}

func (r *PostgresPostRepo) DeleteById(id int) error {
	stmt := `DELETE FROM post
			WHERE id = $1`

	_, err := r.db.Exec(stmt, id)
	if err != nil {
		r.logger.Error(err.Error())
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
		r.logger.Error(err.Error())
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
			&p.Created,
			&p.Expires,
			&p.Username,
			&p.UserID,
		)
		if err != nil {
			r.logger.Error(err.Error())
			return []domain.Post{}, err
		}
		posts = append(posts, p)
	}
	err = rows.Err()
	if err != nil {
		r.logger.Error(err.Error())
		return []domain.Post{}, err
	}

	return posts, nil
}

func (r *PostgresPostRepo) IsPostExpired(id int) error {
	return nil
}
