package router

import (
	ihttp "1337b04rd/internal/interface/http"
	"net/http"
)

func NewRouter(handler *ihttp.Handler) http.Handler {
	router := http.NewServeMux()


	router.HandleFunc("GET /", handler.home)
	router.Handle("GET /post/create", )
	router.Handle("POST /submit-post", )
	router.Handle("GET /post/view",)
	router.Handle("GET /archive/post",)
	router.Handle("GET /archive/view/post" handler.)
	router.Handle("GET /error", )
	return router
}
