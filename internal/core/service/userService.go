package service

import (
	"database/sql"
	"errors"

	"1337b04rd/internal/core/domain"
	"1337b04rd/internal/core/ports"
)

type UserService struct {
	users ports.UserRepository
	rick  ports.RickAPI
}

func NewUserService(users ports.UserRepository, rick ports.RickAPI) *UserService {
	return &UserService{users: users, rick: rick}
}

// EnsureUser returns existing user by session or creates one with Rick&Morty avatar
func (s *UserService) EnsureUser(sessionID string) (domain.User, error) {
	u, err := s.users.GetBySession(sessionID)
	if err == nil {
		return u, nil
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, err
	}
	name, img, _, err := s.rick.RandomCharacter()
	if err != nil {
		return domain.User{}, err
	}
	return s.users.UpsertBySession(sessionID, name, img)
}
