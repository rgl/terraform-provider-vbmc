package vbmc

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVbmc() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVbmcCreate,
		ReadContext:   resourceVbmcRead,
		DeleteContext: resourceVbmcDelete,
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"address": {
				Type:     schema.TypeString,
				Default:  "127.0.0.1",
				Optional: true,
				ForceNew: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"username": {
				Type:     schema.TypeString,
				Default:  "admin",
				Optional: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Default:   "password",
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceVbmcCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	domainName := d.Get("domain_name").(string)
	address := d.Get("address").(string)
	port := d.Get("port").(int)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	_, err := Create(domainName, address, port, username, password)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domainName)

	return resourceVbmcRead(ctx, d, m)
}

func resourceVbmcRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	vbmc, err := Get(d.Id())
	if err != nil {
		if execError, ok := err.(*VbmcExecError); ok {
			if strings.Contains(execError.Stderr, "No domain with matching name") {
				d.SetId("")
				return diag.Diagnostics{}
			}
		}
		return diag.FromErr(err)
	}

	d.Set("port", vbmc.Port)

	return diag.Diagnostics{}
}

func resourceVbmcDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	err := Delete(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}
