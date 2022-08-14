package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Route struct {
	Src string `form:"src"`
	Dst string `form:"dst"`
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

	// setup web server
	app := fiber.New()
	app.Static("/", "./view/public")
	app.Post("/", func(c *fiber.Ctx) error {
		r := new(Route)
		if err := c.BodyParser(r); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		collectionRoutes.InsertOne(ctx, r)
		return c.Status(201).SendString(fmt.Sprintf("{\n\tsrc: \"%s\",\n\tdst: \"%s\"\n}", r.Src, r.Dst))
	})
	app.Get("/*", func(c *fiber.Ctx) error {
		src := c.Params("*")
		if src == "" {
			return c.Status(400).SendString("bad route")
		}
		fmt.Println(src)
		var res Route
		collectionRoutes.FindOne(ctx, bson.D{{"src", src}}).Decode(res)
		fmt.Println(res)
		return c.SendString(res.Dst)
	})

	// start app
	app.Listen(":3000")
}
