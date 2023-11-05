package models

import (
	"context"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 7
)

type UpdateUserParams struct {
	FirstName string `bson:"first_name" json:"first_name" validate:"omitempty,min=2,max=100"`
	LastName  string `bson:"last_name" json:"last_name" validate:"omitempty,min=2,max=100"`
}

type CreateUserParams struct {
	FirstName string `json:"first_name" validate:"required,min=2,max=100"`
	LastName  string `json:"last_name" validate:"required,min=2,max=100"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=7"`
}

func (params UpdateUserParams) Validate(ctx context.Context) error {
	validate := validator.New()
	if err := validate.StructCtx(ctx, params); err != nil {
		return err
	}
	return nil
}

func (params CreateUserParams) Validate(ctx context.Context) error {
	validate := validator.New()
	if err := validate.StructCtx(ctx, params); err != nil {
		return err
	}
	return nil
}

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	FirstName string             `bson:"first_name" json:"first_name"`
	LastName  string             `bson:"last_name" json:"last_name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Password:  string(encryptedPassword),
	}, nil
}
