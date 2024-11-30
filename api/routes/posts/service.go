package posts

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

func getPosts(c fiber.Ctx) error {
	log.Println("Access granted to protected route")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Access granted to protected route",
		"posts":   []string{"Post 1", "Post 2", "Post 3"},
	})
}
