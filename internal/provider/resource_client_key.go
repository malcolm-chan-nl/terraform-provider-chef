package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	chefc "github.com/go-chef/chef"
)

func resourceChefClientKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: CreateClientKey,
		UpdateContext: UpdateClientKey,
		ReadContext:   ReadClientKey,
		DeleteContext: DeleteClientKey,

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

func CreateClientKey(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*chefClient)

	key, err := clientKeyFromResourceData(d)
	if err != nil {
		return err
	}

	if _, err := c.Clients.AddKey(key.Client, key.Key); err != nil {
		return diag.Diagnostics{
			{
				Severity:      diag.Error,
				Summary:       "Error creating client key",
				Detail:        fmt.Sprint(err),
				AttributePath: cty.GetAttrPath("key_name"),
			},
		}
	}

	d.SetId(key.Client + "+" + key.Key.Name)
	return ReadClientKey(ctx, d, meta)
}

func UpdateClientKey(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*chefClient)

	key, err := clientKeyFromResourceData(d)
	if err != nil {
		return err
	}

	if _, err := c.Clients.UpdateKey(key.Client, key.Key.Name, key.Key); err != nil {
		return diag.Diagnostics{
			{
				Severity:      diag.Error,
				Summary:       "Error updating client key",
				Detail:        fmt.Sprint(err),
				AttributePath: cty.GetAttrPath("key_name"),
			},
		}
	}

	d.SetId(key.Client + "+" + key.Key.Name)
	return ReadClientKey(ctx, d, meta)
}

func ReadClientKey(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*chefClient)

	key, err := clientKeyFromResourceData(d)
	if err != nil {
		return err
	}

	if k, err := c.Clients.GetKey(key.Client, key.Key.Name); err == nil {
		d.Set("client", key.Client)
		d.Set("key_name", k.Name)
		d.Set("public_key", k.PublicKey)
	} else {
		if errRes, ok := err.(*chefc.ErrorResponse); ok {
			if errRes.Response.StatusCode == 404 {
				d.SetId("")
			}
		} else {
			return diag.Diagnostics{
				{
					Severity:      diag.Error,
					Summary:       "Error reading client key",
					Detail:        fmt.Sprint(err),
					AttributePath: cty.GetAttrPath("key_name"),
				},
			}
		}
	}
	return nil
}

func DeleteClientKey(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*chefClient)

	key, err := clientKeyFromResourceData(d)
	if err != nil {
		return err
	}
	if _, err := c.Clients.DeleteKey(key.Client, key.Key.Name); err == nil {
		d.SetId("")
		return nil
	} else {
		return diag.Diagnostics{
			{
				Severity:      diag.Error,
				Summary:       "Error deleting client key",
				Detail:        fmt.Sprint(err),
				AttributePath: cty.GetAttrPath("key_name"),
			},
		}
	}
}

func clientKeyFromResourceData(d *schema.ResourceData) (*chefClientKey, diag.Diagnostics) {
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
