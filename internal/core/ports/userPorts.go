package ports

import (
	"1337b04rd/internal/core/domain"
)

type UserRepository interface {
	GetOccupiedCharacters() ([]int, error)
	NewUser(user domain.User) (int, error)
	Exists(id int) (bool, error)
	GetUserByID(id int) (domain.User, error)
	GetUserByToken(token string) (domain.User, error)
	UpdateName(userID int, newName string) error
	UpdateUserToken(userId int, token string) error
}
