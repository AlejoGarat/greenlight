package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"greenlight/pkg/jsonlog"

	healthcheckHandler "greenlight/internal/healthcheck/handlers"
	moviesHandler "greenlight/internal/movies/handlers"

	"github.com/gin-gonic/gin"
)

type info struct {
	cfg                *config
	logger             *jsonlog.Logger
	HealthcheckHandler *healthcheckHandler.Handler
	MoviesHandler      *moviesHandler.Handler
}

func Serve(info info, r *gin.Engine) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", info.cfg.port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		info.logger.PrintInfo("shutting down server", map[string]string{
			"signal": s.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	info.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  info.cfg.env,
	})

	err := r.Run(srv.Addr)
	if err != nil {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	info.logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
