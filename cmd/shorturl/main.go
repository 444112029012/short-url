package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/444112029012/short-url/internal/application"
	"github.com/444112029012/short-url/internal/config"
	"github.com/444112029012/short-url/internal/domain"
	"github.com/444112029012/short-url/internal/httpapi"
	"github.com/444112029012/short-url/internal/infra/sqlite"
	"github.com/444112029012/short-url/internal/observability"
	"github.com/444112029012/short-url/internal/ratelimit"
)

func main() {
	if err := run(); err != nil {
		log.Printf("shorturl stopped: %s", err.Error())
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	store, err := sqlite.Open(cfg.DatabasePath)
	if err != nil {
		return err
	}
	defer store.Close()

	ctx := context.Background()
	if err := store.Migrate(ctx); err != nil {
		return err
	}

	validator := domain.NewURLValidator()
	generator := domain.NewShortCodeGenerator()
	appSvc := application.NewShortURLApplicationService(validator, generator, store, store, cfg.BaseURL)
	redirectSvc := application.NewRedirectService(validator, store, store)
	adapter := httpapi.NewHTTPAPIAdapter(
		appSvc,
		redirectSvc,
		ratelimit.NewGuard(),
		httpapi.NewErrorMapper(),
		observability.NewHooks(),
	)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           httpapi.NewRouter(adapter),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("shorturl listening on %s", cfg.HTTPAddr)
		serveErr := srv.ListenAndServe()
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			errCh <- serveErr
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	select {
	case serveErr := <-errCh:
		return serveErr
	case <-sig:
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
