package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Route struct {
	Src string `form:"src"`
	Dst string `form:"dst"`
}

func main() {
	godotenv.Load()
	app := fiber.New()
	app.Static("/", "./view/public")
	app.Post("/", func(c *fiber.Ctx) error {
		r := new(Route)
		if err := c.BodyParser(r); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		return c.SendString(fmt.Sprintf("{\n\tsrc: \"%s\",\n\tdst: \"%s\"\n}", r.Src, r.Dst))
	})
	app.Get("/*", func(c *fiber.Ctx) error {
		route := c.Params("*")
		if route != "" {
			return c.SendString(route)
		}
		return c.Status(400).SendString("bad route")
	})
	app.Listen(":3000")
}
