package main

import (
	"context"
	"log"
	"zhuchok/go-todo-server/ent"
	"zhuchok/go-todo-server/internal/auth"
	"zhuchok/go-todo-server/internal/middlewares"
	"zhuchok/go-todo-server/internal/tasks"
	"zhuchok/go-todo-server/internal/users"

	"entgo.io/ent/dialect"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found")
	}

	client, err := ent.Open(dialect.SQLite, "file:app.db?_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// migrations
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed create schema resource: %v", err)
	}

	app := fiber.New()

	authGroup := app.Group("/auth")
	auth.RegisterRoutes(authGroup, client)

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Use(func(c fiber.Ctx) error {
		return middlewares.AuthRequired(c, client)
	})

	tasks.RegisterRoutes(v1, client)
	users.RegisterRoutes(v1, client)

	log.Fatal(app.Listen(":3000"))

}
