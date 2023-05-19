package main

import (
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

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit
		info.logger.PrintInfo("caught signal", map[string]string{
			"signal": s.String(),
		})
		os.Exit(0)
	}()

	info.logger.PrintInfo("starting server", map[string]string{
		"addr": ":4000",
		"env":  info.cfg.env,
	})

	return r.Run(srv.Addr)
}
