package auth

import (
	"context"
	"github.com/alexatsu/jwt-in-go/db"
	"github.com/alexatsu/jwt-in-go/db/sqlc/generated/session"
	"github.com/alexatsu/jwt-in-go/utils"
	"github.com/gofiber/fiber/v3"
	"log"
)

var ctx = context.Background()

func WithAuthMiddleware(c fiber.Ctx) error {
	log.Println("withAuthMiddleware")
	accessToken := c.Cookies(utils.AccessTokenKey)
	if accessToken == "" {
		log.Println("Access token missing")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Access token missing"})
	}

	claims, err := utils.ValidateToken(accessToken)
	if err != nil {
		log.Println("Invalid access token")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid access token"})
	}

	if utils.IsTokenExpired(claims) {
		conn, err := db.AcqureConnection(ctx)
		if err != nil {
			return err
		}
		defer conn.Release()

		sessionQuery := session.New(conn)
		sessionData, err := sessionQuery.GetSession(ctx, claims.UserId)
		if err != nil {
			log.Println("Session not found")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Session not found"})
		}

		claims, err = utils.ValidateToken(sessionData.RefreshToken)
		if err != nil {
			log.Println("Invalid session")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid session"})
		}

		if utils.IsTokenExpired(claims) {
			log.Println("Session expired")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Session expired"})
		}

		newAccessToken, newRefreshToken, err := utils.GenerateTokens(claims.UserId)
		if err != nil {
			log.Println("Failed to refresh token")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to refresh token"})
		}

		err = sessionQuery.CreateSession(ctx, session.CreateSessionParams{
			UserID:       claims.UserId,
			RefreshToken: newRefreshToken,
		})
		if err != nil {
			log.Println("Failed to create session")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create session"})
		}

		c.Cookie(&fiber.Cookie{
			Name:     utils.AccessTokenKey,
			Value:    newAccessToken,
			Expires:  utils.TimeOneMonth,
			HTTPOnly: true,
			Secure:   true,
		})

		log.Println("Session refreshed")
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Session refreshed",
		})
	}

	return c.Next()
}
