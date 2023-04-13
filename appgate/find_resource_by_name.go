// Code generated by go generate; DO NOT EDIT.
package appgate

import (
	"context"
	"github.com/appgate/sdp-api-client-go/api/v18/openapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func findEntitlementByUUID(ctx context.Context, api *openapi.EntitlementsApiService, id, token string) (*openapi.Entitlement, diag.Diagnostics) {
	log.Printf("[DEBUG] Data source Entitlement get by UUID %s", id)
	resource, _, err := api.EntitlementsIdGet(ctx, id).Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return resource, nil
}

func findEntitlementByName(ctx context.Context, api *openapi.EntitlementsApiService, name, token string) (*openapi.Entitlement, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Data source Entitlement get by name %s", name)

	resource, _, err := api.EntitlementsGet(ctx).Query(name).OrderBy("name").Range_("0-10").Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	for _, r := range resource.GetData() {
		if r.GetName() == name {
			return &r, nil
		}
	}
	if len(resource.GetData()) > 1 {
		return nil, AppendErrorf(diags, "multiple Entitlement matched; use additional constraints to reduce matches to a single Entitlement")
	}
	return nil, AppendErrorf(diags, "Failed to find Entitlement %s", name)
}

func ResolveEntitlementFromResourceData(ctx context.Context, d *schema.ResourceData, api *openapi.EntitlementsApiService, token string) (*openapi.Entitlement, diag.Diagnostics) {
	var diags diag.Diagnostics
	resourceID, iok := d.GetOk("entitlement_id")
	resourceName, nok := d.GetOk("entitlement_name")

	if !iok && !nok {
		return nil, AppendErrorf(diags, "please provide one of entitlement_id or entitlement_name attributes")
	}

	if iok {
		return findEntitlementByUUID(ctx, api, resourceID.(string), token)
	}
	return findEntitlementByName(ctx, api, resourceName.(string), token)
}

func findAdministrativeRoleByUUID(ctx context.Context, api *openapi.AdminRolesApiService, id, token string) (*openapi.AdministrativeRole, diag.Diagnostics) {
	log.Printf("[DEBUG] Data source AdministrativeRole get by UUID %s", id)
	resource, _, err := api.AdministrativeRolesIdGet(ctx, id).Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return resource, nil
}

func findAdministrativeRoleByName(ctx context.Context, api *openapi.AdminRolesApiService, name, token string) (*openapi.AdministrativeRole, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Data source AdministrativeRole get by name %s", name)

	resource, _, err := api.AdministrativeRolesGet(ctx).Query(name).OrderBy("name").Range_("0-10").Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	for _, r := range resource.GetData() {
		if r.GetName() == name {
			return &r, nil
		}
	}
	if len(resource.GetData()) > 1 {
		return nil, AppendErrorf(diags, "multiple AdministrativeRole matched; use additional constraints to reduce matches to a single AdministrativeRole")
	}
	return nil, AppendErrorf(diags, "Failed to find AdministrativeRole %s", name)
}

func ResolveAdministrativeRoleFromResourceData(ctx context.Context, d *schema.ResourceData, api *openapi.AdminRolesApiService, token string) (*openapi.AdministrativeRole, diag.Diagnostics) {
	var diags diag.Diagnostics
	resourceID, iok := d.GetOk("administrative_role_id")
	resourceName, nok := d.GetOk("administrative_role_name")

	if !iok && !nok {
		return nil, AppendErrorf(diags, "please provide one of administrative_role_id or administrative_role_name attributes")
	}

	if iok {
		return findAdministrativeRoleByUUID(ctx, api, resourceID.(string), token)
	}
	return findAdministrativeRoleByName(ctx, api, resourceName.(string), token)
}

func findApplianceCustomizationByUUID(ctx context.Context, api *openapi.ApplianceCustomizationsApiService, id, token string) (*openapi.ApplianceCustomization, diag.Diagnostics) {
	log.Printf("[DEBUG] Data source ApplianceCustomization get by UUID %s", id)
	resource, _, err := api.ApplianceCustomizationsIdGet(ctx, id).Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return resource, nil
}

