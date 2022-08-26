package appgate

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/appgate/sdp-api-client-go/api/v17/openapi"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// errDefaultTagsError is used when trying to use default tags on privileges that does not allow it.
// The items in this list would be added automatically to the newly created objects' tags.
// Only applicable on "Create" type and targets with tagging capability.
// This field must be omitted if not applicable.
var errDefaultTagsError = errors.New("default tags are only applicable on \"Create\" type and targets with tagging capability")

func resourceAppgateAdministrativeRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppgateAdministrativeRoleCreate,
		ReadContext:   resourceAppgateAdministrativeRoleRead,
		UpdateContext: resourceAppgateAdministrativeRoleUpdate,
		DeleteContext: resourceAppgateAdministrativeRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{

			"administrative_role_id": resourceUUID(),

			"name": {
				Type:        schema.TypeString,
				Description: "Name of the object.",
				Required:    true,
			},

			"notes": {
				Type:        schema.TypeString,
				Description: "Notes for the object. Used for documentation purposes.",
				Default:     DefaultDescription,
				Optional:    true,
			},

			"tags": tagsSchema(),

			"privileges": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"type": {
							Type:     schema.TypeString,
							Required: true,
						},

						"target": {
							Type:     schema.TypeString,
							Required: true,
						},

						"scope": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							DefaultFunc: func() (interface{}, error) {
								var out = make([]map[string]interface{}, 0, 0)
								m := make(map[string]interface{})
								m["all"] = false
								emptyList := make([]string, 0)
								m["ids"] = emptyList
								m["tags"] = emptyList
								out = append(out, m)
								return out, nil
							},
							DiffSuppressFunc: suppressMissingOptionalConfigurationBlock,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"all": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},

									"ids": {
										Type:     schema.TypeList,
										Optional: true,
										DefaultFunc: func() (interface{}, error) {
											emptyList := make([]string, 0)
											return emptyList, nil
										},
										Elem: &schema.Schema{Type: schema.TypeString},
									},

									"tags": {
										Type:     schema.TypeList,
										Optional: true,
										DefaultFunc: func() (interface{}, error) {
											emptyList := make([]string, 0)
											return emptyList, nil
										},
										Elem: &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},

						"default_tags": {
							Type:     schema.TypeSet,
							Optional: true,
							DefaultFunc: func() (interface{}, error) {
								return nil, nil
							},
							Elem: &schema.Schema{Type: schema.TypeString},
						},

						"functions": {
							Type:     schema.TypeList,
							Optional: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return strings.EqualFold(old, new)
							},
							Elem: &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceAppgateAdministrativeRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	log.Printf("[DEBUG] Creating Administrative role: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	currentVersion := meta.(*Client).ApplianceVersion
	api := meta.(*Client).API.AdminRolesApi
	args := openapi.NewAdministrativeRoleWithDefaults()
	if v, ok := d.GetOk("administrative_role_id"); ok {
		args.SetId(v.(string))
	}
	args.SetName(d.Get("name").(string))
	args.SetNotes(d.Get("notes").(string))
	args.SetTags(schemaExtractTags(d))

	if v, ok := d.GetOk("privileges"); ok {
		privileges, err := readAdminIstrativeRolePrivileges(v.([]interface{}), currentVersion)
		if err != nil {
			return diag.FromErr(err)
		}
		args.SetPrivileges(privileges)
	}
	request := api.AdministrativeRolesPost(ctx)
	administrativeRole, _, err := request.AdministrativeRole(*args).Authorization(token).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Could not create Administrative role %w", prettyPrintAPIError(err)))
	}

	d.SetId(administrativeRole.GetId())
	d.Set("administrative_role_id", administrativeRole.GetId())

	resourceAppgateAdministrativeRoleRead(ctx, d, meta)

	return diags
}

