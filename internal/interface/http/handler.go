package http

import (
	"1337b04rd/internal/domain/post"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"text/template"
)

type Handler struct {
	postService *post.Service
	// userService   *user.Service
	templateCache map[string]*template.Template
}

type createPost struct {
	Name    string
	Subject string
	Comment string
	File    multipart.File
}

func NewHandler(postService *post.Service, templateCache map[string]*template.Template) *Handler {
	return &Handler{
		postService:   postService,
		templateCache: templateCache,
	}
}

func (h *Handler) Catalog(w http.ResponseWriter, r *http.Request) {
	// posts := h.postService.GetPost()
	w.Write([]byte("Showing posts"))
}

func (h *Handler) Archive(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Showing archived posts"))

}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	data := h.NewTemplateData(r)

	h.render(w, r, http.StatusOK, "create-post.html", data)
	w.Write([]byte("Create page"))
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.serverError(w, r, err)
		return
	}

	var post createPost

	post.Name = r.PostForm.Get("name")
	post.Subject = r.PostForm.Get("subject")
	post.Comment = r.PostForm.Get("comment")
	file, _, err := r.FormFile("file")
	if err != nil {
		h.serverError(w, r, err)
		return
	}
	post.File = file
	fmt.Println(post)
	// err = h.postService.CreatePost(post)
	w.Write([]byte("Create page submission"))

}

func (h *Handler) ViewPost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Showing post"))

}

func (h *Handler) Error(errId int, message string) http.Handler {
	data := h.NewTemplateData(r)

	h.render(w, r, http.StatusOK, "error.html", data)
}

func (h *Handler) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := h.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		h.serverError(w, r, err)
		return
	}
	w.WriteHeader(status)

	err := ts.Execute(w, nil)
	if err != nil {
		h.serverError(w, r, err)
		return
	}

}

func (h *Handler) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	slog.Error(err.Error(), slog.String("uri", uri), slog.String("method", method))
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
