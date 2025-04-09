package appgate

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v22/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateOidcProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateOidcProviderRuleCreate,
		Read:   resourceAppgateOidcProviderRuleRead,
		Update: resourceAppgateOidcProviderRuleUpdate,
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
			s["type"].Default = identityProviderOidc

			s["issuer"] = &schema.Schema{
				Type:        schema.TypeString,
				Description: "OIDC issuer URL",
				Required:    true,
			}

			s["audience"] = &schema.Schema{
				Type:        schema.TypeString,
				Description: "Audience/Client ID to make sure the recipient of the token is this controller",
				Required:    true,
			}

			s["scope"] = &schema.Schema{
				Type:        schema.TypeString,
				Description: "Scope to use for tokens",
				Optional:    true,
			}

			s["google"] = &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"client_secret": {
							Type:      schema.TypeString,
							Sensitive: true,
							Optional:  true,
						},
						"refresh_token": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			}
			return s
		}(),
	}
}

func resourceAppgateOidcProviderRuleCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating OidcProvider: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.OidcIdentityProvidersApi
	ctx := context.TODO()
	currentVersion := meta.(*Client).ApplianceVersion
	provider := &openapi.ConfigurableIdentityProvider{}
	provider.Type = identityProviderOidc
	provider, err = readProviderFromConfig(d, *provider)
	if err != nil {
		return fmt.Errorf("Failed to read and create basic identity provider for %s %w", identityProviderOidc, err)
	}
	args := openapi.OidcProvider{}
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

	// oidc
	if v, ok := d.GetOk("issuer"); ok {
		args.SetIssuer(v.(string))
	}
	if v, ok := d.GetOk("audience"); ok {
		args.SetAudience(v.(string))
	}
	if v, ok := d.GetOk("scope"); ok {
		args.SetScope(v.(string))
	}
	if v, ok := d.GetOk("google"); ok {
		g := openapi.NewOidcProviderAllOfGoogle()
		for _, raw := range v.([]interface{}) {
			google := raw.(map[string]interface{})
			if v, ok := google["enabled"]; ok {
				g.SetEnabled(v.(bool))
			}
			if v, ok := google["client_secret"]; ok {
				g.SetClientSecret(v.(string))
			}
			if v, ok := google["refresh_token"]; ok {
				g.SetRefreshToken(v.(bool))
			}
		}
		args.SetGoogle(*g)
	}

	request := api.IdentityProvidersPost(ctx)
	p, _, err := request.Body(args).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create %s provider %w", identityProviderOidc, prettyPrintAPIError(err))
	}
	d.SetId(p.GetId())
	return resourceAppgateOidcProviderRuleRead(d, meta)
}

func resourceAppgateOidcProviderRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading oidc identity provider id: %+v", d.Id())

	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.OidcIdentityProvidersApi
	ctx := context.TODO()
	request := api.IdentityProvidersIdGet(ctx, d.Id())
	oidc, _, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Failed to read LDAP Identity provider, %w", err)
	}
	d.Set("type", identityProviderOidc)
	// base attributes
	d.Set("name", oidc.Name)
	d.Set("notes", oidc.Notes)
	d.Set("tags", oidc.Tags)

	// identity provider attributes

	d.Set("admin_provider", oidc.GetAdminProvider())
	if v, ok := oidc.GetDeviceLimitPerUserOk(); ok {
		d.Set("device_limit_per_user", *v)
	}
	if v, ok := oidc.GetOnBoarding2FAOk(); ok {
		if err := d.Set("on_boarding_two_factor", flattenIdentityProviderOnboarding2fa(*v)); err != nil {
			return err
		}
	}

	d.Set("inactivity_timeout_minutes", oidc.GetInactivityTimeoutMinutes())
	d.Set("network_inactivity_timeout_enabled", oidc.GetNetworkInactivityTimeoutEnabled())
	if v, ok := oidc.GetIpPoolV4Ok(); ok {
		d.Set("ip_pool_v4", *v)
	}
	if v, ok := oidc.GetIpPoolV6Ok(); ok {
		d.Set("ip_pool_v6", *v)
	}

	d.Set("user_scripts", oidc.GetUserScripts())
	d.Set("dns_servers", oidc.GetDnsServers())
	d.Set("dns_search_domains", oidc.GetDnsSearchDomains())
	d.Set("block_local_dns_requests", oidc.GetBlockLocalDnsRequests())
	if v, ok := oidc.GetClaimMappingsOk(); ok {
		if err := d.Set("claim_mappings", flattenIdentityProviderClaimsMappning(v)); err != nil {
			return err
		}
	}

	if v, ok := oidc.GetOnDemandClaimMappingsOk(); ok {
		if err := d.Set("on_demand_claim_mappings", flattenIdentityProviderOnDemandClaimsMappning(v)); err != nil {
			return err
		}
	}

	// oidc attributes
	if v, ok := oidc.GetIssuerOk(); ok {
		d.Set("issuer", v)
	}
	if v, ok := oidc.GetAudienceOk(); ok {
		d.Set("audience", *v)
	}
	if v, ok := oidc.GetScopeOk(); ok {
		d.Set("scope", *v)
	}
	if v, ok := oidc.GetGoogleOk(); ok {
		g := make(map[string]interface{})
		if val, ok := v.GetEnabledOk(); ok {
			g["enabled"] = *val
		}
		if val, ok := v.GetClientSecretOk(); ok {
			g["client_secret"] = *val
		} else {
			g["client_secret"] = d.Get("google.0.client_secret")
		}
		if val, ok := v.GetRefreshTokenOk(); ok {
			g["refresh_token"] = *val
		}

		d.Set("google", []interface{}{g})
	}
	return nil
}

func resourceAppgateOidcProviderRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating radius identity provider id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.OidcIdentityProvidersApi
	ctx := context.TODO()
	request := api.IdentityProvidersIdGet(ctx, d.Id())
	originalOidcProvider, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to read LDAP Identity provider, %w", err)
	}
	// base attributes
	if d.HasChange("name") {
		originalOidcProvider.SetName(d.Get("name").(string))
	}

	if d.HasChange("notes") {
		originalOidcProvider.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("tags") {
		originalOidcProvider.SetTags(schemaExtractTags(d))
	}

	// identity provider attributes
	if d.HasChange("admin_provider") {
		originalOidcProvider.SetAdminProvider(d.Get("admin_provider").(bool))
	}
	if d.HasChange("device_limit_per_user") {
		originalOidcProvider.SetDeviceLimitPerUser(int32(d.Get("device_limit_per_user").(int)))
	}
	if d.HasChange("on_boarding_two_factor") {
		_, v := d.GetChange("on_boarding_two_factor")
		onboarding, err := readOnBoardingTwoFactorFromConfig(v.([]interface{}))
		if err != nil {
			return err
		}
		originalOidcProvider.SetOnBoarding2FA(onboarding)
	}

	if d.HasChange("inactivity_timeout_minutes") {
		originalOidcProvider.SetInactivityTimeoutMinutes(int32(d.Get("inactivity_timeout_minutes").(int)))
	}
	if d.HasChange("network_inactivity_timeout_enabled") {
		originalOidcProvider.SetNetworkInactivityTimeoutEnabled(d.Get("network_inactivity_timeout_enabled").(bool))
	}
	if d.HasChange("ip_pool_v4") {
		originalOidcProvider.SetIpPoolV4(d.Get("ip_pool_v4").(string))
	}
	if d.HasChange("ip_pool_v6") {
		originalOidcProvider.SetIpPoolV6(d.Get("ip_pool_v6").(string))
	}
	if d.HasChange("user_scripts") {
		_, v := d.GetChange("user_scripts")
		us, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("Failed to read user_scripts %w", err)
		}
		originalOidcProvider.SetUserScripts(us)
	}
	if d.HasChange("dns_servers") {
		_, v := d.GetChange("dns_servers")
		servers, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("Failed to read dns servers %w", err)
		}
		originalOidcProvider.SetDnsServers(servers)
	}
	if d.HasChange("dns_search_domains") {
		_, v := d.GetChange("dns_search_domains")
		servers, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("Failed to read dns search domains %w", err)
		}
		originalOidcProvider.SetDnsSearchDomains(servers)
	}
	if d.HasChange("block_local_dns_requests") {
		originalOidcProvider.SetBlockLocalDnsRequests(d.Get("block_local_dns_requests").(bool))
	}
	if d.HasChange("claim_mappings") {
		_, v := d.GetChange("claim_mappings")
		claims := readIdentityProviderClaimMappingFromConfig(v.([]interface{}))
		originalOidcProvider.SetClaimMappings(claims)
	}
	if d.HasChange("on_demand_claim_mappings") {
		_, v := d.GetChange("on_demand_claim_mappings")
		claims := readIdentityProviderOnDemandClaimMappingFromConfig(v.([]interface{}))
		originalOidcProvider.SetOnDemandClaimMappings(claims)
	}

	// radius provider attributes
	if d.HasChange("issuer") {
		_, v := d.GetChange("issuer")
		originalOidcProvider.SetIssuer(v.(string))
	}
	if d.HasChange("scope") {
		_, v := d.GetChange("scope")
		originalOidcProvider.SetScope(v.(string))
	}
	if d.HasChange("port") {
		_, v := d.GetChange("port")
		originalOidcProvider.SetAudience(v.(string))
	}

	if d.HasChange("google") {
		_, v := d.GetChange("google")
		googles := readOidcProviderGoogleFromConfig(v.([]interface{}))
		originalOidcProvider.SetGoogle(googles[0])
	}

	req := api.IdentityProvidersIdPut(ctx, d.Id())
	req = req.Body(*originalOidcProvider)
	_, _, err = req.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update %s provider %w", identityProviderRadius, prettyPrintAPIError(err))
	}
	return resourceAppgateOidcProviderRuleRead(d, meta)
}

func readOidcProviderGoogleFromConfig(input []interface{}) []openapi.OidcProviderAllOfGoogle {
	googles := make([]openapi.OidcProviderAllOfGoogle, 0)
	for _, raw := range input {
		google := raw.(map[string]interface{})
		g := openapi.NewOidcProviderAllOfGoogle()
		if v, ok := google["enabled"]; ok {
			g.SetEnabled(v.(bool))
		}
		if v, ok := google["client_secret"]; ok {
			g.SetClientSecret(v.(string))
		}
		if v, ok := google["refresh_token"]; ok {
			g.SetRefreshToken(v.(bool))
		}
		googles = append(googles, *g)
	}
	return googles
}
