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

	token, err := util.CreateJWT(user.Name, user.ImageURL, secret, expiration)
	if err != nil {
		return "", nil, err
	}

	// build claims that match what we encoded in JWT
	claims := &util.Claims{
		Username: user.Name,
		Avatar:   user.ImageURL,
		Exp:      time.Now().Add(expiration).Unix(),
	}

	// persist user in DB
	if err := s.repo.NewUser(domain.User{
		Name:      user.Name,
		AvatarURL: user.ImageURL,
		SessionID: token,
	}); err != nil {
		return "", nil, err
	}

	return token, claims, nil
}
