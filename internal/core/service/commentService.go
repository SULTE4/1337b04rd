package service

import (
	"1337b04rd/internal/appError"
	"1337b04rd/internal/core/domain"
	"1337b04rd/internal/core/ports"
)

type CommentService struct {
	repo ports.CommentRepository
	s3   ports.S3Repository
}

func NewCommentService(repo ports.CommentRepository, s3 ports.S3Repository) *CommentService {
	return &CommentService{repo: repo, s3: s3}
}

func (c *CommentService) AddComment(com domain.CreateCommentForm) error {
	is, err := c.repo.IsPostExpired(com.PostID)
	if err != nil {
		return err
	}

	if is {
		return appError.ErrPostNotAvailable
	}

	imageUrl, err := c.s3.UploadObject(false, com.CommentFile)
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
