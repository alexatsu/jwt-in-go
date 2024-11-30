package utils

import (
	"net/mail"

	"github.com/gofiber/fiber/v3"
)

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GetValidCredentials(c fiber.Ctx) (*Credentials, error) {
	creds := new(Credentials)
	if err := c.Bind().Body(creds); err != nil {
		return nil, err
	}

	if creds.Email == "" || creds.Password == "" {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Missing email or password"})
	}

	if !validateEmail(creds.Email) {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid email"})
	}

	return creds, nil
}
