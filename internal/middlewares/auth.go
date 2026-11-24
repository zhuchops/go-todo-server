package middlewares

import (
	"fmt"
	"strings"
	"zhuchok/go-todo-server/ent"
	"zhuchok/go-todo-server/internal/share"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v4"
)

func AuthRequired(c fiber.Ctx, db *ent.Client) error {
	h := c.Get("Authorization")
	if !strings.HasPrefix(h, "Bearer ") {
		return c.
			Status(fiber.StatusUnauthorized).
			JSON(share.MyError{
				Message: "you must provide authentication token",
				Error:   "",
			})
	}
	claims := &share.Claims{}
	token, err := jwt.ParseWithClaims(
		strings.TrimPrefix(h, "Bearer "),
		claims,
		func(t *jwt.Token) (any, error) { return share.JwtSecret, nil },
		jwt.WithValidMethods([]string{"HS256"}),
	)
	if err != nil || !token.Valid {
		return c.
			Status(fiber.StatusUnauthorized).
			JSON(share.MyError{
				Message: "failed to parse token",
				Error:   fmt.Sprint(err),
			})
	}
	user, err := db.User.Get(c, claims.UserId)
	if err != nil {
		return c.
			Status(fiber.StatusUnauthorized).
			JSON(share.MyError{
				Message: "failed to find user",
				Error:   fmt.Sprint(err),
			})
	}
	fiber.Locals(c, "inner_user", user)
	return c.Next()
}
