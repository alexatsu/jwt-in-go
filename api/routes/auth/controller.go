package auth

import "github.com/gofiber/fiber/v3"

func AuthRoutes(router fiber.Router) {
	auth := router.Group("/auth")
	auth.Post("/register", register)
	auth.Post("/login", login)
	auth.Post("/logout", logout)
}
