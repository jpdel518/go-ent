package repository

import (
	"context"
	"github.com/jpdel518/go-ent/domain/model"
)

type UserRepository interface {
	Fetch(ctx context.Context, num int) (res []*model.User, err error)
	GetByID(ctx context.Context, id int) (*model.User, error)
	Create(ctx context.Context, u *model.User) (*model.User, error)
	Update(ctx context.Context, u *model.User) (*model.User, error)
	Delete(ctx context.Context, id int) error
}
