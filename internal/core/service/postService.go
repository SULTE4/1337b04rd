package service

import (
	"1337b04rd/internal/core/domain"
	"1337b04rd/internal/core/ports"
	"fmt"
	"time"
)

type PostService struct {
	repo ports.PostRepository
	s3   ports.S3Repository
}

func NewPostService(repo ports.PostRepository, s3 ports.S3Repository) *PostService {
	return &PostService{repo: repo, s3: s3}
}

func (s *PostService) CreatePost(p domain.CreatePostForm) (int, error) {
	imageUrl, err := s.s3.UploadObject(true, p.File) // true if it is post image or not
	if err != nil {
		return 0, err
	}

	// add handle to empty data in post form

	post := domain.NewPost(p.Subject, p.Comment, imageUrl, p.Name)

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
