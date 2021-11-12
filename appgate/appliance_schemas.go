package appgate

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func allowSourcesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"address": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validateIPaddress,
				},
				"netmask": {
					Type:     schema.TypeInt,
					Optional: true,
				},
				"nic": {
					Type:     schema.TypeString,
					Optional: true,
				},
			},
		},
	}
}

func controllerSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Optional: true,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enabled": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	}
}

func adminInterfaceSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"hostname": {
					Type:     schema.TypeString,
					Required: true,
				},
				"https_port": {
					Type:     schema.TypeInt,
					Default:  443,
					Optional: true,
				},
				"https_ciphers": {
					Type:        schema.TypeList,
					Description: "The type of TLS ciphers to allow. See: https://www.openssl.org/docs/man1.0.2/apps/ciphers.html for all supported ciphers.",
					Optional:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
				"allow_sources": allowSourcesSchema(),
			},
		},
	}
}
