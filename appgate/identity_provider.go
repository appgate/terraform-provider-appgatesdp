package appgate

import (
	"context"
	"fmt"

	"github.com/appgate/terraform-provider-appgate/client/v12/openapi"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func identityProviderSchema() map[string]*schema.Schema {
	return mergeSchemaMaps(baseEntitySchema(), map[string]*schema.Schema{
		"type": {
			Type:     schema.TypeString,
			Computed: true,
			ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
				s := v.(string)
				list := []string{"LocalDatabase", "Radius", "Ldap", "Saml", "LdapCertificate", "IotConnector"}
				for _, x := range list {
					if s == x {
						return
					}
				}
				errs = append(errs, fmt.Errorf("type must be on of %v, got %s", list, s))
				return
			},
		},

		"default": {
			Type:     schema.TypeBool,
			Optional: true,
		},

		"admin_provider": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"on_boarding_2fa": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"mfa_provider_id": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"message": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
					"device_limit_per_user": &schema.Schema{
						Type:     schema.TypeInt,
						Optional: true,
					},
				},
			},
		},

		"inactivity_timeout_minutes": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"ip_pool_v4": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"ip_pool_v6": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"dns_servers": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"dns_search_domains": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"block_local_dns_requests": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"claim_mappings": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"attribute_name": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"claim_name": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"list": &schema.Schema{
						Type:     schema.TypeBool,
						Optional: true,
					},
					"encrypted": &schema.Schema{
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"on_demand_claim_mappings": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"command": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
						ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
							s := v.(string)
							list := []string{
								"fileSize",
								"fileExists",
								"fileCreated",
								"fileUpdated",
								"fileVersion",
								"fileSha512",
								"processRunning",
								"processList",
								"serviceRunning",
								"serviceList",
								"regExists",
								"regQuery",
								"runScript",
							}
							for _, x := range list {
								if s == x {
									return
								}
							}
							errs = append(errs, fmt.Errorf("type must be on of %v, got %s", list, s))
							return
						},
					},
					"claim_name": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"parameters": &schema.Schema{
						Type:     schema.TypeMap,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"name": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
								},
								"path": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
								},
								"args": &schema.Schema{
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
					"platform": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
						ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
							s := v.(string)
							list := []string{
								"desktop.windows.all",
								"desktop.macos.all",
								"desktop.linux.all",
								"desktop.all",
								"mobile.android.all",
								"mobile.ios.all",
								"mobile.all",
								"all",
							}
							for _, x := range list {
								if s == x {
									return
								}
							}
							errs = append(errs, fmt.Errorf("type must be on of %v, got %s", list, s))
							return
						},
					},
				},
			},
		},
	})
}

func identityProviderPost(ctx context.Context, api *openapi.IdentityProvidersApiService, token string, body openapi.IdentityProvider) (interface{}, error) {
	request := api.IdentityProvidersPost(ctx)
	request = request.IdentityProvider(body)
	provider, _, err := request.Authorization(token).Execute()
	if err != nil {
		return nil, fmt.Errorf("Could not create provider %+v", prettyPrintAPIError(err))
	}
	return provider, nil
}
