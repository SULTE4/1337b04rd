package app

import (
	"1337b04rd/internal/adapters/postgres"
	"1337b04rd/internal/core/ports"
	"database/sql"
)

type Store struct {
	PostRepo ports.PostRepository
	// UserRepo user.Repository
	CommentRepo ports.CommentRepository
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		PostRepo:    postgres.NewPostRepo(db),
		CommentRepo: postgres.NewCommentRepo(db),
	}
}
