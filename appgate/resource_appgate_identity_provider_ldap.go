package appgate

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/appgate/terraform-provider-appgate/client/v12/openapi"
	"github.com/google/uuid"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAppgateLdapProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateLdapProviderRuleCreate,
		Read:   resourceAppgateLdapProviderRuleRead,
		Update: resourceAppgateLdapProviderRuleUpdate,
		Delete: resourceAppgateLdapProviderRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: func() map[string]*schema.Schema {
			s := identityProviderSchema()
			s["type"].Default = "ldap"

			s["hostnames"] = &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			}
			s["port"] = &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			}
			s["ssl_enabled"] = &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			}
			s["admin_distinguished_name"] = &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			}
			s["admin_password"] = &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			}
			s["base_dn"] = &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			}
			s["object_class"] = &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			}
			s["username_attribute"] = &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			}
			s["membership_filter"] = &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			}
			s["membership_base_dn"] = &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			}
			s["password_warning"] = &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
						},
						"threshold_days": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"message": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			}
			return s
		}(),
	}
}

func resourceAppgateLdapProviderRuleCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating LdapProvider: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.IdentityProvidersApi
	ctx := context.TODO()
	args := openapi.NewLdapProviderWithDefaults()
	args.Id = uuid.New().String()

	request := api.IdentityProvidersPost(ctx)
	// request = request.IdentityProvider(*args.(openapi.IdentityProvider))
	provider, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create provider %+v", prettyPrintAPIError(err))
	}
	fmt.Printf("LDAP %+v", provider)
	return nil
}

func resourceAppgateLdapProviderRuleRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAppgateLdapProviderRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAppgateLdapProviderRuleRead(d, meta)
}

func resourceAppgateLdapProviderRuleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete LdapProvider: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.IdentityProvidersApi

	request := api.IdentityProvidersIdDelete(context.Background(), d.Id())

	_, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not delete LdapProvider %+v", prettyPrintAPIError(err))
	}
	d.SetId("")
	return nil
}
