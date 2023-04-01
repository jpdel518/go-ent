package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	CarIDs    []int  `json:"car_ids"`
	Cars      []Car  `json:"cars"`
	Avatar    string `json:"avatar"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.FirstName, validation.Required, validation.Length(1, 20)),
		validation.Field(&u.LastName, validation.Required, validation.Length(1, 20)),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Age, validation.Min(0)),
	)
}
