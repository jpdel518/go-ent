package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"time"
)

const dateLayout = "2006-01-02T15:04:05"

type Car struct {
	ID           int       `json:"id"`
	Name         string    `json:"name""`
	Model        string    `json:"model"`
	RegisteredAt time.Time `json:"registered_at"`
}

func (c Car) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required, validation.Length(1, 20)),
		validation.Field(&c.Model, validation.Required, validation.Length(1, 20)),
		validation.Field(&c.RegisteredAt, validation.Required, validation.Date(dateLayout)),
	)
}
