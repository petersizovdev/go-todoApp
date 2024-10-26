package main

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	app := fiber.New()

	todos := []Todo{}

	// Get
	app.Get("/api/todos/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})

	// Post
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}
		if err := c.BodyParser(todo); err != nil {
			return err
		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"err": "Todo body is required"})
		}

		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(201).JSON(todo)
	})

	// Patch
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"err": "Invalid ID"})
		}

		for i, td := range todos {
			if td.ID == id {
				// Парсим тело запроса и обновляем поля
				updatedTodo := &Todo{}
				if err := c.BodyParser(updatedTodo); err != nil {
					return err
				}

				// Обновляем только те поля, которые пришли в запросе
				if updatedTodo.Body != "" {
					todos[i].Body = updatedTodo.Body
				}
				if updatedTodo.Completed {
					todos[i].Completed = updatedTodo.Completed
				}

				return c.Status(200).JSON(todos[i])
			}
		}

		return c.Status(404).JSON(fiber.Map{"err": "Todo not found"})
	})

	// Delete

	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"err": "Invalid ID"})
		}

		for i, td := range todos {
			if td.ID == id {
				todos = append (todos[:i], todos[i+1:]...)	
				return c.Status(200).JSON(fiber.Map{"success": true})
		
			}
		}
		return c.Status(404).JSON(fiber.Map{"err": "Todo not found"})
	})

	log.Fatal(app.Listen(":4000"))
}
