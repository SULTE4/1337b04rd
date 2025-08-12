package http

import (
	"1337b04rd/internal/domain/comment"
	"1337b04rd/internal/domain/post"
	"1337b04rd/web"
	"io/fs"
	"net/http"
	"path/filepath"
	"text/template"
	"time"
)

type templateData struct {
	Post       post.Post
	Posts      []post.Post
	Comments   []comment.Comment
	ErrID      int
	ErrMessage string
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(web.Files, "templates/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFS(web.Files, page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}

func (h *Handler) NewTemplateData(r *http.Request) templateData {
	return templateData{}
}
