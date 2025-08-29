package router

import (
	ihttp "1337b04rd/internal/adapters/interface/http"
	"net/http"
)

func NewRouter(handler *ihttp.Handler) http.Handler {
	router := http.NewServeMux()

	// not sure at this aproach to handle image datas like this do you mind if you change ?
	router.Handle("GET /postsImg/", http.StripPrefix("/postsImg/",
		http.FileServer(http.Dir("./s3_storage/postsImg"))))

	router.Handle("GET /commentsImg/", http.StripPrefix("/commentsImg/",
		http.FileServer(http.Dir("./s3_storage/commentsImg"))))

	router.Handle("GET /", handler.LogRequest(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/catalog", http.StatusSeeOther)
	})))
	router.Handle("GET /catalog", handler.LogRequest(http.HandlerFunc(handler.Catalog)))
	router.Handle("GET /post/create", handler.LogRequest(http.HandlerFunc(handler.Create)))
	router.Handle("GET /post/{id}", handler.LogRequest(http.HandlerFunc(handler.ViewPost)))
	router.Handle("GET /archive", handler.LogRequest(http.HandlerFunc(handler.Archive)))
	router.Handle("GET /archive/post/{id}", handler.LogRequest(http.HandlerFunc(handler.ViewArchivePost)))
	// router.HandleFunc("GET /error", handler.Error)

	router.Handle("POST /submit-post", handler.LogRequest(handler.TokenMiddleware(http.HandlerFunc(handler.CreatePost))))
	router.Handle("POST /submit-comment/{id}", handler.LogRequest(handler.TokenMiddleware(http.HandlerFunc(handler.CreateComment))))
	return router
}
