package service

import (
	"1337b04rd/internal/core/domain"
	"1337b04rd/internal/core/ports"
	"1337b04rd/internal/core/util"
	"net/http"
	"os"
	"time"
)

type UserService struct {
	repo        ports.UserRepository
	externalApi ports.ExternalApi
}

func NewUserService(repo ports.UserRepository, ex ports.ExternalApi) *UserService {
	return &UserService{repo: repo, externalApi: ex}
}

func (s *UserService) NewUser(r *http.Request) (string, *util.Claims, error) {
	if err := r.ParseForm(); err != nil {
		return "", nil, err
	}

	occupied, err := s.repo.GetOccupiedCharacters()
	if err != nil {
		return "", nil, err
	}

	user, err := s.externalApi.GetRandomCharacter(occupied)
	if err != nil {
		return "", nil, err
	}

	if name := r.PostForm.Get("name"); name != "" {
		user.Name = name
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "supersecrethahaha"
	}

	expiration := 7 * 24 * time.Hour

	var userID int
	if userID, err = s.repo.NewUser(domain.User{
		Name:      user.Name,
		AvatarURL: user.ImageURL,
	}); err != nil {
		return "", nil, err
	}

	token, err := util.CreateJWT(user.Name, user.ImageURL, secret, userID, expiration)
	if err != nil {
		return "", nil, err
	}

	if err := s.repo.UpdateUserToken(userID, token); err != nil {
		return "", nil, err
	}

	claims := &util.Claims{
		UserID:   userID,
		Username: user.Name,
		Avatar:   user.ImageURL,
		Exp:      time.Now().Add(expiration).Unix(),
	}
	return token, claims, nil
}

func (s *UserService) GetUserIDByToken(token string) (int, error) {
	userData, err := s.repo.GetUserByToken(token)
	if err != nil {
		return 0, err
	}

	return userData.UserID, nil
}

func (s *UserService) UpdateUsername(id int, newName string) error {
	return s.repo.UpdateName(id, newName)
}

func (s *UserService) Exists(id int) (bool, error) {
	return s.repo.Exists(id)
}
