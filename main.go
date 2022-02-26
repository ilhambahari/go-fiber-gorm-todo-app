package main

import (
	// import packages
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"go-fiber-gorm-todo-app/database"
	"go-fiber-gorm-todo-app/todos"
)

func main() {
	app := fiber.New()
	initDatabase()
	setupV1(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	// sets up logger
	// Use middlewares for each route
	// This method will match all HTTP verbs: GET, POST, PUT etc Then create a log when every HTTP verb get invoked
	app.Use(logger.New(logger.Config{ // add Logger middleware with config
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	// Start server on port 2500
	err := app.Listen(":2500")
	if err != nil {
		panic(err)
	}
}

func initDatabase() {
	var err error
	dsn := "root:@tcp(127.0.0.1:3306)/todos-golang?charset=utf8mb4&parseTime=True&loc=Local"
	database.DBConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect database")
	}

	fmt.Println("Database successfully connected")

	database.DBConn.AutoMigrate(&todos.Todo{})
	fmt.Println("Database Migrated")
}

func setupV1(app *fiber.App) {
	v1 := app.Group("/v1")
	setupTodosRoutes(v1)
}

func setupTodosRoutes(grp fiber.Router) {
	todosRoutes := grp.Group("/todos")
	todosRoutes.Get("/", todos.GetAll)
	todosRoutes.Get("/:id", todos.GetOne)
	todosRoutes.Post("/", todos.AddTodo)
	todosRoutes.Delete("/:id", todos.DeleteTodo)
	todosRoutes.Put("/:id", todos.UpdateTodo)
}
