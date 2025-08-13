package ports

import "1337b04rd/internal/core/domain"

type CommentRepository interface {
	Insert(com domain.Comment) error
	GetCommentsByPost(id int) ([]domain.Comment, error)
	UpdatePostExpire(id int) error
	IsPostExpired(id int) (bool, error)
}
