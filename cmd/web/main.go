package main

import (
	"1337b04rd/internal/app"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

/*	TODO
	cookie,
	config (optional),
	user authentication

*/

func main() {
	addr := flag.String("port", ":8080", "server address")
	flag.Parse()

	dsn := "postgres://postgres:pass@localhost:5432/mydb?sslmode=disable"

	app, err := app.NewApplication(dsn)
	if err != nil {
		log.Fatalf("failed to start app: %v", err)
	}

	srv := &http.Server{
		Addr:         *addr,
		Handler:      *app.Router,
		ErrorLog:     slog.NewLogLogger(app.Logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.Logger.Info("Server running at", slog.String("addr", *addr))
	err = srv.ListenAndServe()

	app.Logger.Error(err.Error())
	os.Exit(1)
}
