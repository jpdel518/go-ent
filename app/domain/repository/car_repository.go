package repository

import (
	"context"
	"github.com/jpdel518/go-ent/domain/model"
)

type CarRepository interface {
	Fetch(ctx context.Context, num int) (res []*model.Car, err error)
	GetByID(ctx context.Context, id int) (*model.Car, error)
	Create(ctx context.Context, u *model.Car) error
	Update(ctx context.Context, u *model.Car) error
	Delete(ctx context.Context, id int) error
}
