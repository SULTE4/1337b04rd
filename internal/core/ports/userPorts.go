package ports

import "1337b04rd/internal/core/domain"

type UserRepository interface {
	UpsertBySession(sessionID string, name, avatar string) (domain.User, error)
	GetBySession(sessionID string) (domain.User, error)
}
