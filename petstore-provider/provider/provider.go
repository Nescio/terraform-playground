package provider

import (
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/Nescio/terraform-playground/sdk"
)

func New() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PETSTORE_ADDRESS", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"petstore_pet": resourcePet(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ConfigureFunc:  providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	hostname, _ := d.Get("address").(string)
	address, _ := url.Parse(hostname)
	cfg := &sdk.Config{
		Address: address.String(),
	}
	return sdk.NewClient(cfg)
}
