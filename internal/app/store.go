package app

import (
	"1337b04rd/internal/adapters/postgres"
	"1337b04rd/internal/domain/comment"
	"1337b04rd/internal/domain/post"
	"database/sql"
)

type Store struct {
	PostRepo post.Repository
	// UserRepo user.Repository
	CommentRepo comment.Repository
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		PostRepo:    postgres.NewPostRepo(db),
		CommentRepo: postgres.NewCommentRepo(db),
	}
}
