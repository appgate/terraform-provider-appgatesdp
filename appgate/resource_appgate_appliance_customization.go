package appgate

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v22/openapi"

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

			"appliance_customization_id": resourceUUID(),

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

			"detect_sha256": {
				Type: schema.TypeString,
				// This field is not Computed because it needs to trigger a diff.
				Optional: true,
				// Makes the diff message nicer:
				// detect_sha256:       "1XcnP/iFw/hNrbhXi7QTmQ==" => "different hash"
				// Instead of the more confusing:
				// detect_sha256:       "1XcnP/iFw/hNrbhXi7QTmQ==" => ""
				Default: "different hash",
				// 1. Compute the sha256 hash of the local file
				// 2. Compare the computed sha256 hash with the hash stored in the Controller
				// 3. Don't suppress the diff iff they don't match
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if source, ok := d.GetOkExists("file"); ok {
						localHash, err := getFileSha256Hash(source.(string))
						if err != nil {
							log.Printf("[WARN] Failed to compute sha256 hash for content: %v", err)
							return false
						}
						if localHash == "" {
							return false
						}
						// `old` is the checksum_sha256 we retrieved from the server in the ReadFunc
						if old != localHash {
							return false
						}
					}
					return true
				},
			},

			"size": {
				Type:        schema.TypeInt,
				Description: "Binary file's size in bytes.",
				Computed:    true,
			},
		},
	}
}

func resourceAppgateApplianceCustomizationCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating Appliance customization: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.ApplianceCustomizationsApi
	args := openapi.NewApplianceCustomizationWithDefaults()
	if v, ok := d.GetOk("appliance_customization_id"); ok {
		args.SetId(v.(string))
	}
	args.SetName(d.Get("name").(string))
	args.SetNotes(d.Get("notes").(string))

	args.SetTags(schemaExtractTags(d))

	content, err := getResourceFileContent(d, "file")
	if err != nil {
		return err
	}
	if len(content) > 0 {
		encoded := base64.StdEncoding.EncodeToString(content)
		args.SetFile(encoded)
	}

	request := api.ApplianceCustomizationsPost(BaseAuthContext(token))
	request = request.ApplianceCustomization(*args)

	customization, _, err := request.Execute()
	if err != nil {
		return fmt.Errorf("Could not create Appliance customization %w", prettyPrintAPIError(err))
	}

	d.SetId(customization.GetId())
	d.Set("appliance_customization_id", customization.GetId())

	return resourceAppgateApplianceCustomizationRead(d, meta)
}

func resourceAppgateApplianceCustomizationRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Appliance customization id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.ApplianceCustomizationsApi
	request := api.ApplianceCustomizationsIdGet(BaseAuthContext(token), d.Id())
	customization, res, err := request.Execute()
	if err != nil {
		d.SetId("")
		if res != nil && res.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("Failed to read Appliance customization, %w", err)
	}
	d.SetId(customization.GetId())
	d.Set("appliance_customization_id", customization.GetId())
	if err := d.Set("name", customization.GetName()); err != nil {
		return fmt.Errorf("Error setting name %w", err)
	}
	if err := d.Set("notes", customization.GetNotes()); err != nil {
		return fmt.Errorf("Error setting notes %w", err)
	}
	if err := d.Set("tags", customization.GetTags()); err != nil {
		return fmt.Errorf("Error setting tags %w", err)
	}
	if err := d.Set("size", customization.GetSize()); err != nil {
		return fmt.Errorf("Error setting size %w", err)
	}
	if err := d.Set("checksum_sha256", customization.GetChecksum()); err != nil {
		return fmt.Errorf("Error setting checksum_sha256 %w", err)
	}
	if err := d.Set("detect_sha256", customization.GetChecksum()); err != nil {
		return fmt.Errorf("Error setting detect_sha256: %w", err)
	}

	return nil
}

func resourceAppgateApplianceCustomizationUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating Appliance customization: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.ApplianceCustomizationsApi
	ctx := BaseAuthContext(token)
	request := api.ApplianceCustomizationsIdGet(ctx, d.Id())
	originalApplianceCustomization, _, err := request.Execute()
	if err != nil {
		return fmt.Errorf("Failed to read Appliance customization while updating, %w", err)
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

	if v := d.Get("file").(string); len(v) > 0 && d.HasChange("detect_sha256") {
		var content []byte
		file, err := os.Open(v)
		if err != nil {
			return fmt.Errorf("Error opening file (%s): %w", v, err)
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Printf("[WARN] Error closing file (%s): %s", v, err)
			}
		}()
		reader := bufio.NewReader(file)
		content, err = io.ReadAll(reader)
		if err != nil {
			return fmt.Errorf("Error reading file (%s): %w", v, err)
		}
		encoded := base64.StdEncoding.EncodeToString(content)
		originalApplianceCustomization.SetFile(encoded)
	}

	req := api.ApplianceCustomizationsIdPut(ctx, d.Id())
	req = req.ApplianceCustomization(*originalApplianceCustomization)
	_, _, err = req.Execute()
	if err != nil {
		return fmt.Errorf("Could not update Appliance customization %w", prettyPrintAPIError(err))
	}
	return resourceAppgateApplianceCustomizationRead(d, meta)
}

func resourceAppgateApplianceCustomizationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Appliance customization id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.ApplianceCustomizationsApi

	ctx := BaseAuthContext(token)
	if _, err := api.ApplianceCustomizationsIdDelete(ctx, d.Id()).Execute(); err != nil {
		return fmt.Errorf("Could not delete Appliance customization %w", prettyPrintAPIError(err))
	}
	d.SetId("")
	return nil
}

func getFileSha256Hash(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("[WARN] Error closing file (%s): %s", filename, err)
		}
	}()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
