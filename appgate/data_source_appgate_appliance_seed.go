package appgate

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/appgate/sdp-api-client-go/api/v16/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateApplianceSeed() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppgateApplianceSeedRead,
		Schema: map[string]*schema.Schema{
			"appliance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"activated": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"provide_cloud_ssh_key": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"ssh_key", "password"},
			},
			"ssh_key": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"provide_cloud_ssh_key", "password"},
			},
			"latest_version": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"provide_cloud_ssh_key", "provide_cloud_ssh_key"},
			},
			"seed_file": {
				Type:        schema.TypeString,
				Description: "Seed file (json) generated from appliance used in remote-exec.",
				Sensitive:   true,
				Computed:    true,
			},
		},
	}
}

func dataSourceAppgateApplianceSeedRead(d *schema.ResourceData, meta interface{}) error {
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.AppliancesApi
	currentVersion := meta.(*Client).ApplianceVersion
	ctx := context.TODO()
	applianceID, iok := d.GetOk("appliance_id")

	if !iok {
		return fmt.Errorf("please provide one of appliance_id attribute")
	}

	request := api.AppliancesIdGet(ctx, applianceID.(string))
	appliance, res, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		if res != nil && res.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("Failed to read Appliance, %w", err)
	}

	d.SetId(applianceID.(string))
	d.Set("appliance_id", appliance.Id)
	d.Set("activated", appliance.GetActivated())

	if ok, _ := appliance.GetActivatedOk(); *ok {
		d.Set("seed_file", "")
		log.Printf("[DEBUG] Appliance is already seeded")
		return nil
	}

	exportRequest := api.AppliancesIdExportPost(ctx, appliance.Id)
	password, passwordOk := d.GetOk("password")
	sshKey, sshOk := d.GetOk("ssh_key")
	cloudKey, cloudOk := d.GetOk("provide_cloud_ssh_key")

	sshConfig := openapi.NewSSHConfig()
	// AllowCustomization and ValidityDays is only available in >= 5.5
	if currentVersion.LessThan(Appliance55Version) {
		sshConfig.AllowCustomization = nil
		sshConfig.ValidityDays = nil
	}
	if passwordOk {
		sshConfig.Password = openapi.PtrString(password.(string))
		d.Set("password", password.(string))
	}
	if sshOk {
		sshConfig.SshKey = openapi.PtrString(sshKey.(string))
		d.Set("ssh_key", sshKey.(string))
	}
	if cloudOk {
		sshConfig.ProvideCloudSSHKey = openapi.PtrBool(cloudKey.(bool))
		d.Set("provide_cloud_ssh_key", true)
	}
	if _, lvOk := d.GetOk("latest_version"); lvOk {
		exportRequest = exportRequest.LatestVersion(true)
		d.Set("latest_version", true)
	}
	exportRequest = exportRequest.SSHConfig(*sshConfig)
	seedmap, _, err := exportRequest.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not export appliance %w", prettyPrintAPIError(err))
	}
	encodedSeed, err := json.Marshal(seedmap)
	if err != nil {
		return fmt.Errorf("Could not parse json seed file: %w", err)
	}

	d.Set("seed_file", b64.StdEncoding.EncodeToString([]byte(encodedSeed)))

	return nil
}