func readAdminIstrativeRolePrivileges(privileges []interface{}, currentVersion *version.Version) ([]openapi.AdministrativePrivilege, error) {
	result := make([]openapi.AdministrativePrivilege, 0)
	for _, privilege := range privileges {
		if privilege == nil {
			continue
		}
		a := openapi.NewAdministrativePrivilegeWithDefaults()
		raw := privilege.(map[string]interface{})
		if v, ok := raw["type"]; ok {
			a.SetType(v.(string))
		}
		if v, ok := raw["target"]; ok {
			a.SetTarget(v.(string))
		}

		if v, ok := raw["scope"]; ok {
			rawScopes := v.([]interface{})
			if len(rawScopes) > 0 {
				scope := openapi.NewAdministrativePrivilegeScopeWithDefaults()
				for _, v := range rawScopes {
					rawScope := v.(map[string]interface{})
					if v, ok := rawScope["all"]; ok {
						scope.SetAll(v.(bool))
					}

					if v, ok := rawScope["ids"]; ok {
						ids, err := readArrayOfStringsFromConfig(v.([]interface{}))
						if err != nil {
							return result, fmt.Errorf("Failed to resolve privileges scope ids: %w", err)
						}
						scope.SetIds(ids)
					}
					if v, ok := rawScope["tags"]; ok {
						tags, err := readArrayOfStringsFromConfig(v.([]interface{}))
						if err != nil {
							return result, fmt.Errorf("Failed to resolve privileges scope tags: %w", err)
						}
						scope.SetTags(tags)
					}
				}
				a.SetScope(*scope)
			}
		}

		if v, ok := raw["default_tags"]; ok {
			tags, err := readArrayOfStringsFromConfig(v.(*schema.Set).List())
			if err != nil {
				return result, fmt.Errorf("Failed to resolve privileges default tags: %w", err)
			}
			// The items in this list would be added automatically to the newly created objects' tags.
			// Only applicable on "Create" type and targets with tagging capability.
			// This field must be omitted if not applicable.
			if len(tags) > 0 {
				if a.GetType() != "Create" {
					return result, fmt.Errorf("You used %s, %w", a.GetType(), errDefaultTagsError)
				}
				a.SetDefaultTags(tags)
			}

		}
		// client side validation since the controller API does not yet validate it.
		functionAllowedTargets := []string{"Appliance", "All"}
		// lowercase, server side validation does not care about letter case
		allowedFuncs := []string{"controller", "gateway", "logserver", "logforwarder", "connector", "portal"}
		if v, ok := raw["functions"].([]interface{}); ok && len(v) > 0 {
			if currentVersion.LessThan(Appliance60Version) {
				return result, fmt.Errorf("privileges.functions is only supported on >= 6")
			}
			if a.GetType() != "AssignFunction" {
				return result, fmt.Errorf(
					"functions only applicable on \"AssignFunction\" type with target \"Appliance\" or \"All\"."+
						" Got type %s", a.GetType())
			}
			if !inArray(a.GetTarget(), functionAllowedTargets) {
				return result, fmt.Errorf(
					"functions only applicable on \"AssignFunction\" type with target \"Appliance\" or \"All\"."+
						" Got target %s %+v", a.GetTarget(), v)
			}
			for _, f := range v {
				if !inArray(strings.ToLower(f.(string)), allowedFuncs) {
					return result, fmt.Errorf("function must be one of %s, got %s", allowedFuncs, f)
				}
			}
			if _, ok := a.GetScopeOk(); ok {
				return result, fmt.Errorf("Scope is not applicable in combination with privliges.functions")
			}
			funcs, err := readArrayOfStringsFromConfig(v)
			if err != nil {
				return result, fmt.Errorf("Failed to resolve privileges functions %w", err)
			}

			a.SetFunctions(funcs)
		}
		result = append(result, *a)
	}
	return result, nil
}

func resourceAppgateAdministrativeRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Reading Administrative role id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.AdminRolesApi
	request := api.AdministrativeRolesIdGet(ctx, d.Id())
	administrativeRole, res, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		if res != nil && res.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.FromErr(fmt.Errorf("Failed to read Administrative role, %w", err))
	}
	d.SetId(administrativeRole.GetId())
	d.Set("administrative_role_id", administrativeRole.GetId())
	d.Set("name", administrativeRole.GetName())
	d.Set("notes", administrativeRole.GetNotes())
	d.Set("tags", administrativeRole.GetTags())

	privileges, err := flattenAdministrativeRolePrivileges(administrativeRole.Privileges)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("privileges", privileges); err != nil {
		return diag.FromErr(fmt.Errorf("Failed to read privileges %w", err))
	}

	return diags
}

func flattenAdministrativeRolePrivileges(privileges []openapi.AdministrativePrivilege) ([]map[string]interface{}, error) {
	var out = make([]map[string]interface{}, len(privileges), len(privileges))
	for i, v := range privileges {
		m := make(map[string]interface{})
		if val, ok := v.GetTypeOk(); ok {
			m["type"] = *val
		}
		if val, ok := v.GetTargetOk(); ok {
			m["target"] = *val
		}
		if val, ok := v.GetScopeOk(); ok {
			m["scope"] = flattenAdministrativeRolePrivilegesScope(*val)
		}
		if val, ok := v.GetDefaultTagsOk(); ok {
			// The items in this list would be added automatically to the newly created objects' tags.
			// Only applicable on "Create" type and targets with tagging capability.
			// This field must be omitted if not applicable.
			if m["type"] != "Create" {
				return out, fmt.Errorf("You used %s, %w", m["type"], errDefaultTagsError)
			}
			m["default_tags"] = val
		}
		if val, ok := v.GetFunctionsOk(); ok {
			m["functions"] = val
		}
		out[i] = m
	}
	return out, nil
}

func flattenAdministrativeRolePrivilegesScope(scope openapi.AdministrativePrivilegeScope) []interface{} {
	m := make(map[string]interface{})
	if val, ok := scope.GetAllOk(); ok {
		m["all"] = *val
	}
	if val, ok := scope.GetIdsOk(); ok {
		m["ids"] = val
	}
	if val, ok := scope.GetTagsOk(); ok {
		m["tags"] = val
	}
	return []interface{}{m}
}

func resourceAppgateAdministrativeRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Updating Administrative role: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.AdminRolesApi
	currentVersion := meta.(*Client).ApplianceVersion
	request := api.AdministrativeRolesIdGet(ctx, d.Id())
	originalAdministrativeRole, _, err := request.Authorization(token).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to read Administrative role while updating, %w", err))
	}
	if d.HasChange("name") {
		originalAdministrativeRole.SetName(d.Get("name").(string))
	}

	if d.HasChange("notes") {
		originalAdministrativeRole.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("tags") {
		originalAdministrativeRole.SetTags(schemaExtractTags(d))
	}

	if d.HasChange("privileges") {
		_, v := d.GetChange("privileges")
		privileges, err := readAdminIstrativeRolePrivileges(v.([]interface{}), currentVersion)
		if err != nil {
			return diag.FromErr(fmt.Errorf("Failed to administrative role privileges %w", err))
		}
		originalAdministrativeRole.SetPrivileges(privileges)
	}

	req := api.AdministrativeRolesIdPut(ctx, d.Id())
	req = req.AdministrativeRole(*originalAdministrativeRole)
	_, _, err = req.Authorization(token).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Could not update Administrative role %w", prettyPrintAPIError(err)))
	}
	return resourceAppgateAdministrativeRoleRead(ctx, d, meta)
}

func resourceAppgateAdministrativeRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	log.Printf("[DEBUG] Delete Administrative role: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.AdminRolesApi
	if _, err := api.AdministrativeRolesIdDelete(ctx, d.Id()).Authorization(token).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf("Could not delete Administrative role %w", prettyPrintAPIError(err)))
	}
	d.SetId("")
	return diags
}
