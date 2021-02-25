package appgate

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v14/openapi"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateApplianceCustomizations() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateApplianceCustomizationCreate,
		Read:   resourceAppgateApplianceCustomizationRead,
		Update: resourceAppgateApplianceCustomizationUpdate,
		Delete: resourceAppgateApplianceCustomizationDelete,
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

			"appliance_customization_id": {
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

			"file": {
				Type:        schema.TypeString,
				Description: "Path to the appliance customization binary.",
				Required:    true,
			},

			"checksum_sha256": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"size": {
				Type:        schema.TypeInt,
				Description: "Binary file's size in bytes.",
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

func resourceAppgateApplianceCustomizationCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating Appliance customization: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.ApplianceCustomizationsApi
	args := openapi.NewApplianceCustomizationWithDefaults()
	args.Id = uuid.New().String()
	args.SetName(d.Get("name").(string))
	args.SetNotes(d.Get("notes").(string))

	args.SetTags(schemaExtractTags(d))

	content, err := getResourceFileContent(d)
	if err != nil {
		return err
	}
	if len(content) > 0 {
		encoded := base64.StdEncoding.EncodeToString(content)
		args.SetFile(encoded)
	}

	request := api.ApplianceCustomizationsPost(context.TODO())
	request = request.ApplianceCustomization(*args)

	customization, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create Appliance customization %+v", prettyPrintAPIError(err))
	}

	d.SetId(customization.Id)
	d.Set("appliance_customization_id", customization.Id)

	return resourceAppgateApplianceCustomizationRead(d, meta)
}

func resourceAppgateApplianceCustomizationRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Appliance customization id: %+v", d.Id())
	token := meta.(*Client).Token
	api := meta.(*Client).API.ApplianceCustomizationsApi
	ctx := context.TODO()
	request := api.ApplianceCustomizationsIdGet(ctx, d.Id())
	customization, _, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Failed to read Appliance customization, %+v", err)
	}
	d.SetId(customization.Id)
	d.Set("appliance_customization_id", customization.Id)
	d.Set("name", customization.GetName())
	d.Set("notes", customization.GetNotes())
	d.Set("tags", customization.GetTags())
	d.Set("size", customization.GetSize())
	d.Set("checksum_sha256", customization.GetChecksum())

	return nil
}

func resourceAppgateApplianceCustomizationUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating Appliance customization: %s", d.Get("name").(string))
	token := meta.(*Client).Token
	api := meta.(*Client).API.ApplianceCustomizationsApi
	ctx := context.TODO()
	request := api.ApplianceCustomizationsIdGet(ctx, d.Id())
	originalApplianceCustomization, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to read Appliance customization while updating, %+v", err)
	}

	if d.HasChange("name") {
		originalApplianceCustomization.SetName(d.Get("name").(string))
	}

	if d.HasChange("notes") {
		originalApplianceCustomization.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("tags") {
		originalApplianceCustomization.SetTags(schemaExtractTags(d))
	}

	if v := d.Get("file").(string); len(v) > 0 && d.HasChange("file") {
		var content []byte
		file, err := os.Open(v)
		if err != nil {
			return fmt.Errorf("Error opening file (%s): %s", v, err)
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Printf("[WARN] Error closing file (%s): %s", v, err)
			}
		}()
		reader := bufio.NewReader(file)
		content, err = ioutil.ReadAll(reader)
		if err != nil {
			return fmt.Errorf("Error reading file (%s): %s", v, err)
		}
		encoded := base64.StdEncoding.EncodeToString(content)
		originalApplianceCustomization.SetFile(encoded)
	}

	req := api.ApplianceCustomizationsIdPut(ctx, d.Id())
	req = req.ApplianceCustomization(originalApplianceCustomization)
	_, _, err = req.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update Appliance customization %+v", prettyPrintAPIError(err))
	}
	return resourceAppgateApplianceCustomizationRead(d, meta)
}

func resourceAppgateApplianceCustomizationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Appliance customization id: %+v", d.Id())
	token := meta.(*Client).Token
	api := meta.(*Client).API.ApplianceCustomizationsApi

	request := api.ApplianceCustomizationsIdDelete(context.TODO(), d.Id())

	_, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not delete Appliance customization %+v", prettyPrintAPIError(err))
	}
	d.SetId("")
	return nil
}