func findApplianceCustomizationByName(ctx context.Context, api *openapi.ApplianceCustomizationsApiService, name, token string) (*openapi.ApplianceCustomization, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Data source ApplianceCustomization get by name %s", name)

	resource, _, err := api.ApplianceCustomizationsGet(ctx).Query(name).OrderBy("name").Range_("0-10").Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	for _, r := range resource.GetData() {
		if r.GetName() == name {
			return &r, nil
		}
	}
	if len(resource.GetData()) > 1 {
		return nil, AppendErrorf(diags, "multiple ApplianceCustomization matched; use additional constraints to reduce matches to a single ApplianceCustomization")
	}
	return nil, AppendErrorf(diags, "Failed to find ApplianceCustomization %s", name)
}

func ResolveApplianceCustomizationFromResourceData(ctx context.Context, d *schema.ResourceData, api *openapi.ApplianceCustomizationsApiService, token string) (*openapi.ApplianceCustomization, diag.Diagnostics) {
	var diags diag.Diagnostics
	resourceID, iok := d.GetOk("appliance_customization_id")
	resourceName, nok := d.GetOk("appliance_customization_name")

	if !iok && !nok {
		return nil, AppendErrorf(diags, "please provide one of appliance_customization_id or appliance_customization_name attributes")
	}

	if iok {
		return findApplianceCustomizationByUUID(ctx, api, resourceID.(string), token)
	}
	return findApplianceCustomizationByName(ctx, api, resourceName.(string), token)
}

func findApplianceByUUID(ctx context.Context, api *openapi.AppliancesApiService, id, token string) (*openapi.Appliance, diag.Diagnostics) {
	log.Printf("[DEBUG] Data source Appliance get by UUID %s", id)
	resource, _, err := api.AppliancesIdGet(ctx, id).Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return resource, nil
}

func findApplianceByName(ctx context.Context, api *openapi.AppliancesApiService, name, token string) (*openapi.Appliance, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Data source Appliance get by name %s", name)

	resource, _, err := api.AppliancesGet(ctx).Query(name).OrderBy("name").Range_("0-10").Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	for _, r := range resource.GetData() {
		if r.GetName() == name {
			return &r, nil
		}
	}
	if len(resource.GetData()) > 1 {
		return nil, AppendErrorf(diags, "multiple Appliance matched; use additional constraints to reduce matches to a single Appliance")
	}
	return nil, AppendErrorf(diags, "Failed to find Appliance %s", name)
}

func ResolveApplianceFromResourceData(ctx context.Context, d *schema.ResourceData, api *openapi.AppliancesApiService, token string) (*openapi.Appliance, diag.Diagnostics) {
	var diags diag.Diagnostics
	resourceID, iok := d.GetOk("appliance_id")
	resourceName, nok := d.GetOk("appliance_name")

	if !iok && !nok {
		return nil, AppendErrorf(diags, "please provide one of appliance_id or appliance_name attributes")
	}

	if iok {
		return findApplianceByUUID(ctx, api, resourceID.(string), token)
	}
	return findApplianceByName(ctx, api, resourceName.(string), token)
}

func findConditionByUUID(ctx context.Context, api *openapi.ConditionsApiService, id, token string) (*openapi.Condition, diag.Diagnostics) {
	log.Printf("[DEBUG] Data source Condition get by UUID %s", id)
	resource, _, err := api.ConditionsIdGet(ctx, id).Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return resource, nil
}

func findConditionByName(ctx context.Context, api *openapi.ConditionsApiService, name, token string) (*openapi.Condition, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Data source Condition get by name %s", name)

	resource, _, err := api.ConditionsGet(ctx).Query(name).OrderBy("name").Range_("0-10").Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	for _, r := range resource.GetData() {
		if r.GetName() == name {
			return &r, nil
		}
	}
	if len(resource.GetData()) > 1 {
		return nil, AppendErrorf(diags, "multiple Condition matched; use additional constraints to reduce matches to a single Condition")
	}
	return nil, AppendErrorf(diags, "Failed to find Condition %s", name)
}

