package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	chefc "github.com/go-chef/chef"
)

func dataChefSearch() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataChefSearchRead,

		Schema: map[string]*schema.Schema{
			"index": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "node",
			},
			"query": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"filter": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"unique": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"result": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"total_num": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataChefSearchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*chefc.Client)

	query, err := client.Search.NewQuery(d.Get("index").(string), d.Get("query").(string))
	if err != nil {
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  "Error creating search query",
				Detail:   fmt.Sprint(err),
			},
		}
	}
	query.Rows = 1

	filter, ok := d.Get("filter").(*schema.Set)
	var res chefc.SearchResult
	if ok {
		params := make(map[string]interface{})
		for _, v := range filter.List() {
			m := v.(map[string]interface{})
			params[m["name"].(string)] = m["value"].([]interface{})
		}
		res, err = query.DoPartial(client, params)
	} else {
		res, err = query.Do(client)
	}
	if err != nil {
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  "Error executing search",
				Detail:   fmt.Sprint(err),
			},
		}
	}

	log.Printf("Chef search result: %+v\n", res)
	d.SetId("static")
	d.Set("total_num", res.Total)
	if d.Get("unique").(bool) && res.Total != 1 {
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  "Error executing search",
				Detail:   fmt.Sprintf("Query had %d results, not one.", res.Total),
			},
		}
	}
	if res.Total > 0 {
		result := make(map[string]string)
		row := res.Rows[0].(map[string]interface{})
		// For some indexes the data is returned in data and for others in raw_data
		data, ok := row["data"]
		if !ok {
			data = row["raw_data"]
		}
		for k, v := range data.(map[string]interface{}) {
			switch t := v.(type) {
			case string:
				result[k] = t
			default:
				result[k] = fmt.Sprintf("%v", t)
			}
		}
		d.Set("result", result)
	}
	return nil
}
