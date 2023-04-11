package server

import (
	"encoding/json"
	"fmt"
	"github.com/cpustejovsky/mongotest/models"
	"github.com/cpustejovsky/mongotest/store"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type AnimalServer struct {
	store store.Repository
}

func New(store store.Repository) *AnimalServer {
	ss := AnimalServer{
		store: store,
	}
	return &ss
}

func (s *AnimalServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := httprouter.New()
	router.GET("/", s.pingHandler)
	router.GET("/animals", s.animalsHandler)
	router.GET("/animal/:id", s.animalHandler)
	router.POST("/animal/new", s.newAnimalHandler)
	router.ServeHTTP(w, r)
}

func (s *AnimalServer) pingHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "Howdy Ho!")
}

func (s *AnimalServer) animalHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	s.getAnimal(w, id)
}

func (s *AnimalServer) animalsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	s.getAnimals(w)
}

func (s *AnimalServer) newAnimalHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var af models.AnimalForm
	err := json.NewDecoder(r.Body).Decode(&af)
	if err != nil {
		fmt.Fprint(w, fmt.Sprintf("Error encountered: %v", err))
	}
	if af.Species == nil {
		fmt.Fprint(w, fmt.Sprintf("Error encountered: %v", err))
	}
	s.createAnimal(w, *af.Species)
}

func (s *AnimalServer) createAnimal(w http.ResponseWriter, species string) {
	var animal = models.Animal{
		Species: species,
	}
	fetchedAnimal, err := s.store.Create(animal)
	if err != nil {
		fmt.Fprint(w, fmt.Sprintf("Error encountered: %v", err))
	}
	err = json.NewEncoder(w).Encode(fetchedAnimal)
	if err != nil {
		fmt.Fprint(w, fmt.Sprintf("Error encountered: %v", err))
	}
}

func (s *AnimalServer) getAnimal(w http.ResponseWriter, id string) {
	fetchedAnimal, err := s.store.Fetch(id)
	if err != nil {
		fmt.Fprint(w, fmt.Sprintf("Error encountered: %v", err))
	}
	err = json.NewEncoder(w).Encode(fetchedAnimal)
	if err != nil {
		fmt.Fprint(w, fmt.Sprintf("Error encountered: %v", err))
	}
}

func (s *AnimalServer) getAnimals(w http.ResponseWriter) {
	fetchedAnimals, err := s.store.FetchAll()
	if err != nil {
		fmt.Fprint(w, fmt.Sprintf("Error encountered: %v", err))
	}
	err = json.NewEncoder(w).Encode(fetchedAnimals)
	if err != nil {
		fmt.Fprint(w, fmt.Sprintf("Error encountered: %v", err))
	}
}
