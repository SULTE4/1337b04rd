package post

import (
	"time"
)

type Post struct {
	ID       int
	Title    string
	Content  string
	ImageURL string
	UserID   int
	Created  time.Time
	Expires  time.Time
}

func newPost(title, content, imageURL string, id int) *Post {
	t := time.Now().UTC()
	return &Post{
		Title:    title,
		Content:  content,
		ImageURL: imageURL,
		UserID:   id,
		Created:  t,
		Expires:  t.Add(time.Minute * 10),
	}
}
