package domain

import "time"

type User struct {
	SessionID string // cookie session id
	Name      string
	AvatarURL string
	CreatedAt time.Time
}
