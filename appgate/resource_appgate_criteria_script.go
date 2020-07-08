package appgate

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/appgate/terraform-provider-appgate/client/v12/openapi"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAppgateCriteriaScript() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateCriteriaScriptCreate,
		Read:   resourceAppgateCriteriaScriptRead,
		Update: resourceAppgateCriteriaScriptUpdate,
		Delete: resourceAppgateCriteriaScriptDelete,
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

			"criteria_script_id": {
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

			"expression": {
				Type:        schema.TypeString,
				Description: "Name of the object.",
				Required:    true,
			},
		},
	}
}

func resourceAppgateCriteriaScriptCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating Criteria script: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.CriteriaScriptsApi
	args := openapi.NewCriteriaScriptWithDefaults()
	args.Id = uuid.New().String()

	request := api.CriteriaScriptsPost(context.TODO())
	request = request.CriteriaScript(*args)
	criteraScript, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create Criteria script %+v", prettyPrintAPIError(err))
	}

	d.SetId(criteraScript.Id)
	d.Set("criteria_script_id", criteraScript.Id)

	return resourceAppgateCriteriaScriptRead(d, meta)
}

func resourceAppgateCriteriaScriptRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Criteria script id: %+v", d.Id())
	token := meta.(*Client).Token
	api := meta.(*Client).API.CriteriaScriptsApi
	ctx := context.TODO()
	request := api.CriteriaScriptsIdGet(ctx, d.Id())
	ent, _, err := request.Authorization(token).Execute()
	if err != nil {
		// TODO check if 404
		d.SetId("")
		return fmt.Errorf("Failed to read Criteria script, %+v", err)
	}
	d.SetId(ent.Id)
	d.Set("criteria_script_id", ent.Id)

	return nil
}

func resourceAppgateCriteriaScriptUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating Criteria script: %s", d.Get("name").(string))
	log.Printf("[DEBUG] Updating Criteria script id: %+v", d.Id())
	token := meta.(*Client).Token
	api := meta.(*Client).API.CriteriaScriptsApi
	ctx := context.TODO()
	request := api.CriteriaScriptsIdGet(ctx, d.Id())
	originalCriteriaScript, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to read Criteria script while updating, %+v", err)
	}
	// ....

	req := api.CriteriaScriptsIdPut(ctx, d.Id())
	req = req.CriteriaScript(originalCriteriaScript)
	_, _, err = req.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update Criteria script %+v", prettyPrintAPIError(err))
	}
	return resourceAppgateCriteriaScriptRead(d, meta)
}

func resourceAppgateCriteriaScriptDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete Criteria script: %s", d.Get("name").(string))
	log.Printf("[DEBUG] Reading Criteria script id: %+v", d.Id())
	token := meta.(*Client).Token
	api := meta.(*Client).API.CriteriaScriptsApi

	request := api.CriteriaScriptsIdDelete(context.TODO(), d.Id())

	_, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not delete Criteria script %+v", prettyPrintAPIError(err))
	}
	d.SetId("")
	return nil
}
