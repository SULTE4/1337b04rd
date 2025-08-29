package main

import (
	"1337b04rd/internal/adapters/s3"
	"1337b04rd/internal/app"
	"1337b04rd/internal/core/domain"
	"1337b04rd/internal/core/service"
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func setupApp(t *testing.T) (*service.CommentService, *sql.DB) {
	t.Helper()

	// Connect to test DB (make sure you have a separate test database)
	dsn := "postgres://postgres:pass@localhost:5432/mydb?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	store := app.NewStore(db, logger)
	s3Repo, err := s3.NewS3Repo()
	if err != nil {
		t.Fatalf("failed to init s3: %v", err)
	}

	commentService := service.NewCommentService(store.CommentRepo, store.UserRepo, s3Repo)

	return commentService, db
}

func TestAddCommentAndGetComments(t *testing.T) {
	service, db := setupApp(t)
	defer db.Close()

	// Clean tables before test
	_, _ = db.Exec("DELETE FROM comment")
	_, _ = db.Exec("DELETE FROM post")
	_, _ = db.Exec("DELETE FROM users")

	// Create a test user
	userID := 0
	err := db.QueryRow(`INSERT INTO users (username, userurl, usertoken) VALUES ($1, $2, $3) RETURNING userid`,
		"testuser", "avatar.png", "token123").Scan(&userID)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	// Create a test post
	postID := 0
	err = db.QueryRow(`INSERT INTO post (title, content, imageURL, author, created, expires, userid)
					 VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		"Test Post", "Hello", "image.png", "testuser", time.Now(), time.Now().Add(1*time.Hour), userID).Scan(&postID)
	if err != nil {
		t.Fatalf("failed to create test post: %v", err)
	}

	// Create HTTP request with context containing user ID
	req := &http.Request{}
	ctx := context.WithValue(req.Context(), "user", userID)
	req = req.WithContext(ctx)

	// Add comment
	commentForm := domain.CreateCommentForm{
		PostID:      postID,
		Comment:     "Test Comment",
		ParentID:    nil,
		CommentFile: s3.FileType{Exist: false},
	}

	err = service.AddComment(req, commentForm)
	if err != nil {
		t.Fatalf("AddComment failed: %v", err)
	}

	// Verify comment inserted
	comments, err := service.GetPostComments(postID)
	if err != nil {
		t.Fatalf("GetPostComments failed: %v", err)
	}
	if len(comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(comments))
	}
	if comments[0].Content != "Test Comment" {
		t.Fatalf("unexpected comment content: %s", comments[0].Content)
	}

	// Add a reply
	replyForm := domain.CreateCommentForm{
		PostID:      postID,
		Comment:     "Reply Comment",
		ParentID:    &comments[0].CommentID,
		CommentFile: s3.FileType{Exist: false},
	}
	err = service.AddComment(req, replyForm)
	if err != nil {
		t.Fatalf("AddComment reply failed: %v", err)
	}

	// Verify reply tree
	comments, err = service.GetPostComments(postID)
	if err != nil {
		t.Fatalf("GetPostComments failed: %v", err)
	}
	if len(comments[0].Replies) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(comments[0].Replies))
	}
	if comments[0].Replies[0].Content != "Reply Comment" {
		t.Fatalf("unexpected reply content: %s", comments[0].Replies[0].Content)
	}
}
