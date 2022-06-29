package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	chefc "github.com/go-chef/chef"
)

func resourceChefClientKey() *schema.Resource {
	return &schema.Resource{
		Create: CreateClientKey,
		Update: UpdateClientKey,
		Read:   ReadClientKey,
		Delete: DeleteClientKey,

		Schema: map[string]*schema.Schema{
			"client": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"public_key": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

type chefClientKey struct {
	Client string
	Key    chefc.AccessKey
}

func CreateClientKey(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*chefClient)

	key, err := clientKeyFromResourceData(d)
	if err != nil {
		return err
	}

	if _, err := c.Clients.AddKey(key.Client, key.Key); err != nil {
		return err
	}

	d.SetId(key.Client + "+" + key.Key.Name)
	return ReadClientKey(d, meta)
}

func UpdateClientKey(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*chefClient)

	key, err := clientKeyFromResourceData(d)
	if err != nil {
		return err
	}

	if _, err = c.Clients.UpdateKey(key.Client, key.Key.Name, key.Key); err != nil {
		return err
	}

	d.SetId(key.Client + "+" + key.Key.Name)
	return ReadClientKey(d, meta)
}

func ReadClientKey(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*chefClient)

	key, err := clientKeyFromResourceData(d)
	if err != nil {
		return err
	}

	k, err := c.Clients.GetKey(key.Client, key.Key.Name)
	if err != nil {
		if errRes, ok := err.(*chefc.ErrorResponse); ok {
			if errRes.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		} else {
			return err
		}
	}

	d.Set("client", key.Client)
	d.Set("key_name", k.Name)
	d.Set("public_key", k.PublicKey)

	return nil
}

func DeleteClientKey(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*chefClient)

	key, err := clientKeyFromResourceData(d)
	if err != nil {
		return err
	}
	if _, err = c.Clients.DeleteKey(key.Client, key.Key.Name); err == nil {
		d.SetId("")
	}

	return err
}

func clientKeyFromResourceData(d *schema.ResourceData) (*chefClientKey, error) {
	key := &chefClientKey{
		Client: d.Get("client").(string),
		Key: chefc.AccessKey{
			Name:           d.Get("key_name").(string),
			PublicKey:      d.Get("public_key").(string),
			ExpirationDate: "infinity",
		},
	}
	return key, nil
}
