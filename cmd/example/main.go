package main

import (
	"log"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"

	"greenlight/internal/helloworld/handlers"
	"greenlight/internal/helloworld/logic"
	"greenlight/internal/helloworld/repo"
)

func main() {
	var err error

	// dsn := "user=foo password=bar dbname=foobar host=localhost port=5432 sslmode=disable"

	// db, err := sqlx.Connect("postgres", dsn)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// helloWorldRepo := repo.NewSqlxRepo(db)
	helloWorldRepo := repo.NewSqlxRepo(nil)
	helloWorldLogic := logic.NewHelloWorldLogic(helloWorldRepo)
	helloWorldHandler := handlers.NewHelloWorldHandler(helloWorldLogic)

	r := gin.Default()

	makeRoutes(r, helloWorldHandler)

	err = r.Run(":4000")

	if err != nil {
		log.Fatal(err)
	}
}

func makeRoutes(r *gin.Engine, handler *handlers.GinHelloWorldHandler) {
	v1 := r.Group("/v1")
	{
		helloworld := v1.Group("/helloworld")
		{
			helloworld.POST("", handler.Greet())
			helloworld.GET("", handler.ListUsers())
			helloworld.GET("/:name", handler.GetUserByName())
		}
	}
}
