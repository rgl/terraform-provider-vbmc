package vbmc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"vbmc_vbmc": resourceVbmc(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
	}
}
