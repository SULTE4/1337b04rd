package router

import (
	ihttp "1337b04rd/internal/interface/http"
	"net/http"
)

func NewRouter(handler *ihttp.Handler) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /catalog", handler.Catalog)
	router.HandleFunc("GET /post/create", handler.Create)
	router.HandleFunc("GET /post/{id}", handler.ViewPost)
	router.HandleFunc("GET /archive", handler.Archive)
	router.HandleFunc("GET /archive/post/{id}", handler.ViewArchivePost)
	// router.HandleFunc("GET /error", handler.Error)

	router.HandleFunc("POST /submit-post", handler.CreatePost)
	router.HandleFunc("POST /submit-comment/{id}", handler.CreateComment)
	return router
}
