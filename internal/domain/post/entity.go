package post

import (
	"time"
)

type Post struct {
	ID      int
	Title   string
	Content string
	Image   string
	UserID  int
	Created time.Time
	Expires time.Time
}

func newPost(title, content, image string, id int) *Post {
	return &Post{
		Title:   title,
		Content: content,
		Image:   "-",
		UserID:  id,
		Created: time.Now().UTC(),
		Expires: time.Now().UTC().Add(time.Minute * 10),
	}
}
