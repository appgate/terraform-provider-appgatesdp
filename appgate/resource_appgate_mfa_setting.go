package appgate

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAdminMfaSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceAdminMfaSettingsCreate,
		Read:   resourceAdminMfaSettingsRead,
		Update: resourceAdminMfaSettingsUpdate,
		Delete: resourceAdminMfaSettingsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"provider_id": {
				Type:        schema.TypeString,
				Description: "The MFA provider ID to use during Multi-Factor Authentication. If null, Admin MFA is disabled.",
				Optional:    true,
			},
			"exempted_users": {
				Type:        schema.TypeList,
				Description: "The MFA provider ID to use during Multi-Factor Authentication. If null, Admin MFA is disabled.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAdminMfaSettingsCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceAdminMfaSettingsUpdate(d, meta)
}

func resourceAdminMfaSettingsRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading MFA admin settings id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.MFAForAdminsApi
	ctx := BaseAuthContext(token)
	request := api.AdminMfaSettingsGet(ctx)
	settings, _, err := request.Execute()
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Failed to read MFA admin settings, %w", err)
	}
	d.SetId("admin_mfa_settings")
	if v, o := settings.GetProviderIdOk(); o {
		d.Set("provider_id", v)
	}
	if err := d.Set("exempted_users", settings.GetExemptedUsers()); err != nil {
		return err
	}
	return nil
}

func resourceAdminMfaSettingsUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating MFA admin settings")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.MFAForAdminsApi
	ctx := BaseAuthContext(token)
	request := api.AdminMfaSettingsGet(ctx)
	originalsettings, _, err := request.Execute()
	if err != nil {
		return fmt.Errorf("Failed to read MFA admin settings while updating, %w", err)
	}
	d.SetId("admin_mfa_settings")

	if d.HasChange("provider_id") {
		originalsettings.SetProviderId(d.Get("provider_id").(string))
	}
	if d.HasChange("exempted_users") {
		_, v := d.GetChange("exempted_users")
		exemptedUsers, err := readArrayOfStringsFromConfig(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("Failed to read exempted_users %w", err)
		}
		originalsettings.SetExemptedUsers(exemptedUsers)
	}

	log.Printf("[DEBUG] Updating MFA admin settings %+v", originalsettings)
	req := api.AdminMfaSettingsPut(ctx)
	_, err = req.AdminMfaSettings(*originalsettings).Execute()
	if err != nil {
		return fmt.Errorf("Could not update MFA admin settings %w", prettyPrintAPIError(err))
	}

	return resourceAdminMfaSettingsRead(d, meta)
}

func resourceAdminMfaSettingsDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete/Resetting MFA admin settings")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.MFAForAdminsApi

	if _, err := api.AdminMfaSettingsDelete(BaseAuthContext(token)).Execute(); err != nil {
		return fmt.Errorf("Could reset MFA admin settings %w", prettyPrintAPIError(err))
	}
	d.SetId("")
	return nil
}
