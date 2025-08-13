package ports

import "1337b04rd/internal/core/domain"

type PostRepository interface {
	Insert(p domain.Post) (int, error)
	GetByID(id int) (domain.Post, error)
	GetAll() ([]domain.Post, error)
	DeleteById(id int) error
	GetExpiredPosts() ([]domain.Post, error)
}
