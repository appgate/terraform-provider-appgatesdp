package appgate

import (
	"fmt"
	"log"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v22/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateRadiusProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateRadiusProviderRuleCreate,
		Read:   resourceAppgateRadiusProviderRuleRead,
		Update: resourceAppgateRadiusProviderRuleUpdate,
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
			s["type"].Default = identityProviderRadius

			s["hostnames"] = &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			}
			s["port"] = &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			}
			s["shared_secret"] = &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			}
			s["authentication_protocol"] = &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "CHAP",
				ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
					s := v.(string)
					list := []string{"PAP", "CHAP"}
					for _, x := range list {
						if s == x {
							return
						}
					}
					errs = append(errs, fmt.Errorf("authentication_protocol must be on of %v, got %s", list, s))
					return
				},
			}

			return s
		}(),
	}
}

func resourceAppgateRadiusProviderRuleCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating RadiusProvider: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.RadiusIdentityProvidersApi
	ctx := BaseAuthContext(token)
	currentVersion := meta.(*Client).ApplianceVersion
	provider := &openapi.ConfigurableIdentityProvider{}
	provider.Type = identityProviderRadius
	provider, err = readProviderFromConfig(d, *provider)
	if err != nil {
		return fmt.Errorf("Failed to read and create basic identity provider for %s %w", identityProviderRadius, err)
	}
	args := openapi.RadiusProvider{}
	// base
	args.SetType(provider.GetType())
	args.SetId(provider.GetId())
	args.SetName(provider.GetName())
	args.SetNotes(provider.GetNotes())
	args.SetTags(provider.GetTags())
	// identity provider

	if provider.AdminProvider != nil {
		args.SetAdminProvider(*provider.AdminProvider)
	}
	if provider.DeviceLimitPerUser != nil {
		args.SetDeviceLimitPerUser(*provider.DeviceLimitPerUser)
	}
	if provider.OnBoarding2FA != nil {
		args.SetOnBoarding2FA(*provider.OnBoarding2FA)
	}
	if provider.InactivityTimeoutMinutes != nil {
		args.SetInactivityTimeoutMinutes(*provider.InactivityTimeoutMinutes)
	}
	if provider.NetworkInactivityTimeoutEnabled != nil {
		if currentVersion.LessThan(Appliance61Version) {
			return ErrNetworkInactivityTimeoutEnabled
		}
		args.SetNetworkInactivityTimeoutEnabled(provider.GetNetworkInactivityTimeoutEnabled())
	}
	if provider.IpPoolV4 != nil {
		args.SetIpPoolV4(*provider.IpPoolV4)
	}
	if provider.IpPoolV6 != nil {
		args.SetIpPoolV6(*provider.IpPoolV6)
	}
	if provider.UserScripts != nil {
		args.SetUserScripts(provider.GetUserScripts())
	}
	if provider.DnsServers != nil {
		args.SetDnsServers(provider.GetDnsServers())
	}
	if provider.DnsSearchDomains != nil {
		args.SetDnsSearchDomains(provider.GetDnsSearchDomains())
	}
	if provider.BlockLocalDnsRequests != nil {
		args.SetBlockLocalDnsRequests(*provider.BlockLocalDnsRequests)
	}
	if provider.ClaimMappings != nil {
		args.SetClaimMappings(provider.GetClaimMappings())
	}
	if provider.OnDemandClaimMappings != nil {
		args.SetOnDemandClaimMappings(provider.GetOnDemandClaimMappings())
	}
	// radius
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
	if v, ok := d.GetOk("shared_secret"); ok {
		args.SetSharedSecret(v.(string))
	}
	if v, ok := d.GetOk("authentication_protocol"); ok {
		args.SetAuthenticationProtocol(v.(string))
	}

	request := api.IdentityProvidersPost(ctx)
	p, _, err := request.Body(args).Execute()
	if err != nil {
		return fmt.Errorf("Could not create %s provider %w", identityProviderRadius, prettyPrintAPIError(err))
	}
	d.SetId(p.GetId())
	return resourceAppgateRadiusProviderRuleRead(d, meta)
}

func resourceAppgateRadiusProviderRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading radius identity provider id: %+v", d.Id())

	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.RadiusIdentityProvidersApi
	ctx := BaseAuthContext(token)
	request := api.IdentityProvidersIdGet(ctx, d.Id())
	radius, _, err := request.Execute()
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Failed to read LDAP Identity provider, %w", err)
	}
	d.Set("type", identityProviderRadius)
	// base attributes
	d.Set("name", radius.Name)
	d.Set("notes", radius.Notes)
	d.Set("tags", radius.Tags)

	// identity provider attributes

	d.Set("admin_provider", radius.GetAdminProvider())
	if v, ok := radius.GetDeviceLimitPerUserOk(); ok {
		d.Set("device_limit_per_user", *v)
	}
	if v, ok := radius.GetOnBoarding2FAOk(); ok {
		if err := d.Set("on_boarding_two_factor", flattenIdentityProviderOnboarding2fa(*v)); err != nil {
			return err
		}
	}

	d.Set("inactivity_timeout_minutes", radius.GetInactivityTimeoutMinutes())
	d.Set("network_inactivity_timeout_enabled", radius.GetNetworkInactivityTimeoutEnabled())
	if v, ok := radius.GetIpPoolV4Ok(); ok {
		d.Set("ip_pool_v4", *v)
	}
	if v, ok := radius.GetIpPoolV6Ok(); ok {
		d.Set("ip_pool_v6", *v)
	}

	d.Set("user_scripts", radius.GetUserScripts())
	d.Set("dns_servers", radius.GetDnsServers())
	d.Set("dns_search_domains", radius.GetDnsSearchDomains())
	d.Set("block_local_dns_requests", radius.GetBlockLocalDnsRequests())
	if v, ok := radius.GetClaimMappingsOk(); ok {
		if err := d.Set("claim_mappings", flattenIdentityProviderClaimsMappning(v)); err != nil {
			return err
		}
	}

	if v, ok := radius.GetOnDemandClaimMappingsOk(); ok {
		if err := d.Set("on_demand_claim_mappings", flattenIdentityProviderOnDemandClaimsMappning(v)); err != nil {
			return err
		}
	}

	// radius attributes
	if v, ok := radius.GetHostnamesOk(); ok {
		d.Set("hostnames", v)
	}
	if v, ok := radius.GetPortOk(); ok {
		d.Set("port", *v)
	}
	if val, ok := d.GetOk("shared_secret"); ok {
		d.Set("shared_secret", val)
	}
	if v, ok := radius.GetAuthenticationProtocolOk(); ok {
		d.Set("authentication_protocol", *v)
	}
	return nil
}

func resourceAppgateRadiusProviderRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating radius identity provider id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.RadiusIdentityProvidersApi
	ctx := BaseAuthContext(token)
	request := api.IdentityProvidersIdGet(ctx, d.Id())
	originalRadiusProvider, _, err := request.Execute()
	if err != nil {
		return fmt.Errorf("Failed to read LDAP Identity provider, %w", err)
	}
	// base attributes
	if d.HasChange("name") {
		originalRadiusProvider.SetName(d.Get("name").(string))
	}

	if d.HasChange("notes") {
		originalRadiusProvider.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("tags") {
		originalRadiusProvider.SetTags(schemaExtractTags(d))
	}

	// identity provider attributes
	if d.HasChange("admin_provider") {
		originalRadiusProvider.SetAdminProvider(d.Get("admin_provider").(bool))
	}
	if d.HasChange("device_limit_per_user") {
		originalRadiusProvider.SetDeviceLimitPerUser(int32(d.Get("device_limit_per_user").(int)))
	}
	if d.HasChange("on_boarding_two_factor") {
		_, v := d.GetChange("on_boarding_two_factor")
		onboarding, err := readOnBoardingTwoFactorFromConfig(v.([]interface{}))
		if err != nil {
			return err
		}
		originalRadiusProvider.SetOnBoarding2FA(onboarding)
	}

	if d.HasChange("inactivity_timeout_minutes") {
		originalRadiusProvider.SetInactivityTimeoutMinutes(int32(d.Get("inactivity_timeout_minutes").(int)))
	}
	if d.HasChange("network_inactivity_timeout_enabled") {
		originalRadiusProvider.SetNetworkInactivityTimeoutEnabled(d.Get("network_inactivity_timeout_enabled").(bool))
	}
	if d.HasChange("ip_pool_v4") {
		originalRadiusProvider.SetIpPoolV4(d.Get("ip_pool_v4").(string))
	}
	if d.HasChange("ip_pool_v6") {
		originalRadiusProvider.SetIpPoolV6(d.Get("ip_pool_v6").(string))
	}
	if d.HasChange("user_scripts") {
		_, v := d.GetChange("user_scripts")
		us, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("Failed to read user_scripts %w", err)
		}
		originalRadiusProvider.SetUserScripts(us)
	}
	if d.HasChange("dns_servers") {
		_, v := d.GetChange("dns_servers")
		servers, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("Failed to read dns servers %w", err)
		}
		originalRadiusProvider.SetDnsServers(servers)
	}
	if d.HasChange("dns_search_domains") {
		_, v := d.GetChange("dns_search_domains")
		servers, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("Failed to read dns search domains %w", err)
		}
		originalRadiusProvider.SetDnsSearchDomains(servers)
	}
	if d.HasChange("block_local_dns_requests") {
		originalRadiusProvider.SetBlockLocalDnsRequests(d.Get("block_local_dns_requests").(bool))
	}
	if d.HasChange("claim_mappings") {
		_, v := d.GetChange("claim_mappings")
		claims := readIdentityProviderClaimMappingFromConfig(v.(*schema.Set).List())
		originalRadiusProvider.SetClaimMappings(claims)
	}
	if d.HasChange("on_demand_claim_mappings") {
		_, v := d.GetChange("on_demand_claim_mappings")
		claims := readIdentityProviderOnDemandClaimMappingFromConfig(v.(*schema.Set).List())
		originalRadiusProvider.SetOnDemandClaimMappings(claims)
	}

	// radius provider attributes
	if d.HasChange("hostnames") {
		_, v := d.GetChange("hostnames")
		servers, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("Failed to read hostnames %w", err)
		}
		originalRadiusProvider.SetHostnames(servers)
	}
	if d.HasChange("port") {
		_, v := d.GetChange("port")
		originalRadiusProvider.SetPort(int32(v.(int)))
	}

	originalRadiusProvider.SetSharedSecret(d.Get("shared_secret").(string))

	req := api.IdentityProvidersIdPut(ctx, d.Id())
	req = req.Body(*originalRadiusProvider)
	_, _, err = req.Execute()
	if err != nil {
		return fmt.Errorf("Could not update %s provider %w", identityProviderRadius, prettyPrintAPIError(err))
	}
	return resourceAppgateRadiusProviderRuleRead(d, meta)
}
