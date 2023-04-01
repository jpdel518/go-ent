package seed

import (
	"context"
	"github.com/jpdel518/go-ent/ent"
	"log"
	"time"
)

// CarSeed is a function to seed cars
func CarSeed(client *ent.Client) {
	var v []struct {
		Count int
	}
	err := client.Car.Query().
		Aggregate(ent.Count()).
		Scan(context.Background(), &v)
	if err != nil {
		log.Printf("failed seeding when count cars: %v", err)
		return
	} else {
		log.Printf("car count: %v and skip seeding", v)
	}

	if v[0].Count <= 0 {
		_, err := client.Car.Create().
			SetName("Toyota").
			SetModel("Prius").
			SetRegisteredAt(time.Now()).
			Save(context.Background())
		if err != nil {
			log.Printf("failed seeding when creating car: %v", err)
			return
		}

		_, err = client.Car.Create().
			SetName("Honda").
			SetModel("Civic").
			SetRegisteredAt(time.Now()).
			Save(context.Background())
		if err != nil {
			log.Printf("failed seeding when creating car: %v", err)
			return
		}

		_, err = client.Car.Create().
			SetName("Nissan").
			SetModel("Leaf").
			SetRegisteredAt(time.Now()).
			Save(context.Background())
		if err != nil {
			log.Printf("failed seeding when creating car: %v", err)
			return
		}
	}
}
