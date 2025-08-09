package router

import (
	ihttp "1337b04rd/internal/interface/http"
	"net/http"
)

func NewRouter(handler *ihttp.Handler) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /", handler.Catalog)
	router.HandleFunc("GET /post/create", handler.Create)
	router.HandleFunc("GET /post/{id}", handler.ViewPost)
	router.HandleFunc("GET /archive/post", handler.Archive)
	router.HandleFunc("GET /archive/view/post", handler.ViewPost)
	// router.HandleFunc("GET /error", handler.Error)

	router.HandleFunc("POST /submit-post", handler.CreatePost)
	router.HandleFunc("POST /post/{id}", handler.CommentPost)
	return router
}
