package tasks

import (
	"fmt"
	"strconv"
	"zhuchok/go-todo-server/ent"
	"zhuchok/go-todo-server/ent/task"
	"zhuchok/go-todo-server/internal/share"

	"github.com/gofiber/fiber/v3"
	"github.com/samber/lo"
)

type ListTasksQuery struct {
	Status string `query:"status"`
}

type ListTasksItem struct {
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type ListTasksOutput struct {
	Tasks []ListTasksItem `json:"tasks"`
}

func listTasks(c fiber.Ctx, db *ent.Client) error {
	var query ListTasksQuery
	if err := c.Bind().Query(&query); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(share.MyError{
				Message: "empty required query: status",
				Error:   fmt.Sprint(err),
			})
	}
	user := fiber.Locals[*ent.User](c, "inner_user")

	var status task.Status
	switch query.Status {
	case "todo":
		status = task.StatusTodo
	case "done":
		status = task.StatusTodo
	default:
		return c.
			Status(fiber.StatusBadRequest).
			JSON(share.MyError{
				Message: "status must be <todo> or <done>",
				Error:   "",
			})
	}

	tasks, err := db.User.QueryTasks(user).
		Where(task.StatusEQ(status)).
		All(c)
	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(share.MyError{
				Message: "cannot find tasks from this user",
				Error:   fmt.Sprint(err),
			})
	}
	tasksOutput := lo.Map(tasks, func(task *ent.Task, index int) ListTasksItem {
		return ListTasksItem{Content: task.Content, CreatedAt: task.CreatedAt.GoString()}
	})
	return c.JSON(tasksOutput)
}

type CreateTaskInput struct {
	Content string `json:"content"`
	Status  string `json:"status"`
}

func createTask(c fiber.Ctx, db *ent.Client) error {
	var input CreateTaskInput
	if err := c.Bind().JSON(&input); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(share.MyError{
				Message: "failed to parse json",
				Error:   fmt.Sprint(err),
			})
	}

	var status task.Status
	switch input.Status {
	case "todo":
		status = task.StatusTodo
	case "done":
		status = task.StatusDone
	default:
		return c.
			Status(fiber.StatusBadRequest).
			JSON(share.MyError{
				Message: "status must be <todo> or <done>",
				Error:   "",
			})
	}
	owner := fiber.Locals[*ent.User](c, "inner_user")

	_, err := db.Task.
		Create().
		SetOwner(owner).
		SetStatus(status).
		SetContent(input.Content).
		Save(c.Context())
	if err != nil {
		return c.
			Status(fiber.StatusInternalServerError).
			JSON(share.MyError{
				Message: "failed to create task",
				Error:   fmt.Sprint(err),
			})
	}
	return c.SendStatus(fiber.StatusCreated)
}

type ToggleTaskInput struct {
	TaskId    string `json:"task_id"`
	NewStatus string `json:"status"`
}

func toggleTask(c fiber.Ctx, db *ent.Client) error {
	var input ToggleTaskInput
	if err := c.Bind().JSON(&input); err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(share.MyError{
				Message: "failed to parse json",
				Error:   fmt.Sprint(err),
			})
	}

	var status task.Status
	switch input.NewStatus {
	case "todo":
		status = task.StatusTodo
	case "done":
		status = task.StatusDone
	default:
		return c.
			Status(fiber.StatusBadRequest).
			JSON(share.MyError{
				Message: "status must be <todo> or <done>",
				Error:   "",
			})
	}
	user := fiber.Locals[*ent.User](c, "inner_user")
	userId := user.ID

	taskId, err := strconv.Atoi(input.TaskId)
	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(share.MyError{
				Message: "id must be an int",
				Error:   "",
			})
	}
	err = db.Task.Update().
		Where(
			task.OwnerIDEQ(userId),
			task.IDEQ(taskId),
		).
		SetStatus(status).
		Exec(c)
	if err != nil {
		return c.
			Status(fiber.StatusBadRequest).
			JSON(share.MyError{
				Message: "failed to update task, no such task belongs to such user",
				Error:   fmt.Sprint(err),
			})
	}
	return c.SendStatus(fiber.StatusOK)
}
