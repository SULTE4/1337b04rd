package comment

import "1337b04rd/internal/adapters/s3"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (c *Service) AddComment(com CreateCommentForm) error {
	imageUrl, err := s3.UploadObject(false, com.CommentFile)
	if err != nil {
		return err
	}

	comment := newComment(0, com.PostID, com.Comment, imageUrl)

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

func (c *Service) GetPostComments(id int) ([]Comment, error) {
	comments, err := c.repo.GetCommentsByPost(id)
	if err != nil {
		return []Comment{}, err
	}

	return comments, nil
}
