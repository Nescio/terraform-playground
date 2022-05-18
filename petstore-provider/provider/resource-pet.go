package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/Nescio/terraform-playground/sdk"
)

func resourcePet() *schema.Resource {
	return &schema.Resource{
		Create: resourcePetCreate,
		Read:   resourcePetRead,
		Update: resourcePetUpdate,
		Delete: resourcePetDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"species": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"age": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourcePetCreate(d *schema.ResourceData, m interface{}) error {
	conn := m.(*sdk.Client)
	options := sdk.PetCreateOptions{
		Name:    d.Get("name").(string),
		Species: d.Get("species").(string),
		Age:     d.Get("age").(int),
	}

	pet, err := conn.Pets.Create(options)
	if err != nil {
		return err
	}

	d.SetId(string(pet.ID))
	return resourcePetRead(d, m)
}

func resourcePetRead(d *schema.ResourceData, m interface{}) error {
	conn := m.(*sdk.Client)
	pet, err := conn.Pets.Read(d.Id())
	if err != nil {
		return err
	}
	d.Set("name", pet.Name)
	d.Set("species", pet.Species)
	d.Set("age", pet.Age)
	return nil
}

func resourcePetUpdate(d *schema.ResourceData, m interface{}) error {
	conn := m.(*sdk.Client)
	options := sdk.PetUpdateOptions{}
	if d.HasChange("name") {
		options.Name = d.Get("name").(string)
	}
	if d.HasChange("age") {
		options.Age = d.Get("age").(int)
	}
	if _, err := conn.Pets.Update(d.Id(), options); err != nil {
		return err
	}
	return resourcePetRead(d, m)
}

func resourcePetDelete(d *schema.ResourceData, m interface{}) error {
	conn := m.(*sdk.Client)
	return conn.Pets.Delete(d.Id())
}
