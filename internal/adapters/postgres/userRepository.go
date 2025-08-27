package postgres

import (
	"1337b04rd/internal/core/domain"
	"database/sql"
	"log/slog"
)

type PostgresUserRepo struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewUserRepo(db *sql.DB, logger *slog.Logger) *PostgresUserRepo {
	return &PostgresUserRepo{db: db, logger: logger}
}

// написать тест на переполнение персонажев
func (r *PostgresUserRepo) GetOccupiedCharacters() ([]int, error) {
	stmt := `SELECT userid from comment`

	rows, err := r.db.Query(stmt)
	if err != nil {
		r.logger.Error(err.Error())
		return []int{}, err
	}

	var occupied []int
	for rows.Next() {
		var id int
		err = rows.Scan(
			&id,
		)
		if err != nil {
			r.logger.Error(err.Error())
			return []int{}, err
		}
		occupied = append(occupied, id)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error(err.Error())
		return []int{}, err
	}

	return occupied, nil
}
func (r *PostgresUserRepo) NewUser(user domain.User) error {
	stmt := `INSERT INTO users (username, userurl, usertoken) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(stmt, user.Name, user.AvatarURL, user.SessionID)
	return err
}
