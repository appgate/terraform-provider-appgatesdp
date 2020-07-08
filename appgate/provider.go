package appgate

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider function returns the object that implements the terraform.ResourceProvider interface, specifically a schema.Provider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPGATE_ADDRESS", ""),
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPGATE_USERNAME", ""),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPGATE_PASSWORD", ""),
			},
			"provider": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPGATE_PROVIDER", "local"),
			},
			"insecure": {
				Type:        schema.TypeBool,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPGATE_INSECURE", true),
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"appgate_appliance":       dataSourceAppgateAppliance(),
			"appgate_site":            dataSourceAppgateSite(),
			"appgate_condition":       dataSourceAppgateCondition(),
			"appgate_policy":          dataSourceAppgatePolicy(),
			"appgate_ringfence_rule":  dataSourceAppgateRingfenceRule(),
			"appgate_criteria_script": dataSourceCriteriaScript(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"appgate_appliance":       resourceAppgateAppliance(),
			"appgate_entitlement":     resourceAppgateEntitlement(),
			"appgate_site":            resourceAppgateSite(),
			"appgate_ringfence_rule":  resourceAppgateRingfenceRule(),
			"appgate_condition":       resourceAppgateCondition(),
			"appgate_policy":          resourceAppgatePolicy(),
			"appgate_criteria_script": resourceAppgateCriteriaScript(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		URL:      d.Get("url").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		Provider: d.Get("provider").(string),
		Insecure: d.Get("insecure").(bool),
		Timeout:  20,
	}
	return config.Client()
}
