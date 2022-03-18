package appgate

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v16/openapi"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateDeviceScript() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateDeviceScriptCreate,
		Read:   resourceAppgateDeviceScriptRead,
		Update: resourceAppgateDeviceScriptUpdate,
		Delete: resourceAppgateDeviceScriptDelete,
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

			"device_script_id": {
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

			"tags": tagsSchema(),

			"filename": {
				Type:        schema.TypeString,
				Description: "The name of the file to be downloaded as to the client devices.",
				Required:    true,
			},

			"file": {
				Type:          schema.TypeString,
				Description:   "Path to the Device Script binary.",
				Optional:      true,
				ConflictsWith: []string{"content"},
			},

			"content": {
				Type:          schema.TypeString,
				Description:   "The Device Script content.",
				Optional:      true,
				ConflictsWith: []string{"file"},
			},

			"checksum_sha256": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAppgateDeviceScriptCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating Device script: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.DeviceClaimScriptsApi
	args := openapi.NewDeviceScriptWithDefaults()
	args.Id = uuid.New().String()
	args.SetName(d.Get("name").(string))
	args.SetNotes(d.Get("notes").(string))
	args.SetFilename(d.Get("filename").(string))
	args.SetTags(schemaExtractTags(d))

	content, err := getResourceFileContent(d, "file")
	if err != nil {
		return err
	}

	encoded := base64.StdEncoding.EncodeToString(content)
	args.SetFile(encoded)

	request := api.DeviceScriptsPost(context.TODO())
	request = request.DeviceScript(*args)

	deviceScript, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create Device script %w", prettyPrintAPIError(err))
	}

	d.SetId(deviceScript.Id)
	d.Set("device_script_id", deviceScript.Id)

	return resourceAppgateDeviceScriptRead(d, meta)
}

func resourceAppgateDeviceScriptRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Device script id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.DeviceClaimScriptsApi
	ctx := context.TODO()
	request := api.DeviceScriptsIdGet(ctx, d.Id())
	deviceScript, res, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		if res.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("Failed to read Device script, %w", err)
	}
	d.SetId(deviceScript.Id)
	d.Set("device_script_id", deviceScript.Id)
	d.Set("name", deviceScript.Name)
	d.Set("notes", deviceScript.Notes)
	d.Set("tags", deviceScript.Tags)
	d.Set("checksum_sha256", deviceScript.GetChecksumSha256())

	return nil
}

func resourceAppgateDeviceScriptUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating Device script: %s", d.Get("name").(string))
	log.Printf("[DEBUG] Updating Device script id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.DeviceClaimScriptsApi
	ctx := context.TODO()
	request := api.DeviceScriptsIdGet(ctx, d.Id())
	originalDeviceScript, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to read Device script while updating, %w", err)
	}

	if d.HasChange("name") {
		originalDeviceScript.SetName(d.Get("name").(string))
	}

	if d.HasChange("notes") {
		originalDeviceScript.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("tags") {
		originalDeviceScript.SetTags(schemaExtractTags(d))
	}

	if d.HasChange("file") || d.HasChange("content") {
		content, err := getResourceFileContent(d, "file")
		if err != nil {
			return err
		}

		encoded := base64.StdEncoding.EncodeToString(content)
		originalDeviceScript.SetFile(encoded)
	}

	req := api.DeviceScriptsIdPut(ctx, d.Id())
	req = req.DeviceScript(originalDeviceScript)
	_, _, err = req.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update Device script %w", prettyPrintAPIError(err))
	}
	return resourceAppgateDeviceScriptRead(d, meta)
}

func resourceAppgateDeviceScriptDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete Device script: %s", d.Get("name").(string))
	log.Printf("[DEBUG] Reading Device script id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.DeviceClaimScriptsApi

	if _, err := api.DeviceScriptsIdDelete(context.Background(), d.Id()).Authorization(token).Execute(); err != nil {
		return fmt.Errorf("Could not delete Device script %w", prettyPrintAPIError(err))
	}
	d.SetId("")
	return nil
}
