package router

import (
	ihttp "1337b04rd/internal/core/interface/http"
	"net/http"
)

func NewRouter(handler *ihttp.Handler) http.Handler {
	router := http.NewServeMux()

	// not sure at this aproach to handle iamge datas like this do mind if you change ?
	router.Handle("/postsImg/", http.StripPrefix("/postsImg/",
		http.FileServer(http.Dir("./s3_storage/postsImg"))))

	router.Handle("/commentsImg/", http.StripPrefix("/commentsImg/",
		http.FileServer(http.Dir("./s3_storage/commentsImg"))))

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
