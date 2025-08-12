package post

import (
	"1337b04rd/internal/adapters/s3"
	"fmt"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreatePost(p CreatePostForm) (int, error) {
	imageUrl, err := s3.UploadObject(false, p.File)
	if err != nil {
		return 0, err
	}

	post := newPost(p.Subject, p.Comment, imageUrl, 99)

	// add handle to empty data in post form
	id, err := s.repo.Insert(*post)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) GetPost(id int) (Post, error) {
	p, err := s.repo.GetByID(id)
	if err != nil {
		return Post{}, err
	}
	return p, nil
}

func (s *Service) GetAll() ([]Post, error) {
	posts, err := s.repo.GetAll()
	if err != nil {
		return []Post{}, err
	}

	return posts, nil
}

func (s *Service) GetArchive() ([]Post, error) {
	posts, err := s.repo.GetExpiredPosts()
	if err != nil {
		return []Post{}, err
	}

	return posts, nil
}

func (s *Service) GetArchivePost(id int) (Post, error) {
	p, err := s.repo.GetByID(id)
	if err != nil {
		return Post{}, err
	}
	if p.Expires.After(time.Now()) {
		return Post{}, fmt.Errorf("post is not expired: %d", id)
	}
	return p, nil

}
