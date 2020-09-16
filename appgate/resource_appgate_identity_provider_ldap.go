package appgate

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/appgate/terraform-provider-appgate/client/v12/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAppgateLdapProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateLdapProviderRuleCreate,
		Read:   resourceAppgateLdapProviderRuleRead,
		Update: resourceAppgateLdapProviderRuleUpdate,
		Delete: identityProviderDelete,
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
			s["type"].Default = identityProviderLdap

			s["hostnames"] = &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
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
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
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
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"threshold_days": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"message": {
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
	api := meta.(*Client).API.LdapIdentityProvidersApi
	ctx := context.TODO()
	provider := openapi.NewIdentityProvider(identityProviderLdap)
	provider, err := readProviderFromConfig(d, *provider)
	if err != nil {
		return fmt.Errorf("Failed to read and create basic identity provider for %s %s", identityProviderLdap, err)
	}

	args := openapi.NewLdapProviderWithDefaults()
	args.SetType(provider.GetType())
	args.SetId(provider.GetId())
	args.SetName(provider.GetName())
	args.SetNotes(provider.GetNotes())
	args.SetTags(provider.GetTags())

	if provider.Default != nil {
		args.SetDefault(provider.GetDefault())
	}
	if provider.ClientProvider != nil {
		args.SetClientProvider(*provider.ClientProvider)
	}
	if provider.AdminProvider != nil {
		args.SetAdminProvider(*provider.AdminProvider)
	}
	if provider.OnBoarding2FA != nil {
		args.SetOnBoarding2FA(*provider.OnBoarding2FA)
	}
	if provider.OnBoardingType != nil {
		args.SetOnBoardingType(*provider.OnBoardingType)
	}
	if provider.OnBoardingOtpProvider != nil {
		args.SetOnBoardingOtpProvider(*provider.OnBoardingOtpProvider)
	}
	if provider.OnBoardingOtpMessage != nil {
		args.SetOnBoardingOtpMessage(*provider.OnBoardingOtpMessage)
	}
	if provider.InactivityTimeoutMinutes != nil {
		args.SetInactivityTimeoutMinutes(*provider.InactivityTimeoutMinutes)
	}
	if provider.IpPoolV4 != nil {
		args.SetIpPoolV4(*provider.IpPoolV4)
	}
	if provider.IpPoolV6 != nil {
		args.SetIpPoolV6(*provider.IpPoolV6)
	}
	if provider.DnsServers != nil {
		args.SetDnsServers(*provider.DnsServers)
	}
	if provider.DnsSearchDomains != nil {
		args.SetDnsSearchDomains(*provider.DnsSearchDomains)
	}
	if provider.BlockLocalDnsRequests != nil {
		args.SetBlockLocalDnsRequests(*provider.BlockLocalDnsRequests)
	}
	if provider.ClaimMappings != nil {
		args.SetClaimMappings(*provider.ClaimMappings)
	}
	if provider.OnDemandClaimMappings != nil {
		args.SetOnDemandClaimMappings(*provider.OnDemandClaimMappings)
	}
	if v, ok := d.GetOk("hostnames"); ok {
		hostnames, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return err
		}
		args.SetHostnames(hostnames)
	}
	if v, ok := d.GetOk("port"); ok {
		args.SetPort(int32(v.(int)))
	}
	if v, ok := d.GetOk("ssl_enabled"); ok {
		args.SetSslEnabled(v.(bool))
	}
	if v, ok := d.GetOk("admin_distinguished_name"); ok {
		args.SetAdminDistinguishedName(v.(string))
	}
	if v, ok := d.GetOk("admin_password"); ok {
		args.SetAdminPassword(v.(string))
	}
	if v, ok := d.GetOk("base_dn"); ok {
		args.SetBaseDn(v.(string))
	}
	if v, ok := d.GetOk("object_class"); ok {
		args.SetObjectClass(v.(string))
	}
	if v, ok := d.GetOk("username_attribute"); ok {
		args.SetUsernameAttribute(v.(string))
	}
	if v, ok := d.GetOk("membership_filter"); ok {
		args.SetMembershipFilter(v.(string))
	}
	if v, ok := d.GetOk("membership_base_dn"); ok {
		args.SetMembershipBaseDn(v.(string))
	}
	if v, ok := d.GetOk("password_warning"); ok {
		pw := openapi.LdapProviderAllOfPasswordWarning{}
		raw := v.(map[string]interface{})
		if v, ok := raw["enabled"]; ok {
			pw.SetEnabled(v.(bool))
		}
		if v, ok := raw["threshold_days"]; ok {
			pw.SetThresholdDays(int32(v.(int)))
		}
		if v, ok := raw["message"]; ok {
			pw.SetMessage(v.(string))
		}
		args.SetPasswordWarning(pw)
	}
	request := api.IdentityProvidersPost(ctx)
	_, _, err = request.IdentityProvider(*args).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create %s provider %+v", identityProviderLdap, prettyPrintAPIError(err))
	}
	return resourceAppgateLdapProviderRuleRead(d, meta)
}

func resourceAppgateLdapProviderRuleRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAppgateLdapProviderRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAppgateLdapProviderRuleRead(d, meta)
}
