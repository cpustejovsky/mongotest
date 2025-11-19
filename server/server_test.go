package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cpustejovsky/mongotest/models"
	"github.com/cpustejovsky/mongotest/server"
)

type fakeStore struct {
	hit       bool
	expectIDs []string
}

func (fs *fakeStore) Create(animal models.Animal) (string, error) {
	return animal.ID, nil
}

func (fs *fakeStore) Fetch(id string) (*models.Animal, error) {
	if id != fs.expectIDs[0] {
		return nil, fmt.Errorf("wrong ID; got %v, expected %v", id, fs.expectIDs)
	}
	return &models.Animal{ID: fs.expectIDs[0]}, nil
}

func (fs *fakeStore) FetchAll() ([]models.Animal, error) {
	return []models.Animal{{ID: fs.expectIDs[0]}, {ID: fs.expectIDs[1]}}, nil
}

var expectID1 = "6373c7112476fec678ed0d3b"
var expectID2 = "6373c7112476fec678ed0d3b"
var fs = fakeStore{hit: false, expectIDs: []string{expectID1, expectID2}}

func TestGetAnimalByID(t *testing.T) {
	t.Run("returns number of snakes", func(t *testing.T) {
		expectId := "6373c7112476fec678ed0d3b"
		uri := "/animal/" + expectId
		req, err := http.NewRequest(http.MethodGet, uri, nil)
		if err != nil {
			t.Fatal(err)
		}
		res := httptest.NewRecorder()
		animalServer := server.New(&fs)
		animalServer.ServeHTTP(res, req)
		var animals models.Animal
		err = json.Unmarshal(res.Body.Bytes(), &animals)
		if err != nil {
			t.Fatal(err)
		}
		if got, expect := animals.ID, expectId; got != expect {
			t.Fatalf("got: %v; expect: %v\n", got, expect)
		}
	})
}

func TestGetAnimals(t *testing.T) {
	t.Run("returns number of snakes", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/animals", nil)
		res := httptest.NewRecorder()
		if err != nil {
			t.Fatal(err)
		}
		fs := fakeStore{hit: false, expectIDs: []string{expectID1, expectID2}}
		animalServer := server.New(&fs)
		animalServer.ServeHTTP(res, req)
		var animals []models.Animal
		err = json.Unmarshal(res.Body.Bytes(), &animals)
		if err != nil {
			t.Fatal(err)
		}
		if got, expect := len(fs.expectIDs), len(animals); got != expect {
			t.Fatalf("got: %v; expect: %v\n", got, expect)
		}
		for _, animal := range animals {
			if got := animal.ID; got != expectID1 && got != expectID2 {
				t.Fatalf("got %v; expected %v or %v\n", got, expectID1, expectID2)
			}
		}
	})
}
