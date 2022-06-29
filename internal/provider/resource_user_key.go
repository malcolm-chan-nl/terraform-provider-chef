package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	chefc "github.com/go-chef/chef"
)

func resourceChefUserKey() *schema.Resource {
	return &schema.Resource{
		Create: CreateUserKey,
		Update: UpdateUserKey,
		Read:   ReadUserKey,
		Delete: DeleteUserKey,

		Schema: map[string]*schema.Schema{
			"user": {
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

type chefUserKey struct {
	User string
	Key  chefc.AccessKey
}

func CreateUserKey(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*chefClient)

	key, err := userKeyFromResourceData(d)
	if err != nil {
		return err
	}

	if _, err := c.Global.Users.AddKey(key.User, key.Key); err != nil {
		return err
	}

	d.SetId(key.User + "+" + key.Key.Name)
	return ReadUserKey(d, meta)
}

func UpdateUserKey(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*chefClient)

	key, err := userKeyFromResourceData(d)
	if err != nil {
		return err
	}

	if _, err = c.Global.Users.UpdateKey(key.User, key.Key.Name, key.Key); err != nil {
		return err
	}

	d.SetId(key.User + "+" + key.Key.Name)
	return ReadUserKey(d, meta)
}

func ReadUserKey(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*chefClient)

	key, err := userKeyFromResourceData(d)
	if err != nil {
		return err
	}

	k, err := c.Global.Users.GetKey(key.User, key.Key.Name)
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

	d.Set("user", key.User)
	d.Set("key_name", k.Name)
	d.Set("public_key", k.PublicKey)

	return nil
}

func DeleteUserKey(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*chefClient)

	key, err := userKeyFromResourceData(d)
	if err != nil {
		return err
	}
	if _, err = c.Global.Users.DeleteKey(key.User, key.Key.Name); err == nil {
		d.SetId("")
	}

	return err
}

func userKeyFromResourceData(d *schema.ResourceData) (*chefUserKey, error) {
	key := &chefUserKey{
		User: d.Get("user").(string),
		Key: chefc.AccessKey{
			Name:           d.Get("key_name").(string),
			PublicKey:      d.Get("public_key").(string),
			ExpirationDate: "infinity",
		},
	}
	return key, nil
}
