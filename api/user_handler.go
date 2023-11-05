package api

import (
	"errors"
	"net/http"

	"github.com/Narven/hotel-reservation/db"
	"github.com/Narven/hotel-reservation/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.userStore.DeleteUser(c.Context(), id); err != nil {
		return err
	}
	return c.JSON(map[string]string{"msg": "user deleted"})
}

func (h *UserHandler) HandleCreateUser(c *fiber.Ctx) error {
	var body models.CreateUserParams
	if err := c.BodyParser(&body); err != nil {
		return err
	}

	if err := body.Validate(c.Context()); err != nil {
		return err
	}

	user, err := models.NewUserFromParams(body)
	if err != nil {
		return err
	}

	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	var (
		userID = c.Params("id")
		params models.UpdateUserParams
	)
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.Validate(c.Context()); err != nil {
		return err
	}

	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	if err := h.userStore.UpdateUser(c.Context(), filter, params); err != nil {
		return err
	}
	user, err := h.userStore.GetUserByID(c.Context(), userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.Status(http.StatusNotFound)
			return c.JSON(map[string]string{"error": "not found"})
		}
		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	userID := c.Params("id")

	user, err := h.userStore.GetUserByID(c.Context(), userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.Status(http.StatusNotFound)
			return c.JSON(map[string]string{"error": "not found"})
		}
		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}