func ResolveConditionFromResourceData(ctx context.Context, d *schema.ResourceData, api *openapi.ConditionsApiService, token string) (*openapi.Condition, diag.Diagnostics) {
	var diags diag.Diagnostics
	resourceID, iok := d.GetOk("condition_id")
	resourceName, nok := d.GetOk("condition_name")

	if !iok && !nok {
		return nil, AppendErrorf(diags, "please provide one of condition_id or condition_name attributes")
	}

	if iok {
		return findConditionByUUID(ctx, api, resourceID.(string), token)
	}
	return findConditionByName(ctx, api, resourceName.(string), token)
}

func findCriteriaScriptByUUID(ctx context.Context, api *openapi.CriteriaScriptsApiService, id, token string) (*openapi.CriteriaScript, diag.Diagnostics) {
	log.Printf("[DEBUG] Data source CriteriaScript get by UUID %s", id)
	resource, _, err := api.CriteriaScriptsIdGet(ctx, id).Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return resource, nil
}

func findCriteriaScriptByName(ctx context.Context, api *openapi.CriteriaScriptsApiService, name, token string) (*openapi.CriteriaScript, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Data source CriteriaScript get by name %s", name)

	resource, _, err := api.CriteriaScriptsGet(ctx).Query(name).OrderBy("name").Range_("0-10").Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	for _, r := range resource.GetData() {
		if r.GetName() == name {
			return &r, nil
		}
	}
	if len(resource.GetData()) > 1 {
		return nil, AppendErrorf(diags, "multiple CriteriaScript matched; use additional constraints to reduce matches to a single CriteriaScript")
	}
	return nil, AppendErrorf(diags, "Failed to find CriteriaScript %s", name)
}

func ResolveCriteriaScriptFromResourceData(ctx context.Context, d *schema.ResourceData, api *openapi.CriteriaScriptsApiService, token string) (*openapi.CriteriaScript, diag.Diagnostics) {
	var diags diag.Diagnostics
	resourceID, iok := d.GetOk("criteria_script_id")
	resourceName, nok := d.GetOk("criteria_script_name")

	if !iok && !nok {
		return nil, AppendErrorf(diags, "please provide one of criteria_script_id or criteria_script_name attributes")
	}

	if iok {
		return findCriteriaScriptByUUID(ctx, api, resourceID.(string), token)
	}
	return findCriteriaScriptByName(ctx, api, resourceName.(string), token)
}

func findDeviceScriptByUUID(ctx context.Context, api *openapi.DeviceClaimScriptsApiService, id, token string) (*openapi.DeviceScript, diag.Diagnostics) {
	log.Printf("[DEBUG] Data source DeviceScript get by UUID %s", id)
	resource, _, err := api.DeviceScriptsIdGet(ctx, id).Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return resource, nil
}

func findDeviceScriptByName(ctx context.Context, api *openapi.DeviceClaimScriptsApiService, name, token string) (*openapi.DeviceScript, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Data source DeviceScript get by name %s", name)

	resource, _, err := api.DeviceScriptsGet(ctx).Query(name).OrderBy("name").Range_("0-10").Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	for _, r := range resource.GetData() {
		if r.GetName() == name {
			return &r, nil
		}
	}
	if len(resource.GetData()) > 1 {
		return nil, AppendErrorf(diags, "multiple DeviceScript matched; use additional constraints to reduce matches to a single DeviceScript")
	}
	return nil, AppendErrorf(diags, "Failed to find DeviceScript %s", name)
}

func ResolveDeviceScriptFromResourceData(ctx context.Context, d *schema.ResourceData, api *openapi.DeviceClaimScriptsApiService, token string) (*openapi.DeviceScript, diag.Diagnostics) {
	var diags diag.Diagnostics
	resourceID, iok := d.GetOk("device_script_id")
	resourceName, nok := d.GetOk("device_script_name")

	if !iok && !nok {
		return nil, AppendErrorf(diags, "please provide one of device_script_id or device_script_name attributes")
	}

	if iok {
		return findDeviceScriptByUUID(ctx, api, resourceID.(string), token)
	}
	return findDeviceScriptByName(ctx, api, resourceName.(string), token)
}

func findEntitlementScriptByUUID(ctx context.Context, api *openapi.EntitlementScriptsApiService, id, token string) (*openapi.EntitlementScript, diag.Diagnostics) {
	log.Printf("[DEBUG] Data source EntitlementScript get by UUID %s", id)
	resource, _, err := api.EntitlementScriptsIdGet(ctx, id).Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return resource, nil
}

