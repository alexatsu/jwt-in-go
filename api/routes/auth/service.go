package auth

import (
	"context"
	"log"
	"time"

	"github.com/alexatsu/jwt-in-go/db"
	"github.com/alexatsu/jwt-in-go/db/sqlc/generated/session"
	"github.com/alexatsu/jwt-in-go/db/sqlc/generated/user"
	"github.com/alexatsu/jwt-in-go/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
)

var ctx = context.Background()

func register(c fiber.Ctx) error {
	conn, err := db.AcqureConnection(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	userQuery := user.New(conn.Conn())

	creds, err := utils.GetValidCredentials(c)
	if err != nil {
		return err
	}

	_, err = userQuery.GetUser(ctx, creds.Email)
	if err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "User already exists"})
	}

	hashedPass, err := utils.HashPassword(creds.Password)
	if err != nil {
		return err
	}

	userQuery.CreateUser(ctx, user.CreateUserParams{
		Email:    creds.Email,
		Password: pgtype.Text{String: hashedPass, Valid: true},
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User created"})
}

func login(c fiber.Ctx) error {
	creds, err := utils.GetValidCredentials(c)
	if err != nil {
		return err
	}

	conn, err := db.AcqureConnection(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	userQuery := user.New(conn)
	sessionQuery := session.New(conn)

	user, err := userQuery.GetUser(ctx, creds.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User not found"})
	}

	if !utils.CheckPasswordHash(creds.Password, user.Password.String) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID)
	if err != nil {
		log.Fatal(err, "some issues creating the token")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to login"})
	}

	err = sessionQuery.CreateSession(ctx, session.CreateSessionParams{
		UserID:       user.ID,
		RefreshToken: refreshToken,
	})
	if err != nil {
		log.Fatal(err, "failed to create session")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to login"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     utils.AccessTokenKey,
		Value:    accessToken,
		Expires:  utils.TimeOneMonth,
		HTTPOnly: true,
		Secure:   true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Login successful"})
}

func logout(c fiber.Ctx) error {
	conn, err := db.AcqureConnection(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	accessToken := c.Cookies(utils.AccessTokenKey)
	if accessToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Access token missing"})
	}

	claims, err := utils.ValidateToken(accessToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid access token"})
	}

	sessionQuery := session.New(conn)
	err = sessionQuery.DeleteSession(ctx, claims.UserId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to logout"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     utils.AccessTokenKey,
		Value:    "",
		Expires:  time.Unix(0, 0),
		HTTPOnly: true,
		Secure:   true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logout successful"})
}
