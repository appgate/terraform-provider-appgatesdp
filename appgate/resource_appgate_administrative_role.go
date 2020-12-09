package appgate

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v13/openapi"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateAdministrativeRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateAdministrativeRoleCreate,
		Read:   resourceAppgateAdministrativeRoleRead,
		Update: resourceAppgateAdministrativeRoleUpdate,
		Delete: resourceAppgateAdministrativeRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{

			"administrative_role_id": {
				Type:        schema.TypeString,
				Description: "ID of the object.",
				Computed:    true,
			},

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

			"tags": {
				Type:        schema.TypeSet,
				Description: "Array of tags.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"privileges": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
								s := v.(string)
								list := []string{
									"All",
									"View",
									"Create",
									"Edit",
									"Tag",
									"Delete",
									"Revoke",
									"Export",
									"Upgrade",
									"RenewCertificate",
									"DownloadLogs",
									"Test",
									"GetUserAttributes",
									"Backup",
									"CheckStatus",
									"Reevaluate",
								}
								for _, x := range list {
									if s == x {
										return
									}
								}
								errs = append(errs, fmt.Errorf("type must be on of %v, got %s", list, s))
								return
							},
						},

						"target": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, name string) (warns []string, errs []error) {
								s := v.(string)
								list := []string{
									"All",
									"Appliance",
									"Condition",
									"CriteriaScript",
									"Entitlement",
									"AdministrativeRole",
									"IdentityProvider",
									"MfaProvider",
									"IpPool",
									"LocalUser",
									"Policy",
									"Site",
									"DeviceScript",
									"EntitlementScript",
									"RingfenceRule",
									"ApplianceCustomization",
									"OtpSeed",
									"TokenRecord",
									"Blacklist",
									"UserLicense",
									"OnBoardedDevice",
									"AllocatedIp",
									"SessionInfo",
									"AuditLog",
									"AdminMessage",
									"GlobalSetting",
									"CaCertificate",
									"File",
									"FailedAuthentication",
								}
								for _, x := range list {
									if s == x {
										return
									}
								}
								errs = append(errs, fmt.Errorf("type must be on of %v, got %s", list, s))
								return
							},
						},

						"scope": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{

									"all": {
										Type:     schema.TypeBool,
										Optional: true,
									},

									"ids": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},

									"tags": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},

						"default_tags": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceAppgateAdministrativeRoleCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating Administrative role: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.AdministrativeRolesApi
	args := openapi.NewAdministrativeRoleWithDefaults()
	args.Id = uuid.New().String()
	args.SetName(d.Get("name").(string))
	args.SetNotes(d.Get("notes").(string))
	args.SetTags(schemaExtractTags(d))

	if v, ok := d.GetOk("privileges"); ok {
		privileges, err := readAdminIstrativeRolePrivileges(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("Faild to read privileges %s", err)
		}
		args.SetPrivileges(privileges)
	}
	request := api.AdministrativeRolesPost(context.TODO())
	administrativeRole, _, err := request.AdministrativeRole(*args).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create Administrative role %+v", prettyPrintAPIError(err))
	}

	d.SetId(administrativeRole.Id)
	d.Set("administrative_role_id", administrativeRole.Id)

	return resourceAppgateAdministrativeRoleRead(d, meta)
}

