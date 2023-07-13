package appgate

import (
	"context"
	"fmt"
	"log"

	"github.com/appgate/sdp-api-client-go/api/v19/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateLocalDatabaseProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateLocalDatabaseProviderRuleCreate,
		Read:   resourceAppgateLocalDatabaseProviderRuleRead,
		Update: resourceAppgateLocalDatabaseProviderRuleUpdate,
		Delete: resourceAppgateLocalDatabaseProviderRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: func() map[string]*schema.Schema {
			s := identityProviderSchema()
			s["name"] = &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  builtinProviderLocal,
			}
			s["type"].Default = identityProviderLocalDatabase

			s["user_lockout_threshold"] = &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			}
			s["min_password_length"] = &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			}

			return s
		}(),
	}
}

func resourceAppgateLocalDatabaseProviderRuleDelete(d *schema.ResourceData, meta interface{}) error {
	// We can't delete the builtin local database identity provider, but we can remove it from the terraform state file.
	d.SetId("")
	return nil
}

func resourceAppgateLocalDatabaseProviderRuleCreate(d *schema.ResourceData, meta interface{}) error {
	// we aren'ẗ allowed to create new additional local identity providers, but we can update existing
	// with terraform import.
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.LocalDatabaseIdentityProvidersApi
	ctx := context.TODO()
	localDatabase, err := getBuiltinLocalDatabaseProviderUUID(ctx, *api, token)
	if err != nil {
		return err
	}

	d.SetId(localDatabase.GetId())

	return resourceAppgateLocalDatabaseProviderRuleUpdate(d, meta)
}

func getBuiltinLocalDatabaseProviderUUID(ctx context.Context, api openapi.LocalDatabaseIdentityProvidersApiService, token string) (*openapi.LocalDatabaseProvider, error) {
	var localDatabase *openapi.LocalDatabaseProvider
	request := api.IdentityProvidersGet(ctx)

	provider, _, err := request.Query(builtinProviderLocal).OrderBy("name").Range_("0-25").Authorization(token).Execute()
	if err != nil {
		return localDatabase, err
	}
	for _, s := range provider.GetData() {
		if s.GetName() == builtinProviderLocal {
			return &s, nil
		}
	}
	return localDatabase, fmt.Errorf("Could not find builtin local database identity provider")
}

func resourceAppgateLocalDatabaseProviderRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading localDatabase identity provider")

	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.LocalDatabaseIdentityProvidersApi
	ctx := context.TODO()
	currentVersion := meta.(*Client).ApplianceVersion
	localDatabase, err := getBuiltinLocalDatabaseProviderUUID(ctx, *api, token)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Failed to read LocalDatabase Identity provider, %w", err)
	}
	d.SetId(localDatabase.GetId())

	d.Set("type", identityProviderLocalDatabase)
	// base attributes
	d.Set("name", localDatabase.Name)
	d.Set("notes", localDatabase.Notes)
	d.Set("tags", localDatabase.Tags)

	// identity provider attributes
	d.Set("admin_provider", localDatabase.GetAdminProvider())
	if v, ok := localDatabase.GetOnBoarding2FAOk(); ok {
		if err := d.Set("on_boarding_two_factor", flattenIdentityProviderOnboarding2fa(*v, currentVersion)); err != nil {
			return err
		}
	}

	d.Set("inactivity_timeout_minutes", localDatabase.GetInactivityTimeoutMinutes())
	d.Set("network_inactivity_timeout_enabled", localDatabase.GetNetworkInactivityTimeoutEnabled())
	if v, ok := localDatabase.GetIpPoolV4Ok(); ok {
		d.Set("ip_pool_v4", *v)
	}
	if v, ok := localDatabase.GetIpPoolV6Ok(); ok {
		d.Set("ip_pool_v6", v)
	}

	d.Set("user_scripts", localDatabase.GetUserScripts())
	d.Set("dns_servers", localDatabase.GetDnsServers())
	d.Set("dns_search_domains", localDatabase.GetDnsSearchDomains())
	d.Set("block_local_dns_requests", localDatabase.GetBlockLocalDnsRequests())
	if v, ok := localDatabase.GetClaimMappingsOk(); ok {
		if err := d.Set("claim_mappings", flattenIdentityProviderClaimsMappning(v)); err != nil {
			return err
		}
	}
	if v, ok := localDatabase.GetOnDemandClaimMappingsOk(); ok {
		d.Set("on_demand_claim_mappings", flattenIdentityProviderOnDemandClaimsMappning(v))
	}
	// localDatabase attributes
	d.Set("user_lockout_threshold", localDatabase.GetUserLockoutThreshold())
	d.Set("min_password_length", localDatabase.GetMinPasswordLength())

	return nil
}

func resourceAppgateLocalDatabaseProviderRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating localDatabase identity provider id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.LocalDatabaseIdentityProvidersApi
	ctx := context.TODO()
	currentVersion := meta.(*Client).ApplianceVersion
	request := api.IdentityProvidersIdGet(ctx, d.Id())
	originalLocalDatabaseProvider, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to read LocalDatabase Identity provider, %w", err)
	}
	// base attributes
	if d.HasChange("name") {
		originalLocalDatabaseProvider.SetName(d.Get("name").(string))
	}

	if d.HasChange("notes") {
		originalLocalDatabaseProvider.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("tags") {
		originalLocalDatabaseProvider.SetTags(schemaExtractTags(d))
	}

	// identity provider attributes

	if d.HasChange("admin_provider") {
		originalLocalDatabaseProvider.SetAdminProvider(d.Get("admin_provider").(bool))
	}
	if d.HasChange("on_boarding_two_factor") {
		_, v := d.GetChange("on_boarding_two_factor")
		onboarding, err := readOnBoardingTwoFactorFromConfig(v.([]interface{}), currentVersion)
		if err != nil {
			return err
		}
		originalLocalDatabaseProvider.SetOnBoarding2FA(onboarding)
	}

	if d.HasChange("inactivity_timeout_minutes") {
		originalLocalDatabaseProvider.SetInactivityTimeoutMinutes(int32(d.Get("inactivity_timeout_minutes").(int)))
	}
	if d.HasChange("network_inactivity_timeout_enabled") {
		originalLocalDatabaseProvider.SetNetworkInactivityTimeoutEnabled(d.Get("network_inactivity_timeout_enabled").(bool))
	}
	if d.HasChange("ip_pool_v4") {
		originalLocalDatabaseProvider.SetIpPoolV4(d.Get("ip_pool_v4").(string))
	}
	if d.HasChange("ip_pool_v6") {
		originalLocalDatabaseProvider.SetIpPoolV6(d.Get("ip_pool_v6").(string))
	}
	if d.HasChange("user_scripts") {
		_, v := d.GetChange("user_scripts")
		us, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("Failed to read user_scripts %w", err)
		}
		originalLocalDatabaseProvider.SetUserScripts(us)
	}
	if d.HasChange("dns_servers") {
		_, v := d.GetChange("dns_servers")
		servers, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("Failed to read dns servers %w", err)
		}
		originalLocalDatabaseProvider.SetDnsServers(servers)
	}
	if d.HasChange("dns_search_domains") {
		_, v := d.GetChange("dns_search_domains")
		servers, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("Failed to read dns search domains %w", err)
		}
		originalLocalDatabaseProvider.SetDnsSearchDomains(servers)
	}
	if d.HasChange("block_local_dns_requests") {
		originalLocalDatabaseProvider.SetBlockLocalDnsRequests(d.Get("block_local_dns_requests").(bool))
	}
	if d.HasChange("claim_mappings") {
		_, v := d.GetChange("claim_mappings")
		claims := readIdentityProviderClaimMappingFromConfig(v.([]interface{}))
		originalLocalDatabaseProvider.SetClaimMappings(claims)
	}
	if d.HasChange("on_demand_claim_mappings") {
		_, v := d.GetChange("on_demand_claim_mappings")
		claims := readIdentityProviderOnDemandClaimMappingFromConfig(v.([]interface{}))
		originalLocalDatabaseProvider.SetOnDemandClaimMappings(claims)
	}

	// localDatabase provider attributes
	if d.HasChange("user_lockout_threshold") {
		originalLocalDatabaseProvider.SetUserLockoutThreshold(int32(d.Get("user_lockout_threshold").(int)))
	}
	if d.HasChange("min_password_length") {
		originalLocalDatabaseProvider.SetMinPasswordLength(int32(d.Get("min_password_length").(int)))
	}

	req := api.IdentityProvidersIdPut(ctx, d.Id())
	req = req.Body(*originalLocalDatabaseProvider)
	_, _, err = req.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update %s provider %w", identityProviderLocalDatabase, prettyPrintAPIError(err))
	}
	return resourceAppgateLocalDatabaseProviderRuleRead(d, meta)
}