func findEntitlementScriptByName(ctx context.Context, api *openapi.EntitlementScriptsApiService, name, token string) (*openapi.EntitlementScript, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Data source EntitlementScript get by name %s", name)

	resource, _, err := api.EntitlementScriptsGet(ctx).Query(name).OrderBy("name").Range_("0-10").Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	for _, r := range resource.GetData() {
		if r.GetName() == name {
			return &r, nil
		}
	}
	if len(resource.GetData()) > 1 {
		return nil, AppendErrorf(diags, "multiple EntitlementScript matched; use additional constraints to reduce matches to a single EntitlementScript")
	}
	return nil, AppendErrorf(diags, "Failed to find EntitlementScript %s", name)
}

func ResolveEntitlementScriptFromResourceData(ctx context.Context, d *schema.ResourceData, api *openapi.EntitlementScriptsApiService, token string) (*openapi.EntitlementScript, diag.Diagnostics) {
	var diags diag.Diagnostics
	resourceID, iok := d.GetOk("entitlement_script_id")
	resourceName, nok := d.GetOk("entitlement_script_name")

	if !iok && !nok {
		return nil, AppendErrorf(diags, "please provide one of entitlement_script_id or entitlement_script_name attributes")
	}

	if iok {
		return findEntitlementScriptByUUID(ctx, api, resourceID.(string), token)
	}
	return findEntitlementScriptByName(ctx, api, resourceName.(string), token)
}

func findIpPoolByUUID(ctx context.Context, api *openapi.IPPoolsApiService, id, token string) (*openapi.IpPool, diag.Diagnostics) {
	log.Printf("[DEBUG] Data source IpPool get by UUID %s", id)
	resource, _, err := api.IpPoolsIdGet(ctx, id).Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return resource, nil
}

func findIpPoolByName(ctx context.Context, api *openapi.IPPoolsApiService, name, token string) (*openapi.IpPool, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Data source IpPool get by name %s", name)

	resource, _, err := api.IpPoolsGet(ctx).Query(name).OrderBy("name").Range_("0-10").Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	for _, r := range resource.GetData() {
		if r.GetName() == name {
			return &r, nil
		}
	}
	if len(resource.GetData()) > 1 {
		return nil, AppendErrorf(diags, "multiple IpPool matched; use additional constraints to reduce matches to a single IpPool")
	}
	return nil, AppendErrorf(diags, "Failed to find IpPool %s", name)
}

func ResolveIpPoolFromResourceData(ctx context.Context, d *schema.ResourceData, api *openapi.IPPoolsApiService, token string) (*openapi.IpPool, diag.Diagnostics) {
	var diags diag.Diagnostics
	resourceID, iok := d.GetOk("ip_pool_id")
	resourceName, nok := d.GetOk("ip_pool_name")

	if !iok && !nok {
		return nil, AppendErrorf(diags, "please provide one of ip_pool_id or ip_pool_name attributes")
	}

	if iok {
		return findIpPoolByUUID(ctx, api, resourceID.(string), token)
	}
	return findIpPoolByName(ctx, api, resourceName.(string), token)
}

func findLocalUserByUUID(ctx context.Context, api *openapi.LocalUsersApiService, id, token string) (*openapi.LocalUser, diag.Diagnostics) {
	log.Printf("[DEBUG] Data source LocalUser get by UUID %s", id)
	resource, _, err := api.LocalUsersIdGet(ctx, id).Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return resource, nil
}

func findLocalUserByName(ctx context.Context, api *openapi.LocalUsersApiService, name, token string) (*openapi.LocalUser, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Data source LocalUser get by name %s", name)

	resource, _, err := api.LocalUsersGet(ctx).Query(name).OrderBy("name").Range_("0-10").Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	for _, r := range resource.GetData() {
		if r.GetName() == name {
			return &r, nil
		}
	}
	if len(resource.GetData()) > 1 {
		return nil, AppendErrorf(diags, "multiple LocalUser matched; use additional constraints to reduce matches to a single LocalUser")
	}
	return nil, AppendErrorf(diags, "Failed to find LocalUser %s", name)
}

