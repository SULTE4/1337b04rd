package app

import (
	"1337b04rd/internal/adapters/external"
	"1337b04rd/internal/adapters/interface/router"
	"1337b04rd/internal/adapters/s3"
	"1337b04rd/internal/core/ports"
	"1337b04rd/internal/core/service"
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	ihttp "1337b04rd/internal/adapters/interface/http"
)

type Application struct {
	DB     *sql.DB
	Store  *Store
	Logger *slog.Logger
	Router *http.Handler
}

type Services struct {
	Post    ports.PostService
	Comment ports.CommentService
	User    ports.UserService
}

func NewApplication(dsn string, logger slog.Logger) (*Application, error) {
	// Connect DB
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// parse to prepare html file in advance
	templateCache, err := ihttp.NewTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Init Store (adapters)
	store := NewStore(db, &logger)

	// S3 storage init
	s3Storage, err := s3.NewS3Repo()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	//
	externalApi := external.NewRandMApi(logger)

	// Init Domain Services
	postService := service.NewPostService(store.PostRepo, store.UserRepo, s3Storage)
	commentService := service.NewCommentService(store.CommentRepo, store.UserRepo, s3Storage)
	userService := service.NewUserService(store.UserRepo, externalApi)
	services := Services{
		Post:    postService,
		Comment: commentService,
		User:    userService,
	}

	// Init Handlers
	handler := ihttp.NewHandler(services.Post, services.Comment, services.User, templateCache, &logger)
	router := router.NewRouter(handler)

	return &Application{
		DB:     db,
		Store:  store,
		Logger: &logger,
		Router: &router,
	}, nil
}
