package service

import (
	"1337b04rd/internal/appError"
	"1337b04rd/internal/core/domain"
	"1337b04rd/internal/core/ports"
	"1337b04rd/internal/core/util"
)

type CommentService struct {
	repo  ports.CommentRepository
	userR ports.UserRepository
	s3    ports.S3Repository
}

func NewCommentService(repo ports.CommentRepository, userR ports.UserRepository, s3 ports.S3Repository) ports.CommentService {
	return &CommentService{repo: repo, userR: userR, s3: s3}
}

func (c *CommentService) AddComment(userID int, com domain.CreateCommentForm) error {
	if util.MinChars(com.Comment, 1) {
		return appError.ErrContentShouldNotBeEmpty
	}
	if util.MaxChars(com.Comment, 500) {
		return appError.ErrContentOutOfRange
	}

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

	comment := domain.NewComment(userID, com.PostID, com.ParentID, com.Comment, imageUrl)

	if err := c.repo.Insert(comment); err != nil {
		return err
	}

	if err := c.repo.UpdatePostExpire(com.PostID); err != nil {
		return err
	}

	return nil
}

func (c *CommentService) GetPostComments(id int) ([]*domain.Comment, error) {
	comments, err := c.repo.GetCommentsByPost(id)
	if err != nil {
		return nil, err
	}

	for i := range comments {
		userData, err := c.userR.GetUserByID(comments[i].UserID)
		if err != nil {
			return nil, err
		}
		comments[i].Username = userData.Name
		comments[i].UserAvatar = userData.AvatarURL
	}

	commentMap := make(map[int]*domain.Comment)
	var roots []*domain.Comment

	for i := range comments {
		comment := &comments[i]
		comment.Replies = []*domain.Comment{} // ensure non-nil slice
		commentMap[comment.CommentID] = comment
	}

	for i := range comments {
		comment := &comments[i]
		if comment.ParentID != nil {
			if parent, ok := commentMap[*comment.ParentID]; ok {
				parent.Replies = append(parent.Replies, comment)
			}
		} else {
			roots = append(roots, comment)
		}
	}

	return roots, nil
}
