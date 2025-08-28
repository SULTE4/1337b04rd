package service

import (
	"1337b04rd/internal/appError"
	"1337b04rd/internal/core/domain"
	"1337b04rd/internal/core/ports"
	"fmt"
	"net/http"
)

type CommentService struct {
	repo  ports.CommentRepository
	userR ports.UserRepository
	s3    ports.S3Repository
}

func NewCommentService(repo ports.CommentRepository, userR ports.UserRepository, s3 ports.S3Repository) *CommentService {
	return &CommentService{repo: repo, userR: userR, s3: s3}
}

func (c *CommentService) AddComment(r *http.Request, com domain.CreateCommentForm) error {
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

	userid := r.Context().Value("user").(int)
	fmt.Println(userid)
	comment := domain.NewComment(userid, com.PostID, com.Comment, imageUrl)

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
	for i, comment := range comments {
		userData, err := c.userR.GetUserByID(comment.UserID)
		if err != nil {
			return []domain.Comment{}, err
		}

		comments[i].Username = userData.Name
		comments[i].UserAvatar = userData.AvatarURL
	}

	if err != nil {
		return []domain.Comment{}, err
	}

	return comments, nil
}
