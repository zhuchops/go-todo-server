package tasks

import (
	"zhuchok/go-todo-server/ent"

	"github.com/gofiber/fiber/v3"
)

func RegisterRoutes(router fiber.Router, db *ent.Client) {
	router.Get(
		"/tasks",
		func(c fiber.Ctx) error {
			return listTasks(c, db)
		},
	)
	router.Post(
		"/tasks",
		func(c fiber.Ctx) error {
			return createTask(c, db)
		},
	)
	router.Patch(
		"/tasks/status",
		func(c fiber.Ctx) error {
			return toggleTask(c, db)
		},
	)
}
