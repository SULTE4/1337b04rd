package postgres

import (
	"1337b04rd/internal/core/domain"
	"database/sql"
	"fmt"
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

func (r *PostgresUserRepo) NewUser(user domain.User) (int, error) {
	stmt := `INSERT INTO users (username, userurl, usertoken) VALUES ($1, $2, $3) returning userid`
	var id int
	err := r.db.QueryRow(stmt, user.Name, user.AvatarURL, user.SessionID).Scan(&id)
	fmt.Println("here is user id from newUser")
	return id, err
}

func (r *PostgresUserRepo) GetUserByID(id int) (domain.User, error) {
	stmt := `
    SELECT userid, username, userurl, usertoken
    FROM users
    WHERE userid = $1
`

	row := r.db.QueryRow(stmt, id)
	var user domain.User
	err := row.Scan(
		&user.UserID,
		&user.Name,
		&user.AvatarURL,
		&user.SessionID,
	)
	if err != nil {
		r.logger.Error(err.Error())
		return domain.User{}, err
	}
	return user, nil
}

func (r *PostgresUserRepo) GetUserByToken(token string) (domain.User, error) {
	stmt := `SELECT userID, username, userURL, userToken FROM users WHERE userToken = $1`

	fmt.Println("changing by token", token)
	row := r.db.QueryRow(stmt, token)

	var user domain.User
	err := row.Scan(
		&user.UserID,
		&user.Name,
		&user.AvatarURL,
		&user.SessionID,
	)
	if err != nil {
		r.logger.Error(err.Error())
		return domain.User{}, err
	}
	fmt.Println(user)
	return user, nil
}

func (r *PostgresUserRepo) UpdateName(id int, newName string) error {
	stmt := `UPDATE users 
			SET username = $1
			where userid = $2`

	fmt.Println("updating name")
	_, err := r.db.Exec(stmt, newName, id)
	if err != nil {

		r.logger.Error(err.Error(), slog.String("comment=", "not okay in User repo update name"))
		return err
	}

	return nil
}

func (r *PostgresUserRepo) UpdateUserToken(id int, token string) error {
	stmt := `UPDATE users 
			SET usertoken = $1
			where userid = $2`

	fmt.Println("updating token")
	_, err := r.db.Exec(stmt, token, id)
	if err != nil {

		r.logger.Error(err.Error(), slog.String("comment=", "not okay in User repo update token"))
		return err
	}

	return nil
}

func (r *PostgresUserRepo) Exists(id int) (bool, error) {
	stmt := `SELECT EXISTS(SELECT 1 FROM users WHERE userid = $1)`

	var exists bool
	err := r.db.QueryRow(stmt, id).Scan(&exists)
	if err != nil {
		r.logger.Error(err.Error())
		return false, err
	}

	return exists, nil
}
