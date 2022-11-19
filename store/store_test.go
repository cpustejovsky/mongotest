package store_test

import (
	"context"
	"github.com/cpustejovsky/mongotest/models"
	"github.com/cpustejovsky/mongotest/store"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestAnimalStore(t *testing.T) {
	clientOptions := options.Client().
		ApplyURI("mongodb://localhost:27017/mongotest")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)
	var id1 string
	var id2 string
	animalSpecies := "snake"
	var animalStore *store.AnimalStore

	t.Run("Animal Store New", func(t *testing.T) {
		animalStore = store.NewAnimalStore(client, "animalstest")
		_, err = animalStore.Collection.DeleteMany(ctx, bson.D{})
		assert.Nil(t, err)
	})
	t.Run("Animal Store Create", func(t *testing.T) {
		newAnimal := models.Animal{
			Species: animalSpecies,
		}
		id1, err = animalStore.Create(newAnimal)
		assert.Nil(t, err)
		assert.NotEmpty(t, id1)
	})
	t.Run("Animal Store Fetch", func(t *testing.T) {
		animal, err := animalStore.Fetch(id1)
		assert.Nil(t, err)
		assert.Equal(t, animalSpecies, animal.Species)
	})
	t.Run("Animal Store FetchAll", func(t *testing.T) {
		newAnimal := models.Animal{
			Species: animalSpecies,
		}
		id2, err = animalStore.Create(newAnimal)
		animals, err := animalStore.FetchAll()
		assert.Nil(t, err)
		for _, animal := range animals {
			if !(animal.ID == id1 || animal.ID == id2) {
				t.Fatalf("Wanted %v to be %v or %v", animal.ID, id1, id2)
			}
		}
	})
}
