package service

import (
	"1337b04rd/internal/appError"
	"1337b04rd/internal/core/domain"
	"1337b04rd/internal/core/ports"
	"1337b04rd/internal/core/util"
	"fmt"
	"net/http"
	"time"
)

type PostService struct {
	repo  ports.PostRepository
	userR ports.UserRepository
	s3    ports.S3Repository
}

func NewPostService(repo ports.PostRepository, userR ports.UserRepository, s3 ports.S3Repository) *PostService {
	return &PostService{repo: repo, userR: userR, s3: s3}
}

func (s *PostService) CreatePost(r *http.Request, p domain.CreatePostForm) (int, error) {
	imageUrl, err := s.s3.UploadObject(true, p.File) // true if it is post image or not
	if err != nil {
		return 0, err
	}

	userid := r.Context().Value("user").(int)

	if util.MaxChars(p.Subject, 50) {
		return 0, appError.ErrTitleOutOfRange
	}
	if util.MinChars(p.Subject, 3) {
		return 0, appError.ErrTitleTooShort
	}
	if util.MinChars(p.Comment, 1) {
		return 0, appError.ErrContentShouldNotBeEmpty
	}
	if util.MaxChars(p.Comment, 500) {
		return 0, appError.ErrContentOutOfRange
	}

	post := domain.NewPost(p.Subject, p.Comment, imageUrl, p.Name, userid)

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
	userData, err := s.userR.GetUserByID(p.UserID)
	if err != nil {
		return domain.Post{}, err
	}
	p.Username = userData.Name
	p.UserAvatar = userData.AvatarURL

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
