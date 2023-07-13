package appgate

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v19/openapi"
	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var durationPadding, _ = time.ParseDuration("-30s")

func resourceAppgateApplianceControllerActivation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppgateApplianceControllerActivationCreate,
		ReadContext:   resourceAppgateApplianceControllerActivationRead,
		UpdateContext: resourceAppgateApplianceControllerActivationUpdate,
		DeleteContext: resourceAppgateApplianceControllerActivationDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
			Update: schema.DefaultTimeout(45 * time.Minute),
			Delete: schema.DefaultTimeout(45 * time.Minute),
		},
		Schema: func() map[string]*schema.Schema {
			s := map[string]*schema.Schema{
				"appliance_id": {
					Type:        schema.TypeString,
					Description: "ID of the object.",
					Required:    true,
				},
				"controller":      controllerSchema(),
				"admin_interface": adminInterfaceSchema(),
			}
			// all these fields are mandatory within this resource
			s["admin_interface"].Optional = false
			s["admin_interface"].Required = true
			s["controller"].Optional = false
			s["controller"].Required = true
			s["controller"].Computed = false
			return s
		}(),
	}
}

func resourceAppgateApplianceControllerActivationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type  var diags diag.Diagnostics
	var diags diag.Diagnostics
	Appliance54Constraints, err := version.NewConstraint(">= 5.4.0")
	if err != nil {
		return diag.FromErr(err)
	}
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	currentVersion := meta.(*Client).ApplianceVersion
	if !Appliance54Constraints.Check(currentVersion) {
		return diag.FromErr(errors.New("This resource is only available in 5.4 or higher."))
	}

	id := d.Get("appliance_id").(string)
	log.Printf("[DEBUG] Creating appgatesdp_appliance_controller_activation: %s", id)
	api := meta.(*Client).API.AppliancesApi
	request := api.AppliancesIdGet(ctx, id)
	appliance, res, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		if res != nil && res.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.FromErr(fmt.Errorf("Failed to read Appliance, %w", err))
	}

	if v := appliance.GetActivated(); !v {
		return diag.FromErr(fmt.Errorf("Can not activate controller functions on an inactive appliance. The appliance %q need to be seeded first.", appliance.GetName()))
	}
	if v, ok := d.GetOk("controller"); ok {
		ctrl, err := readControllerFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		appliance.SetController(ctrl)
	}
	if a, ok := d.GetOk("admin_interface"); ok {
		ainterface, err := readAdminInterfaceFromConfig(a.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		appliance.SetAdminInterface(ainterface)
	}

	d.SetId(appliance.GetId())
	// if we disable the controller, we want another state.
	state := ApplianceStateApplianceReady
	ctrl := appliance.GetController()
	if ctrl.GetEnabled() {
		state = ApplianceStateControllerReady
	}
	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, _, err := api.AppliancesIdPut(ctx, id).Appliance(*appliance).Authorization(token).Execute()
		if err != nil {
			return resource.NonRetryableError(prettyPrintAPIError(err))
		}
		b := backoff.NewExponentialBackOff()
		b.MaxElapsedTime = d.Timeout(schema.TimeoutCreate) - durationPadding
		// initial sleep; give it a moment for the state to change/update
		time.Sleep(5 * time.Second)
		if err := waitForApplianceState(ctx, meta, id, state, b); err != nil {
			message := fmt.Sprintf("1 or more controller never reached a healthy state after updating controller on %s: %s", appliance.GetName(), err)
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  message,
			})
			return nil
		}
		return nil
	})
	if retryErr != nil {
		if errors.Is(retryErr, context.DeadlineExceeded) {
			return diags
		}
		diags = append(diags, diag.FromErr(retryErr)...)
	}

	diags = append(diags, resourceAppgateApplianceControllerActivationRead(ctx, d, meta)...)
	return diags
}

func resourceAppgateApplianceControllerActivationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type  var diags diag.Diagnostics
	var diags diag.Diagnostics
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.AppliancesApi
	request := api.AppliancesIdGet(ctx, d.Id())
	appliance, res, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		if res != nil && res.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.FromErr(fmt.Errorf("Failed to read Appliance, %w", err))
	}
	d.SetId(appliance.GetId())
	if v, ok := appliance.GetControllerOk(); ok {
		ctrl := make(map[string]interface{})
		ctrl["enabled"] = v.GetEnabled()

		if err := d.Set("controller", []interface{}{ctrl}); err != nil {
			return diag.FromErr(err)
		}
	}
	if v, ok := appliance.GetAdminInterfaceOk(); ok {
		adminInterface, err := flattenApplianceAdminInterface(*v)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("admin_interface", adminInterface); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceAppgateApplianceControllerActivationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	id := d.Get("appliance_id").(string)
	log.Printf("[DEBUG] Updating appgatesdp_appliance_controller_activation: %s", id)
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.AppliancesApi
	request := api.AppliancesIdGet(ctx, id)
	appliance, _, err := request.Authorization(token).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to read Appliance, %w", err))
	}
	if v := appliance.GetActivated(); !v {
		return diag.FromErr(fmt.Errorf("Can not activate controller functions on an inactive appliance. The appliance %q need to be seeded first.", appliance.GetName()))
	}
	if d.HasChange("controller") {
		_, v := d.GetChange("controller")
		ctrl, err := readControllerFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		appliance.SetController(ctrl)
	}
	if d.HasChange("admin_interface") {
		_, v := d.GetChange("admin_interface")
		ainterface, err := readAdminInterfaceFromConfig(v.([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		// since admin_interface is Optional, but admin_interface.hostname is required
		// if it set, we will make sure that hostname is not None before we send the request
		// to avoid sending empty fields in the request body.
		// otherwise, admin_interface has been removed, and we can set admin_interface to nil
		// and it will be omitted from the PUT request.
		if v, ok := ainterface.GetHostnameOk(); ok && v != nil && len(*v) > 0 {
			appliance.SetAdminInterface(ainterface)
		} else {
			appliance.AdminInterface = nil
		}
	}
	// if we disable the controller, we want another state.
	state := ApplianceStateApplianceReady
	ctrl := appliance.GetController()
	if ctrl.GetEnabled() == true {
		state = ApplianceStateControllerReady
	}
	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		_, _, err := api.AppliancesIdPut(ctx, id).Appliance(*appliance).Authorization(token).Execute()
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Could not update appliance %w", prettyPrintAPIError(err)))
		}
		// initial sleep; give it a moment for the state to change/update
		time.Sleep(5 * time.Second)
		b := backoff.NewExponentialBackOff()
		b.MaxElapsedTime = d.Timeout(schema.TimeoutUpdate) - durationPadding

		if err := waitForApplianceState(ctx, meta, id, state, b); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("1 or more controller never reached a healthy state after updating controller on %s: %s", appliance.GetName(), err),
			})
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		if errors.Is(retryErr, context.DeadlineExceeded) {
			return diags
		}
		diags = append(diags, diag.FromErr(retryErr)...)
	}

	diags = append(diags, resourceAppgateApplianceControllerActivationRead(ctx, d, meta)...)
	return diags
}

func resourceAppgateApplianceControllerActivationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type  var diags diag.Diagnostics
	var diags diag.Diagnostics
	id := d.Get("appliance_id").(string)
	log.Printf("[DEBUG] Deleting appgatesdp_appliance_controller_activation: %s", id)
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.AppliancesApi
	request := api.AppliancesIdGet(ctx, id)
	appliance, _, err := request.Authorization(token).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to read Appliance, %w", err))
	}
	c := openapi.ApplianceAllOfController{}
	c.SetEnabled(false)
	appliance.SetController(c)

	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, _, err := api.AppliancesIdPut(ctx, id).Appliance(*appliance).Authorization(token).Execute()
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Could not update appliance when disable controller on %s %w", appliance.Name, prettyPrintAPIError(err)))
		}
		// initial sleep; give it a moment for the state to change/update
		time.Sleep(5 * time.Second)
		b := backoff.NewExponentialBackOff()
		b.MaxElapsedTime = d.Timeout(schema.TimeoutDelete) - durationPadding
		if err := waitForApplianceState(ctx, meta, id, ApplianceStateApplianceReady, b); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("1 or more controller never reached a healthy state after updating controller on %s: %s", appliance.GetName(), err),
			})
		}
		return nil
	})
	if retryErr != nil {
		if errors.Is(retryErr, context.DeadlineExceeded) {
			return diags
		}
		return diag.FromErr(retryErr)
	}
	return diags
}
