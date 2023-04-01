package rdb

import (
	"context"
	"github.com/jpdel518/go-ent/domain/model"
	"github.com/jpdel518/go-ent/domain/repository"
	"github.com/jpdel518/go-ent/ent"
	"github.com/jpdel518/go-ent/ent/user"
	"log"
)

type userRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) repository.UserRepository {
	return &userRepository{client: client}
}

func (r *userRepository) Fetch(ctx context.Context, num int) ([]*model.User, error) {
	res := make([]*model.User, 0)

	// fetch users
	users, err := r.client.User.Query().Limit(num).All(ctx)
	if err != nil {
		log.Printf("failed fetching users: %v", err)
		return res, err
	}

	// ent.User -> model.User
	for _, u := range users {
		cars := make([]model.Car, 0)
		for _, c := range u.Edges.Cars {
			cars = append(cars, model.Car{
				ID:           c.ID,
				Name:         c.Name,
				Model:        c.Model,
				RegisteredAt: c.RegisteredAt,
			})
		}
		res = append(res, &model.User{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Age:       u.Age,
			Cars:      cars,
		})
	}
	return res, nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	// get user
	u, err := r.client.User.Get(ctx, id)
	if err != nil {
		log.Printf("failed getbyid user: %v", err)
		return nil, err
	}

	// ent.User -> model.User
	cars := make([]model.Car, 0)
	for _, c := range u.Edges.Cars {
		cars = append(cars, model.Car{
			ID:           c.ID,
			Name:         c.Name,
			Model:        c.Model,
			RegisteredAt: c.RegisteredAt,
		})
	}
	return &model.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Age:       u.Age,
		Cars:      cars,
	}, nil
}

func (r *userRepository) Create(ctx context.Context, u *model.User) (*model.User, error) {
	data, err := r.client.User.Create().
		SetFirstName(u.FirstName).
		SetLastName(u.LastName).
		SetEmail(u.Email).
		SetAge(u.Age).
		AddCarIDs(u.CarIDs...).
		Save(ctx)

	if err != nil {
		log.Printf("failed creating user: %v", err)
		return nil, err
	}
	log.Printf("user was created: %v", data)

	// ent.User -> model.User
	cars := make([]model.Car, 0)
	for _, c := range data.Edges.Cars {
		cars = append(cars, model.Car{
			ID:           c.ID,
			Name:         c.Name,
			Model:        c.Model,
			RegisteredAt: c.RegisteredAt,
		})
	}
	res := &model.User{
		ID:        data.ID,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     data.Email,
		Age:       data.Age,
		Cars:      cars,
	}

	return res, err
}

func (r *userRepository) Update(ctx context.Context, u *model.User) (*model.User, error) {
	data, err := r.client.User.Update().
		Where(user.ID(u.ID)).
		SetFirstName(u.FirstName).
		SetLastName(u.LastName).
		SetEmail(u.Email).
		SetAge(u.Age).
		ClearCars().
		AddCarIDs(u.CarIDs...).
		Save(ctx)

	if err != nil {
		log.Printf("failed updating user: %v", err)
		return nil, err
	}
	log.Printf("user was updated: %v", data)

	return u, err
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	return r.client.User.DeleteOneID(id).Exec(ctx)
}
