package domain

import (
	"1337b04rd/internal/adapters/s3"
	"time"
)

type Comment struct {
	CommentID int
	UserID    int
	PostID    int
	Content   string
	ImageURL  string
	Created   time.Time
}

type CreateCommentForm struct {
	PostID      int
	Name        string
	Comment     string
	CommentFile s3.FileType
}

func NewComment(usID, postID int, content, imageurl string) Comment {
	return Comment{
		UserID:   usID,
		PostID:   postID,
		Content:  content,
		ImageURL: imageurl,
		Created:  time.Now().UTC(),
	}
}
