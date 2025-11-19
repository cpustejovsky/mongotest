package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/cpustejovsky/mongotest/models"
	"github.com/cpustejovsky/mongotest/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestAnimalStore(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test")
	}
	clientOptions := options.Client().
		ApplyURI("mongodb://localhost:27017/mongotest")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(ctx)
	var id1 string
	var id2 string
	animalSpecies := "snake"
	var animalStore *store.AnimalStore

	t.Run("Animal Store New", func(t *testing.T) {
		animalStore = store.NewAnimalStore(client, "animalstest")
		_, err = animalStore.Collection.DeleteMany(ctx, bson.D{})
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("Animal Store Create", func(t *testing.T) {
		newAnimal := models.Animal{
			Species: animalSpecies,
		}
		id1, err = animalStore.Create(newAnimal)
		if err != nil {
			t.Fatal(err)
		}
		if id1 != "" {
			t.Fatal("id was empty")
		}

	})
	t.Run("Animal Store Fetch", func(t *testing.T) {
		animal, err := animalStore.Fetch(id1)
		if err != nil {
			t.Fatal(err)
		}
		if got, expect := animalSpecies, animal.Species; got != expect {
			t.Fatalf("got %v, expected %v\n", got, expect)
		}
	})
	t.Run("Animal Store FetchAll", func(t *testing.T) {
		newAnimal := models.Animal{
			Species: animalSpecies,
		}
		id2, err = animalStore.Create(newAnimal)
		animals, err := animalStore.FetchAll()
		if err != nil {
			t.Fatal(err)
		}
		for _, animal := range animals {
			if !(animal.ID == id1 || animal.ID == id2) {
				t.Fatalf("Wanted %v to be %v or %v", animal.ID, id1, id2)
			}
		}
	})
}
