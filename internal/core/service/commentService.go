package service

import (
	"1337b04rd/internal/adapters/s3"
	"1337b04rd/internal/appError"
	"1337b04rd/internal/core/domain"
	"1337b04rd/internal/core/ports"
)

type CommentService struct {
	repo ports.CommentRepository
}

func NewCommentService(repo ports.CommentRepository) *CommentService {
	return &CommentService{repo: repo}
}

func (c *CommentService) AddComment(com domain.CreateCommentForm) error {
	is, err := c.repo.IsPostExpired(com.PostID)
	if err != nil {
		return err
	}

	if is {
		return appError.ErrPostNotAvailable
	}

	imageUrl, err := s3.UploadObject(false, com.CommentFile)
	if err != nil {
		return err
	}

	comment := domain.NewComment(0, com.PostID, com.Comment, imageUrl)

	err = c.repo.Insert(comment)
	if err != nil {
		return err
	}
	err = c.repo.UpdatePostExpire(com.PostID)
	if err != nil {
		return err
	}

	return nil
}

func (c *CommentService) GetPostComments(id int) ([]domain.Comment, error) {
	comments, err := c.repo.GetCommentsByPost(id)
	if err != nil {
		return []domain.Comment{}, err
	}

	return comments, nil
}
