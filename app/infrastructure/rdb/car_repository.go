package rdb

import (
	"context"
	"github.com/jpdel518/go-ent/domain/model"
	"github.com/jpdel518/go-ent/domain/repository"
	"github.com/jpdel518/go-ent/ent"
	"github.com/jpdel518/go-ent/ent/car"
	"log"
)

type carRepository struct {
	client *ent.Client
}

func NewCarRepository(client *ent.Client) repository.CarRepository {
	return &carRepository{client: client}
}

func (r *carRepository) Fetch(ctx context.Context, num int) ([]*model.Car, error) {
	res := make([]*model.Car, 0)

	// fetch cars
	cars, err := r.client.Car.Query().Limit(num).All(ctx)
	if err != nil {
		log.Printf("failed fetching cars: %v", err)
		return res, err
	}

	// ent.Car -> model.Car
	for _, c := range cars {
		res = append(res, &model.Car{
			ID:           c.ID,
			Name:         c.Name,
			Model:        c.Model,
			RegisteredAt: c.RegisteredAt,
		})
	}
	return res, nil
}

func (r *carRepository) GetByID(ctx context.Context, id int) (*model.Car, error) {
	// get car
	c, err := r.client.Car.Get(ctx, id)
	if err != nil {
		log.Printf("failed getbyid car: %v", err)
		return nil, err
	}

	// ent.Car -> model.Car
	return &model.Car{
		ID:           c.ID,
		Name:         c.Name,
		Model:        c.Model,
		RegisteredAt: c.RegisteredAt,
	}, nil
}

func (r *carRepository) Create(ctx context.Context, u *model.Car) error {
	data, err := r.client.Car.Create().
		SetName(u.Name).
		SetModel(u.Model).
		SetRegisteredAt(u.RegisteredAt).
		Save(ctx)

	if err != nil {
		log.Printf("failed creating car: %v", err)
	}
	log.Printf("car was created: %v", data)

	return err
}

func (r *carRepository) Update(ctx context.Context, u *model.Car) error {
	data, err := r.client.Car.Update().
		Where(car.ID(u.ID)).
		SetName(u.Name).
		SetModel(u.Model).
		SetRegisteredAt(u.RegisteredAt).
		Save(ctx)

	if err != nil {
		log.Printf("failed updating car: %v", err)
	}
	log.Printf("car was updated: %v", data)

	return err
}

func (r *carRepository) Delete(ctx context.Context, id int) error {
	return r.client.Car.DeleteOneID(id).Exec(ctx)
}
