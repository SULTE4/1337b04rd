package app

import (
	"1337b04rd/internal/adapters/postgres"
	"1337b04rd/internal/domain/post"
	"database/sql"
)

type Store struct {
	PostRepo post.Repository
	// UserRepo user.Repository
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		PostRepo: postgres.NewPostRepo(db),
	}
}
