package domain

import "time"

type User struct {
	SessionID string // cookie session id
	Name      string
	AvatarURL string
	CreatedAt time.Time
}

type Character struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ImageURL string `json:"image"`
}

func CreateUser(sessionid, name, avatarurl string) *User {
	return &User{
		SessionID: sessionid,
		Name:      name,
		AvatarURL: avatarurl,
		CreatedAt: time.Now().UTC(),
	}
}
