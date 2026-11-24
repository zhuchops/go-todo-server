package users

import (
	"zhuchok/go-todo-server/ent"

	"github.com/gofiber/fiber/v3"
)

type GetMeOutput struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

func getMe(c fiber.Ctx) error {
	user := fiber.Locals[*ent.User](c, "inner_user")
	return c.JSON(GetMeOutput{
		Email:    user.Email,
		Nickname: user.Nickname,
	})
}
