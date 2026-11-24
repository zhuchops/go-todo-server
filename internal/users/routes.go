package users

import (
	"zhuchok/go-todo-server/ent"

	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(router fiber.Router, db *ent.Client) {
	router.Get(
		"/users/me",
		func(c fiber.Ctx) error {
			return getMe(c)
		},
	)
}
