package main

import (
	"context"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Route struct {
	Url   string `json:"url"`
	Route string `json:"route"`
}

func main() {
	godotenv.Load()

	// load and check env variables
	mongoUri := os.Getenv("MONGO_URI")
	mongoDb := os.Getenv("MONGO_DB")
	if !strings.HasPrefix(mongoUri, "mongodb://") {
		panic("Invalid MONGO_URI env var")
	}
	if mongoDb == "" {
		panic("MONGO_DB env var have to defined")
	}

	// setup mongo
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	db := client.Database(mongoDb)
	collectionRoutes := db.Collection("routes")

	// setup web server and its handlers
	app := fiber.New()
	app.Static("/", "./view/public")
	app.Post("/", func(c *fiber.Ctx) error {
		r := new(Route)
		if err := c.BodyParser(r); err != nil {
			return c.SendStatus(500)
		}
		if r.Route == "" || r.Url == "" {
			return c.SendStatus(400)
		}
		if _, err := collectionRoutes.InsertOne(ctx, r); err != nil {
			return c.SendStatus(409)
		}
		return c.Status(201).JSON(r)
	})
	app.Get("/*", func(c *fiber.Ctx) error {
		route := c.Params("*")
		if route == "" {
			return c.Status(400).SendString("bad route")
		}
		var res Route
		if err := collectionRoutes.FindOne(ctx, bson.D{{"route", route}}).Decode(&res); err != nil {
			return c.Redirect("/")
		}
		return c.Redirect(res.Url)
	})

	// start app
	app.Listen(":3000")
}
