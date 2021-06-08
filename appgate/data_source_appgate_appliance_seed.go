package appgate

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/appgate/sdp-api-client-go/api/v15/openapi"

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
				Computed:    true,
			},
		},
	}
}

func dataSourceAppgateApplianceSeedRead(d *schema.ResourceData, meta interface{}) error {
	token := meta.(*Client).Token
	api := meta.(*Client).API.AppliancesApi
	ctx := context.TODO()
	applianceID, iok := d.GetOk("appliance_id")

	if !iok {
		return fmt.Errorf("please provide one of appliance_id attribute")
	}

	request := api.AppliancesIdGet(ctx, applianceID.(string))
	appliance, _, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Failed to read Appliance, %+v", err)
	}
	if ok, _ := appliance.GetActivatedOk(); *ok {
		d.Set("seed_file", "")
		return fmt.Errorf("Appliance is already activated, %+v", err)
	}

	exportRequest := api.AppliancesIdExportPost(ctx, appliance.Id)
	password, passwordOk := d.GetOk("password")
	sshKey, sshOk := d.GetOk("ssh_key")
	cloudKey, cloudOk := d.GetOk("provide_cloud_ssh_key")

	sshConfig := openapi.NewSSHConfig()
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
		exportRequest.LatestVersion(true)
		d.Set("latest_version", true)
	}
	exportRequest = exportRequest.SSHConfig(*sshConfig)
	seedmap, _, err := exportRequest.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not export appliance %+v", prettyPrintAPIError(err))
	}
	encodedSeed, err := json.Marshal(seedmap)
	if err != nil {
		return fmt.Errorf("Could not parse json seed file: %+v", err)
	}

	d.SetId(applianceID.(string))
	d.Set("appliance_id", appliance.Id)
	d.Set("seed_file", b64.StdEncoding.EncodeToString([]byte(encodedSeed)))

	return nil
}
