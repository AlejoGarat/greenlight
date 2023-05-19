package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"

	healthcheckHandler "greenlight/internal/healthcheck/handlers"
	healthcheckRoutes "greenlight/internal/healthcheck/routes"
	moviesHandler "greenlight/internal/movies/handlers"
	moviesRepo "greenlight/internal/movies/repository"
	moviesRoutes "greenlight/internal/movies/routes"
	moviesService "greenlight/internal/movies/service"
	usersHandler "greenlight/internal/users/handlers"
	usersRepo "greenlight/internal/users/repo"
	userRoutes "greenlight/internal/users/routes"
	usersService "greenlight/internal/users/service"
	"greenlight/pkg/httphelpers"
	"greenlight/pkg/jsonlog"
	"greenlight/pkg/middlewares"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")
	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	var err error

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	r := gin.Default()
	r.NoRoute(gin.HandlerFunc(httphelpers.StatusNotFoundResponse))
	r.NoMethod(gin.HandlerFunc(httphelpers.StatusMethodNotAllowedResponse))

	healthcheckHandler := &healthcheckHandler.Handler{
		Logger:  logger,
		Version: version,
		Env:     "development",
	}

	mr := moviesRepo.NewMovieRepo(db)
	ms := moviesService.NewMovieService(mr)

	moviesHandler := &moviesHandler.Handler{
		Logger:       logger,
		Version:      version,
		Env:          "development",
		MovieService: ms,
	}

	ur := usersRepo.NewMovieRepo(db)
	us := usersService.NewUserService(ur)

	usersHandler := &usersHandler.Handler{
		Logger:      logger,
		Version:     version,
		Env:         "development",
		UserService: us,
	}

	r.Use(middlewares.RecoverPanic())
	r.Use(middlewares.RateLimit(cfg.limiter.enabled))
	v1 := r.Group("/v1")
	{
		healthcheckRoutes.MakeRoutes(v1, healthcheckHandler)
		moviesRoutes.MakeRoutes(v1, moviesHandler)
		userRoutes.MakeRoutes(v1, usersHandler)
	}

	info := info{
		cfg:                &cfg,
		logger:             logger,
		HealthcheckHandler: healthcheckHandler,
		MoviesHandler:      moviesHandler,
	}

	err = Serve(info, r)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
