package appgate

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v12/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAppgateSamlProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateSamlProviderRuleCreate,
		Read:   resourceAppgateSamlProviderRuleRead,
		Update: resourceAppgateSamlProviderRuleUpdate,
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
			s["type"].Default = identityProviderSaml

			s["redirect_url"] = &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			}
			s["issuer"] = &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			}
			s["audience"] = &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			}
			s["provider_certificate"] = &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			}
			s["decryption_key"] = &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			}
			return s
		}(),
	}
}

func resourceAppgateSamlProviderRuleCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating SamlProvider: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.SamlIdentityProvidersApi
	ctx := context.TODO()
	provider := openapi.NewIdentityProvider(identityProviderSaml)
	provider, err := readProviderFromConfig(d, *provider)
	if err != nil {
		return fmt.Errorf("Failed to read and create basic identity provider for %s %s", identityProviderSaml, err)
	}

	args := openapi.NewSamlProviderWithDefaults()
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

	if v, ok := d.GetOk("redirect_url"); ok {
		args.SetRedirectUrl(v.(string))
	}
	if v, ok := d.GetOk("issuer"); ok {
		args.SetIssuer(v.(string))
	}
	if v, ok := d.GetOk("audience"); ok {
		args.SetAudience(v.(string))
	}
	if v, ok := d.GetOk("provider_certificate"); ok {
		args.SetProviderCertificate(v.(string))
	}
	if v, ok := d.GetOk("decryption_key"); ok {
		args.SetDecryptionKey(v.(string))
	}
	request := api.IdentityProvidersPost(ctx)
	p, _, err := request.IdentityProvider(*args).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create %s provider %+v", identityProviderSaml, prettyPrintAPIError(err))
	}
	d.SetId(p.Id)
	return resourceAppgateSamlProviderRuleRead(d, meta)
}

func resourceAppgateSamlProviderRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading saml identity provider id: %+v", d.Id())

	token := meta.(*Client).Token
	api := meta.(*Client).API.SamlIdentityProvidersApi
	ctx := context.TODO()
	request := api.IdentityProvidersIdGet(ctx, d.Id())
	saml, _, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Failed to read Saml Identity provider, %+v", err)
	}
	d.Set("type", identityProviderSaml)
	// base attributes
	d.Set("name", saml.Name)
	d.Set("notes", saml.Notes)
	d.Set("tags", saml.Tags)

	// identity provider attributes
	d.Set("default", saml.GetDefault())
	d.Set("client_provider", saml.GetClientProvider())
	d.Set("admin_provider", saml.GetAdminProvider())
	if v, ok := saml.GetOnBoarding2FAOk(); ok {
		if err := d.Set("on_boarding_two_factor", flattenIdentityProviderOnboarding2fa(*v)); err != nil {
			return err
		}
	}

	d.Set("inactivity_timeout_minutes", saml.GetInactivityTimeoutMinutes())
	if v, ok := saml.GetIpPoolV4Ok(); ok {
		d.Set("ip_pool_v4", *v)
	}
	if v, ok := saml.GetIpPoolV6Ok(); ok {
		d.Set("ip_pool_v6", v)
	}

	d.Set("dns_servers", saml.GetDnsServers())
	d.Set("dns_search_domains", saml.GetDnsSearchDomains())
	d.Set("block_local_dns_requests", saml.GetBlockLocalDnsRequests())
	if v, ok := saml.GetClaimMappingsOk(); ok {
		if err := d.Set("claim_mappings", flattenIdentityProviderClaimsMappning(*v)); err != nil {
			return err
		}
	}
	if v, ok := saml.GetOnDemandClaimMappingsOk(); ok {
		d.Set("on_demand_claim_mappings", flattenIdentityProviderOnDemandClaimsMappning(*v))
	}
	// saml attributes
	d.Set("redirect_url", saml.GetRedirectUrl())
	d.Set("issuer", saml.GetIssuer())
	d.Set("audience", saml.GetAudience())
	d.Set("provider_certificate", saml.GetProviderCertificate())
	d.Set("decryption_key", saml.GetDecryptionKey())

	return nil
}

func resourceAppgateSamlProviderRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating saml identity provider id: %+v", d.Id())
	token := meta.(*Client).Token
	api := meta.(*Client).API.SamlIdentityProvidersApi
	ctx := context.TODO()
	request := api.IdentityProvidersIdGet(ctx, d.Id())
	originalSamlProvider, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to read Saml Identity provider, %+v", err)
	}
	// base attributes
	if d.HasChange("name") {
		originalSamlProvider.SetName(d.Get("name").(string))
	}

	if d.HasChange("notes") {
		originalSamlProvider.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("tags") {
		originalSamlProvider.SetTags(schemaExtractTags(d))
	}

	// identity provider attributes
	if d.HasChange("default") {
		originalSamlProvider.SetDefault(d.Get("default").(bool))
	}
	if d.HasChange("client_provider") {
		originalSamlProvider.SetClientProvider(d.Get("client_provider").(bool))
	}
	if d.HasChange("admin_provider") {
		originalSamlProvider.SetAdminProvider(d.Get("admin_provider").(bool))
	}
	if d.HasChange("on_boarding_two_factor") {
		_, v := d.GetChange("on_boarding_two_factor")
		onboarding := readOnBoardingTwoFactorFromConfig(v.([]interface{}))
		originalSamlProvider.SetOnBoarding2FA(onboarding)
	}

	if d.HasChange("inactivity_timeout_minutes") {
		originalSamlProvider.SetInactivityTimeoutMinutes(int32(d.Get("inactivity_timeout_minutes").(int)))
	}
	if d.HasChange("ip_pool_v4") {
		originalSamlProvider.SetIpPoolV4(d.Get("ip_pool_v4").(string))
	}
	if d.HasChange("ip_pool_v6") {
		originalSamlProvider.SetIpPoolV6(d.Get("ip_pool_v6").(string))
	}
	if d.HasChange("dns_servers") {
		_, v := d.GetChange("dns_servers")
		servers, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("Failed to read dns servers %s", err)
		}
		originalSamlProvider.SetDnsServers(servers)
	}
	if d.HasChange("dns_search_domains") {
		_, v := d.GetChange("dns_search_domains")
		servers, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("Failed to read dns search domains %s", err)
		}
		originalSamlProvider.SetDnsSearchDomains(servers)
	}
	if d.HasChange("block_local_dns_requests") {
		originalSamlProvider.SetBlockLocalDnsRequests(d.Get("block_local_dns_requests").(bool))
	}
	if d.HasChange("claim_mappings") {
		_, v := d.GetChange("claim_mappings")
		claims := readIdentityProviderClaimMappingFromConfig(v.([]interface{}))
		originalSamlProvider.SetClaimMappings(claims)
	}
	if d.HasChange("on_demand_claim_mappings") {
		_, v := d.GetChange("on_demand_claim_mappings")
		claims := readIdentityProviderOnDemandClaimMappingFromConfig(v.([]interface{}))
		originalSamlProvider.SetOnDemandClaimMappings(claims)
	}

	// saml provider attributes
	if d.HasChange("redirect_url") {
		originalSamlProvider.SetRedirectUrl(d.Get("redirect_url").(string))
	}
	if d.HasChange("issuer") {
		originalSamlProvider.SetIssuer(d.Get("issuer").(string))
	}
	if d.HasChange("audience") {
		originalSamlProvider.SetAudience(d.Get("audience").(string))
	}
	if d.HasChange("provider_certificate") {
		originalSamlProvider.SetProviderCertificate(d.Get("provider_certificate").(string))
	}
	if d.HasChange("decryption_key") {
		originalSamlProvider.SetDecryptionKey(d.Get("decryption_key").(string))
	}

	req := api.IdentityProvidersIdPut(ctx, d.Id())
	req = req.IdentityProvider(originalSamlProvider)
	_, _, err = req.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update %s provider %+v", identityProviderSaml, prettyPrintAPIError(err))
	}
	return resourceAppgateSamlProviderRuleRead(d, meta)
}
