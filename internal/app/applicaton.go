package app

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"1337b04rd/internal/domain/post"
	ihttp "1337b04rd/internal/interface/http"
	"1337b04rd/internal/interface/router"
)

type Application struct {
	DB     *sql.DB
	Store  *Store
	Logger *slog.Logger
	Router *http.Handler
}

func NewApplication(dsn string) (*Application, error) {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	// Connect DB
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Init Store (adapters)
	store := NewStore(db)

	// Init Domain Services
	postService := post.NewService(store.PostRepo)
	// userService := user.NewService(store.UserRepo)

	// Init Handlers
	handler := ihttp.NewHandler(postService)
	router := router.NewRouter(handler)

	return &Application{
		DB:     db,
		Store:  store,
		Logger: logger,
		Router: &router,
	}, nil
}
