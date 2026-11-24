package auth

import (
	"zhuchok/go-todo-server/ent"

	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(router fiber.Router, db *ent.Client) {
	router.Post(
		"/login",
		func(c fiber.Ctx) error {
			return loginUser(c, db)
		},
	)
	router.Post(
		"/register",
		func(c fiber.Ctx) error {
			return registerUser(c, db)
		},
	)
}
