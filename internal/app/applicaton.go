package app

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"1337b04rd/internal/adapters/s3"
	ihttp "1337b04rd/internal/core/interface/http"
	"1337b04rd/internal/core/interface/router"
	"1337b04rd/internal/core/service"
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

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	templateCache, err := ihttp.NewTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	// Init Store (adapters)
	store := NewStore(db)
	// S3 storage init
	err = s3.InitS3()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Init Domain Services
	postService := service.NewPostService(store.PostRepo)
	commentService := service.NewCommentService(store.CommentRepo)
	// userService := user.NewService(store.UserRepo)

	// Init Handlers
	handler := ihttp.NewHandler(postService, commentService, templateCache, logger)
	router := router.NewRouter(handler)

	return &Application{
		DB:     db,
		Store:  store,
		Logger: logger,
		Router: &router,
	}, nil
}
