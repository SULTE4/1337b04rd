package service

import (
	"1337b04rd/internal/adapters/s3"
	"1337b04rd/internal/core/domain"
	"1337b04rd/internal/core/ports"
	"fmt"
	"time"
)

type PostService struct {
	repo ports.PostRepository
}

func NewPostService(repo ports.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(p domain.CreatePostForm) (int, error) {
	imageUrl, err := s3.UploadObject(true, p.File) // true if it is post image or not
	if err != nil {
		return 0, err
	}

	post := domain.NewPost(p.Subject, p.Comment, imageUrl, 99)

	// add handle to empty data in post form
	id, err := s.repo.Insert(*post)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *PostService) GetPost(id int) (domain.Post, error) {
	p, err := s.repo.GetByID(id)
	if err != nil {
		return domain.Post{}, err
	}
	return p, nil
}

func (s *PostService) GetAll() ([]domain.Post, error) {
	posts, err := s.repo.GetAll()
	if err != nil {
		return []domain.Post{}, err
	}

	return posts, nil
}

func (s *PostService) GetArchive() ([]domain.Post, error) {
	posts, err := s.repo.GetExpiredPosts()
	if err != nil {
		return []domain.Post{}, err
	}

	return posts, nil
}

func (s *PostService) GetArchivePost(id int) (domain.Post, error) {
	p, err := s.repo.GetByID(id)
	if err != nil {
		return domain.Post{}, err
	}
	if p.Expires.After(time.Now()) {
		return domain.Post{}, fmt.Errorf("post is not expired: %d", id)
	}
	return p, nil

}
