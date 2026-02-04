package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apphttp "pack-calculator/internal/http"
	"pack-calculator/internal/repo"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	_ "modernc.org/sqlite"
)

func main() {
	port := getEnv("PORT", "8080")
	dbPath := getEnv("DB_PATH", "data/packs.db")
	webDir := getEnv("WEB_DIR", "web")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal().Err(err).Msg("open db")
	}
	defer db.Close()
	db.SetMaxOpenConns(1)

	repo, err := repo.NewSQLiteRepo(db)
	if err != nil {
		log.Fatal().Err(err).Msg("repo init")
	}

	if err := ensureDefaultPackSizes(context.Background(), repo); err != nil {
		log.Fatal().Err(err).Msg("seed pack sizes")
	}

	api := apphttp.NewRouter(repo)
	mux := http.NewServeMux()
	mux.Handle("/v1/", api)
	mux.Handle("/healthz", api)
	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	log.Logger = logger

	h := hlog.NewHandler(logger)(mux)
	h = hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", status).
			Int("bytes", size).
			Dur("duration", duration).
			Msg("request")
	})(h)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: h,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case <-stop:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("server error")
		}
	}
}

func ensureDefaultPackSizes(ctx context.Context, repo repo.PackSizeRepository) error {
	sizes, err := repo.ListPackSizes(ctx)
	if err != nil {
		return err
	}
	if len(sizes) > 0 {
		return nil
	}
	return repo.ReplacePackSizes(ctx, []int{250, 500, 1000, 2000, 5000})
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
