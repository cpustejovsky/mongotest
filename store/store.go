package store

import (
	"context"
	"errors"
	"github.com/cpustejovsky/mongotest/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type AnimalStore struct {
	Collection *mongo.Collection
}

func NewAnimalStore(client *mongo.Client, dbName string) *AnimalStore {
	database := client.Database(dbName)
	collection := database.Collection("animals")
	return &AnimalStore{
		Collection: collection,
	}
}

type Repository interface {
	Create(animal models.Animal) (string, error)
	Fetch(id string) (*models.Animal, error)
	FetchAll() ([]models.Animal, error)
}

func (a *AnimalStore) Create(animal models.Animal) (string, error) {
	insertResult, err := a.Collection.InsertOne(context.TODO(), bson.D{
		{"species", animal.Species},
	})
	if err != nil {
		return "", err
	}
	id, ok := insertResult.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("id did not coerce to primitive.ObjectID")
	}
	return id.Hex(), nil
}

func (a *AnimalStore) Fetch(id string) (*models.Animal, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &models.Animal{}, err
	}
	var animal models.Animal
	err = a.Collection.FindOne(context.TODO(), bson.M{
		"_id": oid,
	}).Decode(&animal)
	if err != nil {
		return &models.Animal{}, err
	}
	return &animal, nil
}

func (a *AnimalStore) FetchAll() ([]models.Animal, error) {
	ctx := context.TODO()
	var animals []models.Animal
	cursor, err := a.Collection.Find(ctx, bson.D{})
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var animal models.Animal
		if err := cursor.Decode(&animal); err != nil {
			return nil, err
		}
		log.Println(animal.ID)
		animals = append(animals, animal)
	}
	if err != nil {
		return nil, err
	}
	return animals, nil
}
