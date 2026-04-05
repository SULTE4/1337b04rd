package ports

import (
	"1337b04rd/internal/core/domain"
	"1337b04rd/internal/core/util"
)

type CommentService interface {
	AddComment(userID int, com domain.CreateCommentForm) error
	GetPostComments(id int) ([]*domain.Comment, error)
}

type PostService interface {
	CreatePost(userID int, p domain.CreatePostForm) (int, error)
	GetPost(id int) (domain.Post, error)
	GetAll() ([]domain.Post, error)
	GetArchive() ([]domain.Post, error)
	GetArchivePost(id int) (domain.Post, error)
}

type UserService interface {
	NewUser(name string) (string, *util.Claims, error)
	Exists(id int) (bool, error)
	UpdateUsername(id int, newName string) error
}
