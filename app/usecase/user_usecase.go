package usecase

import (
	"context"
	"github.com/jpdel518/go-ent/domain/model"
	"github.com/jpdel518/go-ent/domain/repository"
	"log"
	"mime/multipart"
	"sync"
	"time"
)

type UserUsecase interface {
	Fetch(ctx context.Context, num int) ([]*model.User, error)
	GetByID(ctx context.Context, id int) (*model.User, error)
	Create(ctx context.Context, u *model.User, f multipart.File, fh *multipart.FileHeader) error
	Update(ctx context.Context, u *model.User, f multipart.File, fh *multipart.FileHeader) error
	Delete(ctx context.Context, id int) error
}

type userUsecase struct {
	userRepo       repository.UserRepository
	carRepo        repository.CarRepository
	userFileRepo   repository.UserFileRepository
	contextTimeout time.Duration
}

// NewUserUsecase will create new an userUsecase object
func NewUserUsecase(u repository.UserRepository, c repository.CarRepository, f repository.UserFileRepository, timeout time.Duration) UserUsecase {
	return &userUsecase{
		userRepo:       u,
		carRepo:        c,
		userFileRepo:   f,
		contextTimeout: timeout,
	}
}

// fillCarDetails will fill up car details that is concerned with a parameter of user object
func (usecase *userUsecase) getCarDetails(c context.Context, user *model.User) ([]*model.Car, error) {
	// TODO errgroupを使ってエラーをハンドリング + contextで後続の処理をキャンセルするようにした方がいいかも
	var wg sync.WaitGroup
	// Get the car's id
	var cars []*model.Car

	// Using goroutine to fetch the car's detail
	chanCar := make(chan *model.Car)

	for carID := range user.Cars {
		wg.Add(1)
		go func(carID int) {
			defer wg.Done()
			res, err := usecase.carRepo.GetByID(c, carID)
			chanCar <- res
			if err != nil {
				// TODO エラーをハンドリング
				log.Printf("fillCarDetails get car ID: %d error: %v", carID, err)
			}
		}(carID)
	}

	// release channel
	go func() {
		wg.Wait()
		close(chanCar)
	}()

	// retrieve
	for car := range chanCar {
		if car != nil {
			cars = append(cars, car)
		}
	}
	wg.Wait()

	return cars, nil
}

// Fetch will retrieve user
func (usecase *userUsecase) Fetch(c context.Context, num int) ([]*model.User, error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, usecase.contextTimeout)
	defer cancel()

	res, err := usecase.userRepo.Fetch(ctx, num)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetByID will find a user by id
func (usecase *userUsecase) GetByID(c context.Context, id int) (*model.User, error) {
	ctx, cancel := context.WithTimeout(c, usecase.contextTimeout)
	defer cancel()

	res, err := usecase.userRepo.GetByID(ctx, id)
	if err != nil {
		return &model.User{}, err
	}

	// cars, err := res.QueryCars().All(ctx)
	// if err != nil {
	// 	log.Printf("failed querying user cars: %w", err)
	// }
	// // Query the inverse edge.
	// for _, c := range cars {
	// 	owner, err := c.QueryOwner().Only(ctx)
	// 	if err != nil {
	// 		log.Printf("failed querying car %q owner: %w", c.Model, err)
	// 	}
	// 	log.Printf("car %q owner: %q\n", c.Model, owner.Email)
	// }

	return res, nil
}

// Create will register a user
func (usecase *userUsecase) Create(c context.Context, u *model.User, f multipart.File, fh *multipart.FileHeader) error {
	ctx, cancel := context.WithTimeout(c, usecase.contextTimeout)
	defer cancel()

	// create user
	user, err := usecase.userRepo.Create(ctx, u)
	if err != nil {
		return err
	}

	// upload user avatar file
	if f != nil && fh != nil {
		filename, err := usecase.userFileRepo.Create(ctx, user.ID, f, fh)
		if err != nil {
			return err
		}
		user.Avatar = filename

		// update user
		_, err = usecase.userRepo.Update(ctx, user)
	}
	return err
}

// Update will update a user
func (usecase *userUsecase) Update(c context.Context, u *model.User, f multipart.File, fh *multipart.FileHeader) error {
	ctx, cancel := context.WithTimeout(c, usecase.contextTimeout)
	defer cancel()

	if f != nil && fh != nil {
		filename, err := usecase.userFileRepo.Update(ctx, u.ID, f, fh)
		if err != nil {
			return err
		}
		u.Avatar = filename
	}

	_, err := usecase.userRepo.Update(ctx, u)
	return err
}

// Delete will delete a user by id
func (usecase *userUsecase) Delete(c context.Context, id int) error {
	ctx, cancel := context.WithTimeout(c, usecase.contextTimeout)
	defer cancel()

	return usecase.userRepo.Delete(ctx, id)
}
