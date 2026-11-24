package auth

import (
	"fmt"
	"time"
	"zhuchok/go-todo-server/ent"
	"zhuchok/go-todo-server/ent/user"
	"zhuchok/go-todo-server/internal/share"
	simplecrypto "zhuchok/go-todo-server/internal/simple_crypto"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v4"
)

type LoginUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserOutput struct {
	AccessToken string `json:"access_token"`
}

func loginUser(c fiber.Ctx, db *ent.Client) error {
	var input LoginUserInput
	if err := c.Bind().JSON(&input); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(share.MyError{
				Message: "wrong json input",
				Error:   fmt.Sprint(err),
			})
	}
	users, err := db.User.Query().
		Where(
			user.EmailEQ(input.Email),
		).
		All(c)
	if err != nil {
		return c.
			Status(fiber.StatusInternalServerError).
			JSON(share.MyError{
				Message: "failed to find user",
				Error:   fmt.Sprint(err),
			})
	}
	if len(users) == 0 {
		return c.
			Status(fiber.StatusUnauthorized).
			JSON(share.MyError{
				Message: "wrong email or password",
				Error:   fmt.Sprint(err),
			})
	}
	user := users[0]
	if !simplecrypto.VerifyArgon2(user.PasswordHash, input.Password) {
		return c.
			Status(fiber.StatusUnauthorized).
			JSON(share.MyError{
				Message: "wrong email or password",
				Error:   fmt.Sprint(err),
			})
	}
	claims := share.Claims{
		UserId: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(share.JwtSecret)
	if err != nil {
		return c.
			Status(fiber.StatusInternalServerError).
			JSON(share.MyError{
				Message: "failed to sign token",
				Error:   fmt.Sprint(err),
			})
	}
	return c.JSON(LoginUserOutput{AccessToken: token})
}

type RegisterUserInput struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserOutput struct {
	Id       string `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

func registerUser(c fiber.Ctx, db *ent.Client) error {
	var input RegisterUserInput
	if err := c.Bind().JSON(&input); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(share.MyError{
				Message: "bad json",
				Error:   fmt.Sprint(err),
			})
	}
	passwordHash, err := simplecrypto.HashArgon2(input.Password)
	if err != nil {
		return c.
			Status(fiber.StatusInternalServerError).
			JSON(share.MyError{
				Message: "failed to hash password",
				Error:   fmt.Sprint(err),
			})
	}
	_, err = db.User.
		Create().
		SetEmail(input.Email).
		SetNickname(input.Nickname).
		SetPasswordHash(passwordHash).
		Save(c.Context())
	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(share.MyError{
				Message: "such user already exists or internal error",
				Error:   fmt.Sprint(err),
			})
	}
	return c.SendStatus(fiber.StatusCreated)
}
