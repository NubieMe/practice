package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	// Define a simple route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Start the server on port 3000
	if err := app.Listen(":3000"); err != nil {
		panic(err)
	}
}
