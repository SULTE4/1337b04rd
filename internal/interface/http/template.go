package http

import (
	"1337b04rd/internal/domain/post"
	"1337b04rd/web"
	"io/fs"
	"net/http"
	"path/filepath"
	"text/template"
)

type templateData struct {
	Post       post.Post
	Posts      []post.Post
	ErrID      int
	ErrMessage string
}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(web.Files, "templates/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).ParseFS(web.Files, page)
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