func ResolveLocalUserFromResourceData(ctx context.Context, d *schema.ResourceData, api *openapi.LocalUsersApiService, token string) (*openapi.LocalUser, diag.Diagnostics) {
	var diags diag.Diagnostics
	resourceID, iok := d.GetOk("local_user_id")
	resourceName, nok := d.GetOk("local_user_name")

	if !iok && !nok {
		return nil, AppendErrorf(diags, "please provide one of local_user_id or local_user_name attributes")
	}

	if iok {
		return findLocalUserByUUID(ctx, api, resourceID.(string), token)
	}
	return findLocalUserByName(ctx, api, resourceName.(string), token)
}

func findPolicyByUUID(ctx context.Context, api *openapi.PoliciesApiService, id, token string) (*openapi.Policy, diag.Diagnostics) {
	log.Printf("[DEBUG] Data source Policy get by UUID %s", id)
	resource, _, err := api.PoliciesIdGet(ctx, id).Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return resource, nil
}

func findPolicyByName(ctx context.Context, api *openapi.PoliciesApiService, name, token string) (*openapi.Policy, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Data source Policy get by name %s", name)

	resource, _, err := api.PoliciesGet(ctx).Query(name).OrderBy("name").Range_("0-10").Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	for _, r := range resource.GetData() {
		if r.GetName() == name {
			return &r, nil
		}
	}
	if len(resource.GetData()) > 1 {
		return nil, AppendErrorf(diags, "multiple Policy matched; use additional constraints to reduce matches to a single Policy")
	}
	return nil, AppendErrorf(diags, "Failed to find Policy %s", name)
}

func ResolvePolicyFromResourceData(ctx context.Context, d *schema.ResourceData, api *openapi.PoliciesApiService, token string) (*openapi.Policy, diag.Diagnostics) {
	var diags diag.Diagnostics
	resourceID, iok := d.GetOk("policy_id")
	resourceName, nok := d.GetOk("policy_name")

	if !iok && !nok {
		return nil, AppendErrorf(diags, "please provide one of policy_id or policy_name attributes")
	}

	if iok {
		return findPolicyByUUID(ctx, api, resourceID.(string), token)
	}
	return findPolicyByName(ctx, api, resourceName.(string), token)
}

func findRingfenceRuleByUUID(ctx context.Context, api *openapi.RingfenceRulesApiService, id, token string) (*openapi.RingfenceRule, diag.Diagnostics) {
	log.Printf("[DEBUG] Data source RingfenceRule get by UUID %s", id)
	resource, _, err := api.RingfenceRulesIdGet(ctx, id).Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return resource, nil
}

func findRingfenceRuleByName(ctx context.Context, api *openapi.RingfenceRulesApiService, name, token string) (*openapi.RingfenceRule, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Data source RingfenceRule get by name %s", name)

	resource, _, err := api.RingfenceRulesGet(ctx).Query(name).OrderBy("name").Range_("0-10").Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	for _, r := range resource.GetData() {
		if r.GetName() == name {
			return &r, nil
		}
	}
	if len(resource.GetData()) > 1 {
		return nil, AppendErrorf(diags, "multiple RingfenceRule matched; use additional constraints to reduce matches to a single RingfenceRule")
	}
	return nil, AppendErrorf(diags, "Failed to find RingfenceRule %s", name)
}

func ResolveRingfenceRuleFromResourceData(ctx context.Context, d *schema.ResourceData, api *openapi.RingfenceRulesApiService, token string) (*openapi.RingfenceRule, diag.Diagnostics) {
	var diags diag.Diagnostics
	resourceID, iok := d.GetOk("ringfence_rule_id")
	resourceName, nok := d.GetOk("ringfence_rule_name")

	if !iok && !nok {
		return nil, AppendErrorf(diags, "please provide one of ringfence_rule_id or ringfence_rule_name attributes")
	}

	if iok {
		return findRingfenceRuleByUUID(ctx, api, resourceID.(string), token)
	}
	return findRingfenceRuleByName(ctx, api, resourceName.(string), token)
}

func findSiteByUUID(ctx context.Context, api *openapi.SitesApiService, id, token string) (*openapi.Site, diag.Diagnostics) {
	log.Printf("[DEBUG] Data source Site get by UUID %s", id)
	resource, _, err := api.SitesIdGet(ctx, id).Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return resource, nil
}

