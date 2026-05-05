package appgate

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/appgate/sdp-api-client-go/api/v24/openapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateReplicationTarget() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppgateReplicationTargetRead,
		Schema: map[string]*schema.Schema{
			"replication_target_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"replication_target_name"},
			},
			"replication_target_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"replication_target_id"},
			},
			"registration_token": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
		},
	}
}

func dataSourceAppgateReplicationTargetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.ReplicationTargetsApi

	replicationTarget, diags := ResolveReplicationTargetFromResourceData(ctx, d, api, token)
	if diags != nil {
		return diags
	}
	authCtx := context.WithValue(ctx, openapi.ContextAccessToken, token)
	replToken, response, err := api.ReplicationTargetsIdExportGet(authCtx, replicationTarget.GetId()).Execute()
	if err != nil && response.StatusCode != http.StatusPreconditionFailed {
		if response != nil && response.StatusCode == http.StatusNotFound {
			log.Printf("[DEBUG] Replication token not found for Replication Target id: %+v", d.Id())
		} else {
			return diag.FromErr(fmt.Errorf("could not retrieve replication token for Replication Target %w", prettyPrintAPIError(err)))
		}
	}
	d.SetId(replicationTarget.GetId())
	d.Set("replication_target_name", replicationTarget.GetName())
	d.Set("replication_target_id", replicationTarget.GetId())
	if replToken != nil {
		d.Set("registration_token", replToken.GetToken())
	}

	return diags
}
