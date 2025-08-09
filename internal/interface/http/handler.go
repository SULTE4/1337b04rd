package http

import (
	"1337b04rd/internal/domain/post"
	"net/http"
)

type Handler struct {
	postService *post.Service
}

func NewHandler(postService *post.Service) *Handler {
	return &Handler{
		postService: postService,
	}
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {}