func findSiteByName(ctx context.Context, api *openapi.SitesApiService, name, token string) (*openapi.Site, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Data source Site get by name %s", name)

	resource, _, err := api.SitesGet(ctx).Query(name).OrderBy("name").Range_("0-10").Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	for _, r := range resource.GetData() {
		if r.GetName() == name {
			return &r, nil
		}
	}
	if len(resource.GetData()) > 1 {
		return nil, AppendErrorf(diags, "multiple Site matched; use additional constraints to reduce matches to a single Site")
	}
	return nil, AppendErrorf(diags, "Failed to find Site %s", name)
}

func ResolveSiteFromResourceData(ctx context.Context, d *schema.ResourceData, api *openapi.SitesApiService, token string) (*openapi.Site, diag.Diagnostics) {
	var diags diag.Diagnostics
	resourceID, iok := d.GetOk("site_id")
	resourceName, nok := d.GetOk("site_name")

	if !iok && !nok {
		return nil, AppendErrorf(diags, "please provide one of site_id or site_name attributes")
	}

	if iok {
		return findSiteByUUID(ctx, api, resourceID.(string), token)
	}
	return findSiteByName(ctx, api, resourceName.(string), token)
}

func findTrustedCertificateByUUID(ctx context.Context, api *openapi.TrustedCertificatesApiService, id, token string) (*openapi.TrustedCertificate, diag.Diagnostics) {
	log.Printf("[DEBUG] Data source TrustedCertificate get by UUID %s", id)
	resource, _, err := api.TrustedCertificatesIdGet(ctx, id).Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return resource, nil
}

func findTrustedCertificateByName(ctx context.Context, api *openapi.TrustedCertificatesApiService, name, token string) (*openapi.TrustedCertificate, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Data source TrustedCertificate get by name %s", name)

	resource, _, err := api.TrustedCertificatesGet(ctx).Query(name).OrderBy("name").Range_("0-10").Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	for _, r := range resource.GetData() {
		if r.GetName() == name {
			return &r, nil
		}
	}
	if len(resource.GetData()) > 1 {
		return nil, AppendErrorf(diags, "multiple TrustedCertificate matched; use additional constraints to reduce matches to a single TrustedCertificate")
	}
	return nil, AppendErrorf(diags, "Failed to find TrustedCertificate %s", name)
}

func ResolveTrustedCertificateFromResourceData(ctx context.Context, d *schema.ResourceData, api *openapi.TrustedCertificatesApiService, token string) (*openapi.TrustedCertificate, diag.Diagnostics) {
	var diags diag.Diagnostics
	resourceID, iok := d.GetOk("trusted_certificate_id")
	resourceName, nok := d.GetOk("trusted_certificate_name")

	if !iok && !nok {
		return nil, AppendErrorf(diags, "please provide one of trusted_certificate_id or trusted_certificate_name attributes")
	}

	if iok {
		return findTrustedCertificateByUUID(ctx, api, resourceID.(string), token)
	}
	return findTrustedCertificateByName(ctx, api, resourceName.(string), token)
}

func findUserScriptByUUID(ctx context.Context, api *openapi.UserClaimScriptsApiService, id, token string) (*openapi.UserScript, diag.Diagnostics) {
	log.Printf("[DEBUG] Data source UserScript get by UUID %s", id)
	resource, _, err := api.UserScriptsIdGet(ctx, id).Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return resource, nil
}

func findUserScriptByName(ctx context.Context, api *openapi.UserClaimScriptsApiService, name, token string) (*openapi.UserScript, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Data source UserScript get by name %s", name)

	resource, _, err := api.UserScriptsGet(ctx).Query(name).OrderBy("name").Range_("0-10").Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	for _, r := range resource.GetData() {
		if r.GetName() == name {
			return &r, nil
		}
	}
	if len(resource.GetData()) > 1 {
		return nil, AppendErrorf(diags, "multiple UserScript matched; use additional constraints to reduce matches to a single UserScript")
	}
	return nil, AppendErrorf(diags, "Failed to find UserScript %s", name)
}

