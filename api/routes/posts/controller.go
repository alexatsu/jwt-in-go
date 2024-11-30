package posts

import (
	"github.com/gofiber/fiber/v3"

	m "github.com/alexatsu/jwt-in-go/api/middlewares"
)

func PostsRoutes(router fiber.Router) {
	posts := router.Group("/posts")
	posts.Get("/protected", getPosts, m.WithAuthMiddleware)
}
