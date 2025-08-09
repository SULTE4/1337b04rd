package post

import "mime/multipart"

type CreatePostForm struct {
	Name    string
	Subject string
	Comment string
	File    multipart.File
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreatePost(p CreatePostForm) (int, error) {

	post := newPost(p.Subject, p.Comment, "", 99)
	id, err := s.repo.Insert(*post)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Service) GetPost(id string) (*Post, error) {
	return &Post{}, nil
}

func (s *Service) GetAll() ([]Post, error) {
	// s.repo.DeleteById(?
	return []Post{}, nil
}