func ResolveUserScriptFromResourceData(ctx context.Context, d *schema.ResourceData, api *openapi.UserClaimScriptsApiService, token string) (*openapi.UserScript, diag.Diagnostics) {
	var diags diag.Diagnostics
	resourceID, iok := d.GetOk("user_claim_script_id")
	resourceName, nok := d.GetOk("user_claim_script_name")

	if !iok && !nok {
		return nil, AppendErrorf(diags, "please provide one of user_claim_script_id or user_claim_script_name attributes")
	}

	if iok {
		return findUserScriptByUUID(ctx, api, resourceID.(string), token)
	}
	return findUserScriptByName(ctx, api, resourceName.(string), token)
}

func findMfaProviderByUUID(ctx context.Context, api *openapi.MFAProvidersApiService, id, token string) (*openapi.MfaProvider, diag.Diagnostics) {
	log.Printf("[DEBUG] Data source MfaProvider get by UUID %s", id)
	resource, _, err := api.MfaProvidersIdGet(ctx, id).Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return resource, nil
}

func findMfaProviderByName(ctx context.Context, api *openapi.MFAProvidersApiService, name, token string) (*openapi.MfaProvider, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Data source MfaProvider get by name %s", name)

	resource, _, err := api.MfaProvidersGet(ctx).Query(name).OrderBy("name").Range_("0-10").Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	for _, r := range resource.GetData() {
		if r.GetName() == name {
			return &r, nil
		}
	}
	if len(resource.GetData()) > 1 {
		return nil, AppendErrorf(diags, "multiple MfaProvider matched; use additional constraints to reduce matches to a single MfaProvider")
	}
	return nil, AppendErrorf(diags, "Failed to find MfaProvider %s", name)
}

func ResolveMfaProviderFromResourceData(ctx context.Context, d *schema.ResourceData, api *openapi.MFAProvidersApiService, token string) (*openapi.MfaProvider, diag.Diagnostics) {
	var diags diag.Diagnostics
	resourceID, iok := d.GetOk("mfa_provider_id")
	resourceName, nok := d.GetOk("mfa_provider_name")

	if !iok && !nok {
		return nil, AppendErrorf(diags, "please provide one of mfa_provider_id or mfa_provider_name attributes")
	}

	if iok {
		return findMfaProviderByUUID(ctx, api, resourceID.(string), token)
	}
	return findMfaProviderByName(ctx, api, resourceName.(string), token)
}

func findClientProfileByUUID(ctx context.Context, api *openapi.ClientProfilesApiService, id, token string) (*openapi.ClientProfile, diag.Diagnostics) {
	log.Printf("[DEBUG] Data source ClientProfile get by UUID %s", id)
	resource, _, err := api.ClientProfilesIdGet(ctx, id).Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return resource, nil
}

func findClientProfileByName(ctx context.Context, api *openapi.ClientProfilesApiService, name, token string) (*openapi.ClientProfile, diag.Diagnostics) {
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Data source ClientProfile get by name %s", name)

	resource, _, err := api.ClientProfilesGet(ctx).Query(name).OrderBy("name").Range_("0-10").Authorization(token).Execute()
	if err != nil {
		return nil, diag.FromErr(err)
	}
	for _, r := range resource.GetData() {
		if r.GetName() == name {
			return &r, nil
		}
	}
	if len(resource.GetData()) > 1 {
		return nil, AppendErrorf(diags, "multiple ClientProfile matched; use additional constraints to reduce matches to a single ClientProfile")
	}
	return nil, AppendErrorf(diags, "Failed to find ClientProfile %s", name)
}

func ResolveClientProfileFromResourceData(ctx context.Context, d *schema.ResourceData, api *openapi.ClientProfilesApiService, token string) (*openapi.ClientProfile, diag.Diagnostics) {
	var diags diag.Diagnostics
	resourceID, iok := d.GetOk("client_profile_id")
	resourceName, nok := d.GetOk("client_profile_name")

	if !iok && !nok {
		return nil, AppendErrorf(diags, "please provide one of client_profile_id or client_profile_name attributes")
	}

	if iok {
		return findClientProfileByUUID(ctx, api, resourceID.(string), token)
	}
	return findClientProfileByName(ctx, api, resourceName.(string), token)
}
