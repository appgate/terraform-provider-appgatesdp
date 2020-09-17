package appgate

import (
	"context"
	"fmt"
	"log"
	"strconv"
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
				Computed: true,
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
				Computed: true,
				Optional: true,
			}
			s["username_attribute"] = &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			}
			s["membership_filter"] = &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			}
			s["membership_base_dn"] = &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			}
			s["password_warning"] = &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
							Optional: true,
						},
						"threshold_days": {
							Type:     schema.TypeInt,
							Computed: true,
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
	p, _, err := request.IdentityProvider(*args).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create %s provider %+v", identityProviderLdap, prettyPrintAPIError(err))
	}
	d.SetId(p.Id)
	return resourceAppgateLdapProviderRuleRead(d, meta)
}

func resourceAppgateLdapProviderRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading ldap identity provider id: %+v", d.Id())

	token := meta.(*Client).Token
	api := meta.(*Client).API.LdapIdentityProvidersApi
	ctx := context.TODO()
	request := api.IdentityProvidersIdGet(ctx, d.Id())
	ldap, _, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Failed to read LDAP Identity provider, %+v", err)
	}
	log.Printf("[DEBUG] Reading ldap identity provider POOL IPV4: %+v", ldap.GetIpPoolV4())
	d.Set("type", identityProviderLdap)
	// base attributes
	d.Set("name", ldap.Name)
	d.Set("notes", ldap.Notes)
	d.Set("tags", ldap.Tags)

	// identity provider attributes
	d.Set("default", ldap.GetDefault())
	d.Set("client_provider", ldap.GetClientProvider())
	d.Set("admin_provider", ldap.GetAdminProvider())
	if v, ok := ldap.GetOnBoarding2FAOk(); ok {
		if err := d.Set("on_boarding_2fa", flattenIdentityProviderOnboarding2fa(*v)); err != nil {
			return err
		}
	}

	d.Set("on_boarding_type", ldap.GetOnBoardingType())
	d.Set("on_boarding_otp_provider", ldap.GetOnBoardingOtpProvider())
	d.Set("on_boarding_otp_message", ldap.GetOnBoardingOtpMessage())
	d.Set("inactivity_timeout_minutes", ldap.GetInactivityTimeoutMinutes())
	if v, ok := ldap.GetIpPoolV4Ok(); ok {
		d.Set("ip_pool_v4", *v)
	}
	if v, ok := ldap.GetIpPoolV6Ok(); ok {
		d.Set("ip_pool_v6", v)
	}

	d.Set("dns_servers", ldap.GetDnsServers())
	d.Set("dns_search_domains", ldap.GetDnsSearchDomains())
	d.Set("block_local_dns_requests", ldap.GetBlockLocalDnsRequests())
	if v, ok := ldap.GetClaimMappingsOk(); ok {
		if err := d.Set("claim_mappings", flattenIdentityProviderClaimsMappning(*v)); err != nil {
			return err
		}
	}
	if v, ok := ldap.GetOnDemandClaimMappingsOk(); ok {
		d.Set("on_demand_claim_mappings", flattenIdentityProviderOnDemandClaimsMappning(*v))
	}
	// ldap attributes
	d.Set("hostnames", ldap.GetHostnames())
	d.Set("port", ldap.GetPort())
	d.Set("ssl_enabled", ldap.GetSslEnabled())
	d.Set("admin_distinguished_name", ldap.GetAdminDistinguishedName())
	if val, ok := d.GetOk("admin_password"); ok {
		d.Set("admin_password", val.(string))
	}
	if v, ok := ldap.GetBaseDnOk(); ok {
		d.Set("base_dn", v)
	}

	d.Set("object_class", ldap.GetObjectClass())
	d.Set("username_attribute", ldap.GetUsernameAttribute())
	d.Set("membership_filter", ldap.GetMembershipFilter())
	if v, ok := ldap.GetMembershipBaseDnOk(); ok {
		if err := d.Set("membership_base_dn", &v); err != nil {
			return err
		}
	}
	if v, ok := ldap.GetPasswordWarningOk(); ok {
		if err := d.Set("password_warning", flattenLdapPasswordWarning(*v)); err != nil {
			return err
		}
	}
	return nil
}

func flattenLdapPasswordWarning(pw openapi.LdapProviderAllOfPasswordWarning) map[string]interface{} {
	// TODO: wrong types on threshold_days, enabled
	o := make(map[string]interface{})
	if v, ok := pw.GetEnabledOk(); ok {
		o["enabled"] = strconv.FormatBool(*v) // TODOD should be bool
	}
	if v, ok := pw.GetThresholdDaysOk(); ok {
		o["threshold_days"] = strconv.Itoa(int(*v)) // TOOD Should be int
	}
	if v, ok := pw.GetMessageOk(); ok {
		o["message"] = v
	}
	return o
}

func resourceAppgateLdapProviderRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating ldap identity provider id: %+v", d.Id())
	return resourceAppgateLdapProviderRuleRead(d, meta)
}
