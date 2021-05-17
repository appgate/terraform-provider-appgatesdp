package appgate

import (
	"context"
	"fmt"
	"log"

	"github.com/appgate/sdp-api-client-go/api/v15/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateConnectorProvider() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateConnectorProviderRuleCreate,
		Read:   resourceAppgateConnectorProviderRuleRead,
		Update: resourceAppgateConnectorProviderRuleUpdate,
		Delete: resourceAppgateConnectorProviderRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: func() map[string]*schema.Schema {
			s := mergeSchemaMaps(baseEntitySchema(), identityProviderIPPoolSchema(), identityProviderClaimsSchema())
			s["name"] = &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  builtinProviderConnector,
			}
			s["type"] = &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  identityProviderConnector,
			}

			return s
		}(),
	}
}

func resourceAppgateConnectorProviderRuleDelete(d *schema.ResourceData, meta interface{}) error {
	// We can't delete the builtin connector identity provider, but we can remove it from the terraform state file.
	d.SetId("")
	return nil
}

func resourceAppgateConnectorProviderRuleCreate(d *schema.ResourceData, meta interface{}) error {
	// we aren'ẗ allowed to create new additional local identity providers, but we can update existing
	// with terraform import.
	token := meta.(*Client).Token
	api := meta.(*Client).API.ConnectorIdentityProvidersApi
	ctx := context.TODO()
	connectorIP, err := getBuiltinConnectorProviderUUID(ctx, *api, token)
	if err != nil {
		return err
	}

	d.SetId(connectorIP.GetId())

	return resourceAppgateConnectorProviderRuleUpdate(d, meta)
}

func getBuiltinConnectorProviderUUID(ctx context.Context, api openapi.ConnectorIdentityProvidersApiService, token string) (*openapi.ConnectorProvider, error) {
	var connectorIP *openapi.ConnectorProvider
	request := api.IdentityProvidersGet(ctx)

	provider, _, err := request.Query(builtinProviderConnector).OrderBy("name").Range_("0-25").Authorization(token).Execute()
	if err != nil {
		return connectorIP, err
	}
	for _, s := range provider.GetData() {
		if s.GetName() == builtinProviderConnector {
			return &s, nil
		}
	}
	return connectorIP, fmt.Errorf("Could not find builtin connector identity provider")
}

func resourceAppgateConnectorProviderRuleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading connectorIP identity provider")

	token := meta.(*Client).Token
	api := meta.(*Client).API.ConnectorIdentityProvidersApi
	ctx := context.TODO()
	connectorIP, err := getBuiltinConnectorProviderUUID(ctx, *api, token)
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Failed to read Connector Identity provider, %+v", err)
	}
	d.SetId(connectorIP.GetId())

	d.Set("type", identityProviderConnector)
	// base attributes
	d.Set("name", connectorIP.Name)
	d.Set("notes", connectorIP.Notes)
	d.Set("tags", connectorIP.Tags)

	// identity provider attributes
	if v, ok := connectorIP.GetIpPoolV4Ok(); ok {
		d.Set("ip_pool_v4", *v)
	}
	if v, ok := connectorIP.GetIpPoolV6Ok(); ok {
		d.Set("ip_pool_v6", v)
	}

	if v, ok := connectorIP.GetClaimMappingsOk(); ok {
		if err := d.Set("claim_mappings", flattenIdentityProviderClaimsMappning(*v)); err != nil {
			return err
		}
	}
	if v, ok := connectorIP.GetOnDemandClaimMappingsOk(); ok {
		d.Set("on_demand_claim_mappings", flattenIdentityProviderOnDemandClaimsMappning(*v))
	}

	return nil
}

func resourceAppgateConnectorProviderRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating connectorIP identity provider id: %+v", d.Id())
	token := meta.(*Client).Token
	api := meta.(*Client).API.ConnectorIdentityProvidersApi
	ctx := context.TODO()
	request := api.IdentityProvidersIdGet(ctx, d.Id())
	originalConnectorProvider, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to read Connector Identity provider, %+v", err)
	}
	// base attributes
	if d.HasChange("name") {
		originalConnectorProvider.SetName(d.Get("name").(string))
	}

	if d.HasChange("notes") {
		originalConnectorProvider.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("tags") {
		originalConnectorProvider.SetTags(schemaExtractTags(d))
	}

	// identity provider attributes
	if d.HasChange("ip_pool_v4") {
		originalConnectorProvider.SetIpPoolV4(d.Get("ip_pool_v4").(string))
	}
	if d.HasChange("ip_pool_v6") {
		originalConnectorProvider.SetIpPoolV6(d.Get("ip_pool_v6").(string))
	}
	if d.HasChange("claim_mappings") {
		_, v := d.GetChange("claim_mappings")
		claims := readIdentityProviderClaimMappingFromConfig(v.([]interface{}))
		originalConnectorProvider.SetClaimMappings(claims)
	}
	if d.HasChange("on_demand_claim_mappings") {
		_, v := d.GetChange("on_demand_claim_mappings")
		claims := readIdentityProviderOnDemandClaimMappingFromConfig(v.([]interface{}))
		originalConnectorProvider.SetOnDemandClaimMappings(claims)
	}

	req := api.IdentityProvidersIdPut(ctx, d.Id())
	req = req.IdentityProvider(originalConnectorProvider)
	_, _, err = req.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update %s provider %+v", identityProviderConnector, prettyPrintAPIError(err))
	}
	return resourceAppgateConnectorProviderRuleRead(d, meta)
}
