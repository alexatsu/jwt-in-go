package main

import (
	"log"

	"github.com/alexatsu/jwt-in-go/api/routes/auth"
	"github.com/alexatsu/jwt-in-go/api/routes/posts"
	"github.com/alexatsu/jwt-in-go/db"
	"github.com/alexatsu/jwt-in-go/utils"
	"github.com/gofiber/fiber/v3"
)

func main() {
	utils.LoadEnv()

	db.InitStore()

	app := fiber.New()
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	v1 := app.Group("/v1")
	auth.AuthRoutes(v1)
	posts.PostsRoutes(v1)

	log.Fatal(app.Listen(":3010"))
}
