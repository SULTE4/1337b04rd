package http

import (
	"1337b04rd/internal/adapters/s3"
	"1337b04rd/internal/domain/comment"
	"1337b04rd/internal/domain/post"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"text/template"
)

type Handler struct {
	postService *post.Service
	// userService   *user.Service
	commentService *comment.Service
	templateCache  map[string]*template.Template
	logger         *slog.Logger
}

func NewHandler(postService *post.Service, commentService *comment.Service, templateCache map[string]*template.Template, logger *slog.Logger) *Handler {
	return &Handler{
		postService:    postService,
		commentService: commentService,
		templateCache:  templateCache,
		logger:         logger,
	}
}

func (h *Handler) Catalog(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.GetAll()
	if err != nil {
		h.serverError(w, r, err)
		return
	}
	data := h.NewTemplateData(r)
	data.Posts = posts
	h.render(w, r, http.StatusSeeOther, "catalog.html", data)
}

func (h *Handler) Archive(w http.ResponseWriter, r *http.Request) {
	posts, err := h.postService.GetArchive()
	if err != nil {
		h.serverError(w, r, err)
		return
	}
	data := h.NewTemplateData(r)
	data.Posts = posts
	h.render(w, r, http.StatusSeeOther, "archive.html", data)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	data := h.NewTemplateData(r)

	h.render(w, r, http.StatusOK, "create-post.html", data)
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB в памяти
	if err != nil {
		h.serverError(w, r, err)
		return
	}
	defer r.MultipartForm.RemoveAll()

	var post post.CreatePostForm

	post.Name = r.PostForm.Get("name")
	post.Subject = r.PostForm.Get("subject")
	post.Comment = r.PostForm.Get("comment")

	file, handler, err := r.FormFile("file")
	if err != nil {
		h.serverError(w, r, err)
		return
	}

	post.File = s3.FileType{
		File:    file,
		Handler: handler,
	}
	id, err := h.postService.CreatePost(post)
	if err != nil {
		h.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)

}

func (h *Handler) ViewPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.serverError(w, r, err)
		return
	}
	p, err := h.postService.GetPost(id)
	if err != nil {
		h.serverError(w, r, err)
		return
	}
	comments, err := h.commentService.GetPostComments(id)
	if err != nil {
		h.serverError(w, r, err)
		return
	}

	data := h.NewTemplateData(r)
	data.Post = p
	data.Comments = comments

	h.render(w, r, http.StatusOK, "post.html", data)
}

func (h *Handler) ViewArchivePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.serverError(w, r, err)
		return
	}

	p, err := h.postService.GetArchivePost(id)
	if err != nil {
		h.serverError(w, r, err)
		return
	}

	data := h.NewTemplateData(r)
	data.Post = p

	h.render(w, r, http.StatusOK, "archive-post.html", data)

}

// func (h *Handler) Error(w http.ResponseWriter, r *http.Response, errId int, message string) {
// 	data := h.NewTemplateData(r)
// 	data.ErrID = errId
// 	data.ErrMessage = "occured error"
// 	h.render(w, r, http.StatusOK, "error.html", data)
// }

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB в памяти
	if err != nil {
		h.serverError(w, r, err)
		return
	}
	defer r.MultipartForm.RemoveAll()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.serverError(w, r, err)
		return
	}

	var com comment.CreateCommentForm

	// name will be taken from cookie (middleware)
	com.Comment = r.PostForm.Get("comment")
	com.PostID = id

	file, handler, err := r.FormFile("file")

	if err != nil {
		if err != http.ErrMissingFile {
			h.serverError(w, r, err)
			return
		}
	}
	if err == nil {
		com.CommentFile = s3.FileType{
			File:    file,
			Handler: handler,
		}
	}
	err = nil
	err = h.commentService.AddComment(com)
	if err != nil {
		h.serverError(w, r, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post/%d", id), http.StatusSeeOther)
}

func (h *Handler) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := h.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		h.serverError(w, r, err)
		return
	}
	w.WriteHeader(status)

	err := ts.Execute(w, data)
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
