package ports

import (
	"1337b04rd/internal/core/domain"
)

type UserRepository interface {
	GetOccupiedCharacters() ([]int, error)
	NewUser(user domain.User) error
}
