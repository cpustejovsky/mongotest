package server_test

import (
	"encoding/json"
	"fmt"
	"github.com/cpustejovsky/mongotest/models"
	"github.com/cpustejovsky/mongotest/server"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeStore struct {
	hit     bool
	wantIDs []string
}

func (fs *fakeStore) Create(animal models.Animal) (string, error) {
	return animal.ID, nil
}

func (fs *fakeStore) Fetch(id string) (*models.Animal, error) {
	if id != fs.wantIDs[0] {
		return nil, fmt.Errorf("wrong ID; got %v, wanted %v", id, fs.wantIDs)
	}
	return &models.Animal{ID: fs.wantIDs[0]}, nil
}

func (fs *fakeStore) FetchAll() ([]models.Animal, error) {
	return []models.Animal{{ID: fs.wantIDs[0]}, {ID: fs.wantIDs[1]}}, nil
}

var wantID1 = "6373c7112476fec678ed0d3b"
var wantID2 = "6373c7112476fec678ed0d3b"
var fs = fakeStore{hit: false, wantIDs: []string{wantID1, wantID2}}

func TestGetAnimalByID(t *testing.T) {
	t.Run("returns number of snakes", func(t *testing.T) {
		wantId := "6373c7112476fec678ed0d3b"
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/animal/%v", wantId), nil)
		//req, err := http.NewRequest(http.MethodGet, "/animal/:id", nil)
		res := httptest.NewRecorder()
		assert.Nil(t, err)
		animalServer := server.New(&fs)
		animalServer.ServeHTTP(res, req)
		var a models.Animal
		err = json.Unmarshal(res.Body.Bytes(), &a)
		assert.Nil(t, err)
		assert.Equal(t, wantId, a.ID)
	})
}

func TestGetAnimals(t *testing.T) {
	t.Run("returns number of snakes", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/animals", nil)
		res := httptest.NewRecorder()
		assert.Nil(t, err)
		fs := fakeStore{hit: false, wantIDs: []string{wantID1, wantID2}}
		animalServer := server.New(&fs)
		animalServer.ServeHTTP(res, req)
		var a []models.Animal
		err = json.Unmarshal(res.Body.Bytes(), &a)
		assert.Nil(t, err)
		assert.Equal(t, len(fs.wantIDs), len(a))
		for _, animal := range a {
			assert.True(t, assert.Equal(t, animal.ID, wantID1) || assert.Equal(t, animal.ID, wantID2))
		}
	})
}
