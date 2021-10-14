package appgate

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v16/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			s := ldapProviderSchema()
			s["type"].Default = identityProviderLdap

			return s
		}(),
	}
}

func resourceAppgateLdapProviderRuleCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating LdapProvider: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.LdapIdentityProvidersApi
	ctx := context.TODO()
	provider := &openapi.IdentityProvider{}
	provider.Type = identityProviderLdap
	provider, err = readProviderFromConfig(d, *provider)
	if err != nil {
		return fmt.Errorf("Failed to read and create basic identity provider for %s %s", identityProviderLdap, err)
	}

	args := openapi.NewLdapProviderWithDefaults()
	args.SetType(provider.GetType())
	args.SetId(provider.GetId())
	args.SetName(provider.GetName())
	args.SetNotes(provider.GetNotes())
	args.SetTags(provider.GetTags())

	if provider.AdminProvider != nil {
		args.SetAdminProvider(*provider.AdminProvider)
	}
	if provider.OnBoarding2FA != nil {
		args.SetOnBoarding2FA(*provider.OnBoarding2FA)
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
		pw := readLdapPasswordWarningFromConfig(v.([]interface{}))
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

	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.LdapIdentityProvidersApi
	ctx := context.TODO()
	request := api.IdentityProvidersIdGet(ctx, d.Id())
	ldap, res, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		if res.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("Failed to read LDAP Identity provider, %+v", err)
	}
	d.Set("type", identityProviderLdap)
	// base attributes
	d.Set("name", ldap.Name)
	d.Set("notes", ldap.Notes)
	d.Set("tags", ldap.Tags)

	// identity provider attributes
	if v, ok := ldap.GetOnBoarding2FAOk(); ok {
		if err := d.Set("on_boarding_two_factor", flattenIdentityProviderOnboarding2fa(*v)); err != nil {
			return err
		}
	}

	d.Set("inactivity_timeout_minutes", ldap.GetInactivityTimeoutMinutes())
	if v, ok := ldap.GetIpPoolV4Ok(); ok {
		d.Set("ip_pool_v4", *v)
	}
	if v, ok := ldap.GetIpPoolV6Ok(); ok {
		d.Set("ip_pool_v6", *v)
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
		if err := d.Set("on_demand_claim_mappings", flattenIdentityProviderOnDemandClaimsMappning(*v)); err != nil {
			return err
		}
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
		d.Set("base_dn", *v)
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

func flattenLdapPasswordWarning(pw openapi.LdapProviderAllOfPasswordWarning) []interface{} {
	o := make(map[string]interface{})
	if v, ok := pw.GetEnabledOk(); ok {
		o["enabled"] = *v
	}
	if v, ok := pw.GetThresholdDaysOk(); ok {
		o["threshold_days"] = int(*v)
	}
	if v, ok := pw.GetMessageOk(); ok {
		o["message"] = v
	}
	return []interface{}{o}
}

func readLdapPasswordWarningFromConfig(input []interface{}) openapi.LdapProviderAllOfPasswordWarning {
	pw := openapi.LdapProviderAllOfPasswordWarning{}
	for _, r := range input {
		raw := r.(map[string]interface{})
		if v, ok := raw["enabled"]; ok {
			pw.SetEnabled(v.(bool))
		}
		if v, ok := raw["threshold_days"]; ok {
			pw.SetThresholdDays(int32(v.(int)))
		}
		if v, ok := raw["message"]; ok {
			pw.SetMessage(v.(string))
		}
	}
	return pw
}

func resourceAppgateLdapProviderRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating ldap identity provider id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.LdapIdentityProvidersApi
	ctx := context.TODO()
	request := api.IdentityProvidersIdGet(ctx, d.Id())
	originalLdapProvider, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to read LDAP Identity provider, %+v", err)
	}
	// base attributes
	if d.HasChange("name") {
		originalLdapProvider.SetName(d.Get("name").(string))
	}

	if d.HasChange("notes") {
		originalLdapProvider.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("tags") {
		originalLdapProvider.SetTags(schemaExtractTags(d))
	}

	// identity provider attributes
	if d.HasChange("admin_provider") {
		originalLdapProvider.SetAdminProvider(d.Get("admin_provider").(bool))
	}
	if d.HasChange("on_boarding_two_factor") {
		_, v := d.GetChange("on_boarding_two_factor")
		onboarding := readOnBoardingTwoFactorFromConfig(v.([]interface{}))
		originalLdapProvider.SetOnBoarding2FA(onboarding)
	}

	if d.HasChange("inactivity_timeout_minutes") {
		originalLdapProvider.SetInactivityTimeoutMinutes(int32(d.Get("inactivity_timeout_minutes").(int)))
	}
	if d.HasChange("ip_pool_v4") {
		originalLdapProvider.SetIpPoolV4(d.Get("ip_pool_v4").(string))
	}
	if d.HasChange("ip_pool_v6") {
		originalLdapProvider.SetIpPoolV6(d.Get("ip_pool_v6").(string))
	}
	if d.HasChange("dns_servers") {
		_, v := d.GetChange("dns_servers")
		servers, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("Failed to read dns servers %s", err)
		}
		originalLdapProvider.SetDnsServers(servers)
	}
	if d.HasChange("dns_search_domains") {
		_, v := d.GetChange("dns_search_domains")
		servers, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("Failed to read dns search domains %s", err)
		}
		originalLdapProvider.SetDnsSearchDomains(servers)
	}
	if d.HasChange("block_local_dns_requests") {
		originalLdapProvider.SetBlockLocalDnsRequests(d.Get("block_local_dns_requests").(bool))
	}
	if d.HasChange("claim_mappings") {
		_, v := d.GetChange("claim_mappings")
		claims := readIdentityProviderClaimMappingFromConfig(v.([]interface{}))
		originalLdapProvider.SetClaimMappings(claims)
	}
	if d.HasChange("on_demand_claim_mappings") {
		_, v := d.GetChange("on_demand_claim_mappings")
		claims := readIdentityProviderOnDemandClaimMappingFromConfig(v.([]interface{}))
		originalLdapProvider.SetOnDemandClaimMappings(claims)
	}

	// ldap provider attributes
	if d.HasChange("hostnames") {
		_, v := d.GetChange("hostnames")
		hostnames, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return err
		}
		originalLdapProvider.SetHostnames(hostnames)
	}
	if d.HasChange("port") {
		originalLdapProvider.SetPort(int32(d.Get("port").(int)))
	}
	if d.HasChange("ssl_enabled") {
		originalLdapProvider.SetSslEnabled(d.Get("ssl_enabled").(bool))
	}
	if d.HasChange("admin_distinguished_name") {
		originalLdapProvider.SetAdminDistinguishedName(d.Get("admin_distinguished_name").(string))
	}
	if d.HasChange("admin_password") {
		originalLdapProvider.SetAdminPassword(d.Get("admin_password").(string))
	}
	if d.HasChange("base_dn") {
		originalLdapProvider.SetBaseDn(d.Get("base_dn").(string))
	}
	if d.HasChange("object_class") {
		originalLdapProvider.SetObjectClass(d.Get("object_class").(string))
	}
	if d.HasChange("username_attribute") {
		originalLdapProvider.SetUsernameAttribute(d.Get("username_attribute").(string))
	}
	if d.HasChange("membership_filter") {
		originalLdapProvider.SetMembershipFilter(d.Get("membership_filter").(string))
	}
	if d.HasChange("membership_base_dn") {
		originalLdapProvider.SetMembershipBaseDn(d.Get("membership_base_dn").(string))
	}
	if d.HasChange("password_warning") {
		_, v := d.GetChange("password_warning")
		pw := readLdapPasswordWarningFromConfig(v.([]interface{}))
		originalLdapProvider.SetPasswordWarning(pw)
	}
	req := api.IdentityProvidersIdPut(ctx, d.Id())
	req = req.IdentityProvider(originalLdapProvider)
	_, _, err = req.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update %s provider %+v", identityProviderLdap, prettyPrintAPIError(err))
	}
	return resourceAppgateLdapProviderRuleRead(d, meta)
}
