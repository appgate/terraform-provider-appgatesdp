package appgate

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/appgate/sdp-api-client-go/api/v21/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateLdapCertificateProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateLdapCertificateProviderRuleCreate,
		Read:   resourceAppgateLdapCertificateProviderRuleRead,
		Update: resourceAppgateLdapCertificateProviderRuleUpdate,
		Delete: identityProviderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: func() map[string]*schema.Schema {
			s := ldapProviderSchema()
			s["type"].Default = identityProviderLdapCertificate

			s["ca_certificates"] = &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			}
			s["certificate_user_attribute"] = &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			}
			s["certificate_attribute"] = &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			}
			s["skip_x509_external_checks"] = &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			}
			// LDAP Certificate does not use password_warning
			// wrong in the openapi spec.
			delete(s, "password_warning")

			return s
		}(),
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    identityProviderResourcev0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceIdentityProvidereUpgradeV0,
				Version: 0,
			},
		},
	}
}

func resourceAppgateLdapCertificateProviderRuleCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating LdapCertificateProvider: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.LdapCertificateIdentityProvidersApi
	ctx := context.TODO()
	currentVersion := meta.(*Client).ApplianceVersion
	provider := &openapi.ConfigurableIdentityProvider{}
	provider.Type = identityProviderLdapCertificate
	provider, err = readProviderFromConfig(d, *provider, currentVersion)
	if err != nil {
		return fmt.Errorf("Failed to read and create basic identity provider for %s %w", identityProviderLdapCertificate, err)
	}

	args := openapi.LdapCertificateProvider{}

	if currentVersion.LessThan(Appliance55Version) {
		args.DeviceLimitPerUser = nil
	}

	args.SetType(provider.GetType())
	args.SetId(provider.GetId())
	args.SetName(provider.GetName())
	args.SetNotes(provider.GetNotes())
	args.SetTags(provider.GetTags())

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
	if v, ok := d.GetOk("user_filter"); ok {
		args.SetUserFilter(v.(string))
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
	// ldadp certificcate attributes
	if v, ok := d.GetOk("ca_certificates"); ok {
		certificates, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return err
		}
		args.SetCaCertificates(certificates)
	}

	if v, ok := d.GetOk("certificate_user_attribute"); ok {
		args.SetCertificateUserAttribute(v.(string))
	}
	if v, ok := d.GetOk("certificate_attribute"); ok {
		args.SetCertificateAttribute(v.(string))
	}
	if v, ok := d.GetOk("skip_x509_external_checks"); ok {
		args.SetSkipX509ExternalChecks(v.(bool))
	}
	request := api.IdentityProvidersPost(ctx)
	p, _, err := request.Body(args).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create %s provider %w", identityProviderLdapCertificate, prettyPrintAPIError(err))
	}
	d.SetId(p.GetId())
	return resourceAppgateLdapCertificateProviderRuleRead(d, meta)
}

func resourceAppgateLdapCertificateProviderRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading ldap identity provider id: %+v", d.Id())

	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.LdapCertificateIdentityProvidersApi
	ctx := context.TODO()
	request := api.IdentityProvidersIdGet(ctx, d.Id())
	ldap, res, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		if res != nil && res.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("Failed to read LDAP Identity provider, %w", err)
	}
	d.Set("type", identityProviderLdapCertificate)
	// base attributes
	d.Set("name", ldap.GetName())
	d.Set("notes", ldap.GetNotes())
	d.Set("tags", ldap.GetTags())

	// identity provider attributes
	d.Set("admin_provider", ldap.GetAdminProvider())
	if v, ok := ldap.GetDeviceLimitPerUserOk(); ok {
		d.Set("device_limit_per_user", *v)
	}
	if v, ok := ldap.GetOnBoarding2FAOk(); ok {
		if err := d.Set("on_boarding_two_factor", flattenIdentityProviderOnboarding2fa(*v)); err != nil {
			return err
		}
	}

	d.Set("inactivity_timeout_minutes", ldap.GetInactivityTimeoutMinutes())
	d.Set("network_inactivity_timeout_enabled", ldap.GetNetworkInactivityTimeoutEnabled())
	if v, ok := ldap.GetIpPoolV4Ok(); ok {
		d.Set("ip_pool_v4", *v)
	}
	if v, ok := ldap.GetIpPoolV6Ok(); ok {
		d.Set("ip_pool_v6", *v)
	}

	d.Set("user_scripts", ldap.GetUserScripts())
	d.Set("dns_servers", ldap.GetDnsServers())
	d.Set("dns_search_domains", ldap.GetDnsSearchDomains())
	d.Set("block_local_dns_requests", ldap.GetBlockLocalDnsRequests())
	if v, ok := ldap.GetClaimMappingsOk(); ok {
		if err := d.Set("claim_mappings", flattenIdentityProviderClaimsMappning(v)); err != nil {
			return err
		}
	}
	if v, ok := ldap.GetOnDemandClaimMappingsOk(); ok {
		if err := d.Set("on_demand_claim_mappings", flattenIdentityProviderOnDemandClaimsMappning(v)); err != nil {
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
	d.Set("user_filter", ldap.GetUserFilter())
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

	d.Set("ca_certificates", ldap.GetCaCertificates())
	d.Set("certificate_user_attribute", ldap.GetCertificateUserAttribute())
	d.Set("certificate_attribute", ldap.GetCertificateAttribute())
	d.Set("skip_x509_external_checks", ldap.GetSkipX509ExternalChecks())
	return nil
}

func resourceAppgateLdapCertificateProviderRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating ldap identity provider id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.LdapCertificateIdentityProvidersApi
	ctx := context.TODO()
	request := api.IdentityProvidersIdGet(ctx, d.Id())
	originalLdapCertificateProvider, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to read LDAP Identity provider, %w", err)
	}
	// base attributes
	if d.HasChange("name") {
		originalLdapCertificateProvider.SetName(d.Get("name").(string))
	}

	if d.HasChange("notes") {
		originalLdapCertificateProvider.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("tags") {
		originalLdapCertificateProvider.SetTags(schemaExtractTags(d))
	}

	// identity provider attributes
	if d.HasChange("admin_provider") {
		originalLdapCertificateProvider.SetAdminProvider(d.Get("admin_provider").(bool))
	}
	if d.HasChange("device_limit_per_user") {
		originalLdapCertificateProvider.SetDeviceLimitPerUser(int32(d.Get("device_limit_per_user").(int)))
	}
	if d.HasChange("on_boarding_two_factor") {
		_, v := d.GetChange("on_boarding_two_factor")
		onboarding, err := readOnBoardingTwoFactorFromConfig(v.([]interface{}))
		if err != nil {
			return err
		}
		originalLdapCertificateProvider.SetOnBoarding2FA(onboarding)
	}

	if d.HasChange("inactivity_timeout_minutes") {
		originalLdapCertificateProvider.SetInactivityTimeoutMinutes(int32(d.Get("inactivity_timeout_minutes").(int)))
	}
	if d.HasChange("network_inactivity_timeout_enabled") {
		originalLdapCertificateProvider.SetNetworkInactivityTimeoutEnabled(d.Get("network_inactivity_timeout_enabled").(bool))
	}
	if d.HasChange("ip_pool_v4") {
		originalLdapCertificateProvider.SetIpPoolV4(d.Get("ip_pool_v4").(string))
	}
	if d.HasChange("ip_pool_v6") {
		originalLdapCertificateProvider.SetIpPoolV6(d.Get("ip_pool_v6").(string))
	}
	if d.HasChange("user_scripts") {
		_, v := d.GetChange("user_scripts")
		us, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("Failed to read user_scripts %w", err)
		}
		originalLdapCertificateProvider.SetUserScripts(us)
	}
	if d.HasChange("dns_servers") {
		_, v := d.GetChange("dns_servers")
		servers, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("Failed to read dns servers %w", err)
		}
		originalLdapCertificateProvider.SetDnsServers(servers)
	}
	if d.HasChange("dns_search_domains") {
		_, v := d.GetChange("dns_search_domains")
		servers, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("Failed to read dns search domains %w", err)
		}
		originalLdapCertificateProvider.SetDnsSearchDomains(servers)
	}
	if d.HasChange("block_local_dns_requests") {
		originalLdapCertificateProvider.SetBlockLocalDnsRequests(d.Get("block_local_dns_requests").(bool))
	}
	if d.HasChange("claim_mappings") {
		_, v := d.GetChange("claim_mappings")
		claims := readIdentityProviderClaimMappingFromConfig(v.([]interface{}))
		originalLdapCertificateProvider.SetClaimMappings(claims)
	}
	if d.HasChange("on_demand_claim_mappings") {
		_, v := d.GetChange("on_demand_claim_mappings")
		claims := readIdentityProviderOnDemandClaimMappingFromConfig(v.([]interface{}))
		originalLdapCertificateProvider.SetOnDemandClaimMappings(claims)
	}

	// ldap provider attributes
	if d.HasChange("hostnames") {
		_, v := d.GetChange("hostnames")
		hostnames, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return err
		}
		originalLdapCertificateProvider.SetHostnames(hostnames)
	}
	if d.HasChange("port") {
		originalLdapCertificateProvider.SetPort(int32(d.Get("port").(int)))
	}
	if d.HasChange("ssl_enabled") {
		originalLdapCertificateProvider.SetSslEnabled(d.Get("ssl_enabled").(bool))
	}
	if d.HasChange("admin_distinguished_name") {
		originalLdapCertificateProvider.SetAdminDistinguishedName(d.Get("admin_distinguished_name").(string))
	}
	if d.HasChange("admin_password") {
		originalLdapCertificateProvider.SetAdminPassword(d.Get("admin_password").(string))
	}
	if d.HasChange("base_dn") {
		originalLdapCertificateProvider.SetBaseDn(d.Get("base_dn").(string))
	}
	if d.HasChange("object_class") {
		originalLdapCertificateProvider.SetObjectClass(d.Get("object_class").(string))
	}
	if d.HasChange("user_filter") {
		originalLdapCertificateProvider.SetUserFilter(d.Get("user_filter").(string))
	}
	if d.HasChange("username_attribute") {
		originalLdapCertificateProvider.SetUsernameAttribute(d.Get("username_attribute").(string))
	}
	if d.HasChange("membership_filter") {
		originalLdapCertificateProvider.SetMembershipFilter(d.Get("membership_filter").(string))
	}
	if d.HasChange("membership_base_dn") {
		originalLdapCertificateProvider.SetMembershipBaseDn(d.Get("membership_base_dn").(string))
	}
	if d.HasChange("password_warning") {
		_, v := d.GetChange("password_warning")
		pw := readLdapPasswordWarningFromConfig(v.([]interface{}))
		originalLdapCertificateProvider.SetPasswordWarning(pw)
	}

	// ldadp certificcate attributes
	if d.HasChange("ca_certificates") {
		_, v := d.GetChange("ca_certificates")
		certificates, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return err
		}
		originalLdapCertificateProvider.SetCaCertificates(certificates)
	}

	if d.HasChange("certificate_user_attribute") {
		originalLdapCertificateProvider.SetCertificateUserAttribute(d.Get("certificate_user_attribute").(string))
	}
	if d.HasChange("certificate_attribute") {
		originalLdapCertificateProvider.SetCertificateAttribute(d.Get("certificate_attribute").(string))
	}
	if d.HasChange("skip_x509_external_checks") {
		originalLdapCertificateProvider.SetSkipX509ExternalChecks(d.Get("skip_x509_external_checks").(bool))
	}
	req := api.IdentityProvidersIdPut(ctx, d.Id())
	req = req.Body(*originalLdapCertificateProvider)
	_, _, err = req.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update %s provider %w", identityProviderLdapCertificate, prettyPrintAPIError(err))
	}
	return resourceAppgateLdapCertificateProviderRuleRead(d, meta)
}
