package app

import (
	"1337b04rd/internal/adapters/postgres"
	"1337b04rd/internal/core/ports"
	"database/sql"
	"log/slog"
)

type Store struct {
	PostRepo    ports.PostRepository
	UserRepo    ports.UserRepository
	CommentRepo ports.CommentRepository
}

func NewStore(db *sql.DB, logger *slog.Logger) *Store {
	return &Store{
		PostRepo:    postgres.NewPostRepo(db, logger),
		CommentRepo: postgres.NewCommentRepo(db, logger),
		UserRepo:    postgres.NewUserRepo(db, logger),
	}
}
