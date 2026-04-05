package domain

import (
	"time"
)

type Comment struct {
	CommentID  int
	PostID     int
	UserID     int
	ParentID   *int
	Content    string
	ImageURL   string
	Created    time.Time
	Username   string
	UserAvatar string
	Replies    []*Comment
}

type CreateCommentForm struct {
	PostID      int
	ParentID    *int
	Name        string
	Comment     string
	CommentFile UploadFile
}

func NewComment(usID, postID int, parentID *int, content, imageurl string) Comment {
	return Comment{
		UserID:   usID,
		PostID:   postID,
		Content:  content,
		ImageURL: imageurl,
		Created:  time.Now().UTC(),
		ParentID: parentID,
	}
}
