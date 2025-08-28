package domain

import (
	"1337b04rd/internal/adapters/s3"
	"time"
)

type Post struct {
	ID         int
	Title      string
	Content    string
	ImageURL   string
	Username   string
	Created    time.Time
	Expires    time.Time
	UserID     int
	UserAvatar string
}

type CreatePostForm struct {
	Name    string
	Subject string
	Comment string
	File    s3.FileType
}

func NewPost(title, content, imageURL, username string, userid int) *Post {
	t := time.Now().UTC()
	return &Post{
		Title:    title,
		Content:  content,
		ImageURL: imageURL,
		Username: username,
		Created:  t,
		Expires:  t.Add(time.Minute * 10),
		UserID:   userid,
	}
}
