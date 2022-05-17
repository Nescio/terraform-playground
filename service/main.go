//Petstore is a simple HTTP API that provides a RESTful set of services for creating, updating, deleting, and retrieving pets.
package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

var pets = map[int]*pet{}

func main() {
	http.HandleFunc("/api/v1/pets", handlePets)
	http.ListenAndServe(":8080", nil)
}

type pet struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Species string `json:"species"`
	Age     int    `json:"age"`
}

func handlePets(w http.ResponseWriter, r *http.Request) {
	log.Printf(r.Method + " " + r.URL.Path + " " + r.Proto + " " + r.UserAgent() + " " + r.RemoteAddr)
	switch r.Method {
	case http.MethodPost:
		handlePost(w, r)
	case http.MethodGet:
		handleGetById(w, r)
	case http.MethodPatch:
		handlePatch(w, r)
	case http.MethodDelete:
		handleDelete(w, r)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	pet := pet{}
	if err := json.NewDecoder(r.Body).Decode(&pet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pet.ID = rand.Int()
	pets[pet.ID] = &pet
	log.Printf("Created pet %v", pet)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pet)
}

func handleGetById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pet, ok := pets[intId]
	if !ok {
		http.Error(w, "pet not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(pet); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handlePatch(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pet, ok := pets[intId]
	if !ok {
		http.Error(w, "pet not found", http.StatusNotFound)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&pet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pets[intId] = pet //TODO: should we map the individual fields?
	w.WriteHeader(http.StatusNoContent)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	intId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	delete(pets, intId)
	w.WriteHeader(http.StatusNoContent)
}
