package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"

	healthcheckHandler "greenlight/internal/healthcheck/handlers"
	healthcheckRoutes "greenlight/internal/healthcheck/routes"
	moviesHandler "greenlight/internal/movies/handlers"
	moviesRoutes "greenlight/internal/movies/routes"
	"greenlight/pkg/httphelpers"
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
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	var err error

	// dsn := "user=foo password=bar dbname=foobar host=localhost port=5432 sslmode=disable"

	// db, err := sqlx.Connect("postgres", dsn)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// helloWorldRepo := repo.NewSqlxRepo(db)
	// helloWorldRepo := repo.NewSqlxRepo(nil)
	// helloWorldLogic := logic.NewHelloWorldLogic(helloWorldRepo)
	// helloWorldHandler := handlers.NewHelloWorldHandler(helloWorldLogic)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
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

	moviesHandler := &moviesHandler.Handler{
		Logger:  logger,
		Version: version,
		Env:     "development",
	}

	v1 := r.Group("/v1")
	{
		healthcheckRoutes.MakeRoutes(v1, healthcheckHandler)
		moviesRoutes.MakeRoutes(v1, moviesHandler)
	}

	logger.Printf("starting %s server on %s", cfg.env, ":4000")

	err = r.Run(":4000")

	if err != nil {
		log.Fatal(err)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	// Set the maximum idle timeout.
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
