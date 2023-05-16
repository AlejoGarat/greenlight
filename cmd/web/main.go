package main

import (
	"flag"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"

	healthcheckHandler "greenlight/internal/healthcheck/handlers"
	healthcheckRoutes "greenlight/internal/healthcheck/routes"
	moviesHandler "greenlight/internal/movies/handlers"
	moviesRoutes "greenlight/internal/movies/routes"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
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

	r := gin.Default()

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