func readAdminIstrativeRolePrivileges(privileges []interface{}) ([]openapi.AdministrativePrivilege, error) {
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
				emptyList := make([]string, 0)
				scope := openapi.NewAdministrativePrivilegeScopeWithDefaults()
				scope.SetIds(emptyList)
				scope.SetTags(emptyList)
				scope.SetAll(true)
				for _, v := range rawScopes {
					rawScope := v.(map[string]interface{})
					if v, ok := rawScope["all"]; ok {
						scope.SetAll(v.(bool))
					}

					if v, ok := rawScope["ids"]; ok {
						ids, err := readArrayOfStringsFromConfig(v.([]interface{}))
						if err != nil {
							return result, fmt.Errorf("Failed to resolve privileges scope ids: %+v", err)
						}
						scope.SetIds(ids)
					}
					if v, ok := rawScope["tags"]; ok {
						log.Printf("[DEBUG] readAdminIstrativeRolePrivileges: TAGS CHECK %s ", v.([]interface{}))
						tags, err := readArrayOfStringsFromConfig(v.([]interface{}))
						if err != nil {
							return result, fmt.Errorf("Failed to resolve privileges scope tags: %+v", err)
						}
						scope.SetTags(tags)
					}
				}
				a.SetScope(*scope)
			}
		}

		if v := raw["default_tags"]; len(v.([]interface{})) > 0 {
			tags, err := readArrayOfStringsFromConfig(v.([]interface{}))
			if err != nil {
				return result, fmt.Errorf("Failed to resolve privileges default tags: %+v", err)
			}
			a.SetDefaultTags(tags)
		}
		result = append(result, *a)
	}
	return result, nil
}

func resourceAppgateAdministrativeRoleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Administrative role id: %+v", d.Id())
	token := meta.(*Client).Token
	api := meta.(*Client).API.AdministrativeRolesApi
	ctx := context.TODO()
	request := api.AdministrativeRolesIdGet(ctx, d.Id())
	administrativeRole, _, err := request.Authorization(token).Execute()
	if err != nil {
		// TODO check if 404
		d.SetId("")
		return fmt.Errorf("Failed to read Administrative role, %+v", err)
	}
	d.SetId(administrativeRole.Id)
	d.Set("administrative_role_id", administrativeRole.Id)
	d.Set("name", administrativeRole.Name)
	d.Set("notes", administrativeRole.Notes)
	d.Set("tags", administrativeRole.Tags)

	if administrativeRole.Privileges != nil {
		if err = d.Set("privileges", flattenAdministrativeRolePrivileges(administrativeRole.Privileges)); err != nil {
			return fmt.Errorf("Failed to read privileges %s", err)
		}
	}

	return nil
}

func flattenAdministrativeRolePrivileges(privileges []openapi.AdministrativePrivilege) []map[string]interface{} {
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
			m["default_tags"] = *val
		}

		out[i] = m
	}
	return out
}

func flattenAdministrativeRolePrivilegesScope(scope openapi.AdministrativePrivilegeScope) []interface{} {
	m := make(map[string]interface{})
	if val, ok := scope.GetAllOk(); ok {
		m["all"] = *val
	}
	if val, ok := scope.GetIdsOk(); ok {
		m["ids"] = *val
	}
	if val, ok := scope.GetTagsOk(); ok {
		m["tags"] = *val
	}
	// the response body always include 1 scope
	// example:
	// { "scope":{"all":false,"ids":[],"tags":[]} }
	// but if we dont have it defined in our state, we should not add it now either.
	if len(m["tags"].([]string)) > 0 && len(m["ids"].([]string)) > 0 {
		return []interface{}{m}
	}

	return []interface{}{}
}

func resourceAppgateAdministrativeRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating Administrative role: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.AdministrativeRolesApi
	ctx := context.TODO()
	request := api.AdministrativeRolesIdGet(ctx, d.Id())
	originalAdministrativeRole, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to read Administrative role while updating, %+v", err)
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
		_, n := d.GetChange("privileges")
		privileges, err := readAdminIstrativeRolePrivileges(n.([]interface{}))
		if err != nil {
			return fmt.Errorf("Failed to administrative role privileges %s", err)
		}
		originalAdministrativeRole.SetPrivileges(privileges)
	}

	req := api.AdministrativeRolesIdPut(ctx, d.Id())
	req = req.AdministrativeRole(originalAdministrativeRole)
	_, _, err = req.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update Administrative role %+v", prettyPrintAPIError(err))
	}
	return resourceAppgateAdministrativeRoleRead(d, meta)
}

func resourceAppgateAdministrativeRoleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete Administrative role: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.AdministrativeRolesApi

	request := api.AdministrativeRolesIdDelete(context.TODO(), d.Id())

	_, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not delete Administrative role %+v", prettyPrintAPIError(err))
	}
	d.SetId("")
	return nil
}
