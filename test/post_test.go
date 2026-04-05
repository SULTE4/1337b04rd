package main

import (
	"1337b04rd/internal/app"
	"1337b04rd/internal/core/domain"
	"1337b04rd/internal/core/service"
	"database/sql"
	"log/slog"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func setupPostService(t *testing.T) (*service.PostService, *sql.DB) {
	t.Helper()

	dsn := "postgres://postgres:pass@localhost:5432/mydb?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}
	err = db.Ping()
	if err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	store := app.NewStore(db, logger)

	s3Repo, err := s3.NewS3Repo()
	if err != nil {
		t.Fatalf("failed to init s3: %v", err)
	}

	postService := service.NewPostService(store.PostRepo, store.UserRepo, s3Repo)
	return postService, db
}

func TestCreateAndGetPost(t *testing.T) {
	service, db := setupPostService(t)
	defer db.Close()

	// Clean tables
	_, _ = db.Exec("DELETE FROM post")
	_, _ = db.Exec("DELETE FROM users")

	// Create a user
	userID := 0
	err := db.QueryRow(`INSERT INTO users (username, userurl, usertoken) VALUES ($1,$2,$3) RETURNING userid`,
		"testuser", "avatar.png", "token123").Scan(&userID)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Create post form
	form := domain.CreatePostForm{
		Subject: "Test Post",
		Comment: "This is a post content",
		Name:    "Test Author",
		File:    domain.UploadFile{Exist: false},
	}

	postID, err := service.CreatePost(userID, form)
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}

	// Get post and verify
	post, err := service.GetPost(postID)
	if err != nil {
		t.Fatalf("GetPost failed: %v", err)
	}
	if post.Title != "Test Post" || post.Content != "This is a post content" || post.Username != "testuser" {
		t.Fatalf("post data mismatch: %+v", post)
	}

	// Test GetAll returns at least one post
	posts, err := service.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}
	if len(posts) == 0 {
		t.Fatalf("expected at least 1 post, got 0")
	}

	// Expire post manually for archive test
	_, _ = db.Exec("UPDATE post SET expires = $1 WHERE id = $2", time.Now().UTC().Add(-1*time.Hour), postID)

	// Test GetArchive
	archived, err := service.GetArchive()
	if err != nil {
		t.Fatalf("GetArchive failed: %v", err)
	}
	found := false
	for _, p := range archived {
		if p.ID == postID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expired post not found in archive")
	}

	// Test GetArchivePost
	p, err := service.GetArchivePost(postID)
	if err != nil {
		t.Fatalf("GetArchivePost failed: %v", err)
	}
	if p.ID != postID {
		t.Fatalf("GetArchivePost returned wrong post id: %d", p.ID)
	}
}
