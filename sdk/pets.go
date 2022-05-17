package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
)

// Compile-time proof of interface implementation.
var _ Pets = (*pets)(nil)

// Pets describes all the pet related methods that the Petstore API supports.
type Pets interface {
	// Create a new pet with the given options.
	Create(options PetCreateOptions) (*Pet, error)

	// Read an pet by its ID.
	Read(petID string) (*Pet, error)

	// Update an pet by its ID.
	Update(petID string, options PetUpdateOptions) (*Pet, error)

	// Delete an pet by its ID.
	Delete(petID string) error
}

func newPets(client *Client) *pets {
	return &pets{
		ctx:    context.Background(),
		client: client,
		path:   "pets",
	}
}

// pets implements Pets.
type pets struct {
	ctx    context.Context
	client *Client
	path   string
}

// Pet represents a Petstore pet.
type Pet struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Species string `json:"species"`
	Age     int    `json:"age"`
}

func (p Pet) MarshalJSON() ([]byte, error) {
	id, err := strconv.Atoi(p.ID)
	if err != nil {
		return nil, err
	}
	return json.Marshal([]interface{}{
		id,
		p.Name,
		p.Species,
		p.Age,
	})
}

func (p *Pet) UnmarshalJSON(data []byte) error {
	log.Println("UnmarshalJSON:", string(data))

	var v []interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if len(v) != 4 {
		return errors.New("invalid pet")
	}

	id, ok := v[0].(int)
	if !ok {
		return errors.New("invalid pet id")
	}

	p.ID = strconv.Itoa(id)
	p.Name, _ = v[1].(string)
	p.Species = v[2].(string)
	p.Age = v[3].(int)

	return nil
}

// PetCreateOptions represents the options for creating an pet.
type PetCreateOptions struct {
	Name    string `json:"name"`
	Species string `json:"species"`
	Age     int    `json:"age"`
}

func (o PetCreateOptions) valid() error {
	if !validString(&o.Name) {
		return errors.New("pet name is required")
	}
	if !validString(&o.Species) {
		return errors.New("pet species is required")
	}
	if !validInt(&o.Age) {
		return errors.New("age is required")
	}
	return nil
}

// Create a new pet with the given options.
func (p *pets) Create(options PetCreateOptions) (*Pet, error) {
	if err := options.valid(); err != nil {
		return nil, err
	}

	req, err := p.client.newRequest("POST", p.path, &options)
	if err != nil {
		return nil, err
	}

	v := &Pet{}
	err = p.client.do(p.ctx, req, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Read an pet by id.
func (p *pets) Read(id string) (*Pet, error) {
	if !validID(&id) {
		return nil, errors.New("invalid id")
	}

	path := fmt.Sprintf("%s/%s", p.path, url.QueryEscape(id))
	req, err := p.client.newRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	v := &Pet{}
	err = p.client.do(p.ctx, req, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// PetUpdateOptions represents the options for updating an pet.
type PetUpdateOptions struct {
	Name    string `json:"name"`
	Species string `json:"species"`
	Age     int    `json:"age"`
}

// Update attributes of an existing pet.
func (p *pets) Update(id string, options PetUpdateOptions) (*Pet, error) {
	if !validID(&id) {
		return nil, errors.New("invalid id")
	}

	path := fmt.Sprintf("%s/%s", p.path, url.QueryEscape(id))
	req, err := p.client.newRequest("PATCH", path, &options)
	if err != nil {
		return nil, err
	}

	v := &Pet{}
	err = p.client.do(p.ctx, req, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

// Delete an pet by its ID.
func (p *pets) Delete(id string) error {
	if !validID(&id) {
		return errors.New("invalid id")
	}

	path := fmt.Sprintf("%s/%s", p.path, url.QueryEscape(id))
	req, err := p.client.newRequest("DELETE", path, nil)
	if err != nil {
		return err
	}

	return p.client.do(p.ctx, req, nil)
}
