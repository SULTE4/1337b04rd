package main

import (
	"1337b04rd/internal/app"
	"flag"
	"fmt"
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
	addr := flag.String("port", "8080", "Port number.")
	help := flag.Bool("help", false, "Show this screen.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "hacker board\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  1337b04rd [--port <N>]\n")
		fmt.Fprintf(os.Stderr, "  1337b04rd --help\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	*addr = ":" + *addr

	dsn := "postgres://postgres:pass@localhost:5432/mydb?sslmode=disable"

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	app, err := app.NewApplication(dsn, *logger)
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
