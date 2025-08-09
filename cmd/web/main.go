package main

import (
	"1337b04rd/internal/app"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	dsn := "postgres://user:pass@localhost:5432/mydb?sslmode=disable"

	app, err := app.NewApplication(dsn)
	if err != nil {
		log.Fatalf("failed to start app: %v", err)
	}

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      *app.Router,
		ErrorLog:     slog.NewLogLogger(app.Logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.Logger.Info("Server running at :8080")
	err = srv.ListenAndServe()

	app.Logger.Error(err.Error())
	os.Exit(1)
}
