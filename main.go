package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/Narven/hotel-reservation/api"
	"github.com/Narven/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbURI   = "mongodb://localhost:27017"
	dbName  = "hotel"
	userCol = "users"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	ctx := context.Background()
	port := flag.String("port", "5000", "The listen address of the API server")
	dbName := flag.String("db", "hotel", "The database name")
	flag.Parse()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatal(err)
	}

	userHandler := api.NewUserHandler(db.NewMongoUsersStore(client, *dbName))

	app := fiber.New(config)
	apiV1 := app.Group("/api/v1")
	apiV1.Get("/users", userHandler.HandleGetUsers)
	apiV1.Post("/users", userHandler.HandleCreateUser)
	apiV1.Get("/users/:id", userHandler.HandleGetUser)
	apiV1.Delete("/users/:id", userHandler.HandleDeleteUser)
	apiV1.Put("/users/:id", userHandler.HandleUpdateUser)

	err = app.Listen(fmt.Sprintf(":%s", *port))
	log.Fatal(err)
}
