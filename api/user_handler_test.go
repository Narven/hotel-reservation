package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/Narven/hotel-reservation/db"
	"github.com/Narven/hotel-reservation/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testDbURI = "mongodb://localhost:27017"
	dbName    = "hotel_dev"
)

type testDB struct {
	db.UserStore
}

func setup(t *testing.T) *testDB {
	fmt.Println("->> test setup")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testDbURI))
	if err != nil {
		log.Fatal(err)
	}

	return &testDB{
		UserStore: db.NewMongoUsersStore(client, dbName),
	}
}

func (tdb *testDB) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		log.Fatal(err)
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandleCreateUser)

	payload := models.CreateUserParams{
		FirstName: "Some",
		LastName:  "Foo",
		Email:     "some@foo.com",
		Password:  "01234567",
	}

	b, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var user models.User
	json.NewDecoder(resp.Body).Decode(&user)

	if user.FirstName != payload.FirstName {
		t.Errorf("expected %s, got %s", payload.FirstName, user.FirstName)
	}

	if user.LastName != payload.LastName {
		t.Errorf("expected %s, got %s", payload.LastName, user.LastName)
	}

	if user.Email != payload.Email {
		t.Errorf("expected %s, got %s", payload.Email, user.Email)
	}
}
