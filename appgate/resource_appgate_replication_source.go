package appgate

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v24/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateReplicationSource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppgateReplicationSourceCreate,
		ReadContext:   resourceAppgateReplicationSourceRead,
		DeleteContext: resourceAppgateReplicationSourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		SchemaVersion: 1,
		Schema: func() map[string]*schema.Schema {
			return map[string]*schema.Schema{
				"replication_source_id": resourceUUID(),
				"registration_token": {
					Type:      schema.TypeString,
					Sensitive: true,
					Required:  true,
					ForceNew:  true,
				},
				"status": {
					Type:     schema.TypeString,
					Computed: true,
				},
			}
		}(),
	}
}

/*
Note: The Replication Source resource is unique and does not have an ID that can be used for CRUD operations.

Therefore, we will use a fixed ID for the Terraform resource and manage the lifecycle based
on the status of the Replication Source API.
*/
var id = "e791af5e-cb24-4205-b447-1ea0160830f6"

func resourceAppgateReplicationSourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating Replication Source")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		log.Printf("[DEBUG] Error retrieving token: %v", err)
		return diag.FromErr(err)
	}
	if d.Get("status") != "notConnected" && d.Get("status") != "" {
		log.Printf("[DEBUG] Replication Source already exists with status: %s, skipping creation", d.Get("status").(string))
		return resourceAppgateReplicationSourceRead(ctx, d, meta)
	}
	api := meta.(*Client).API.ReplicationSourceApi
	registrationToken := openapi.ReplicationRegistrationToken{}
	if v, ok := d.GetOk("registration_token"); ok {
		registrationToken.SetToken(v.(string))
	}

	log.Printf("[DEBUG] Registration token for Replication Source creation: %s", registrationToken.GetToken())

	ctx = context.WithValue(ctx, openapi.ContextAccessToken, token)
	_, _, err = api.ReplicationSourcePost(ctx).ReplicationRegistrationToken(registrationToken).Execute()
	if err != nil {
		log.Printf("[DEBUG] Error creating Replication Source: %v", err)
		return diag.FromErr(prettyPrintAPIError(err))
	}

	d.SetId(id)
	return resourceAppgateReplicationSourceRead(ctx, d, meta)
}

func resourceAppgateReplicationSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading Replication Source")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.ReplicationSourceApi
	ctx = context.WithValue(ctx, openapi.ContextAccessToken, token)
	replSource, response, err := api.ReplicationSourceGet(ctx).Execute()
	if err != nil {
		d.SetId("")
		if response != nil && response.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.FromErr(prettyPrintAPIError(err))
	}

	d.Set("status", replSource.GetStatus())
	d.SetId(id)
	log.Printf("[DEBUG] Replication Source status: %s", replSource.GetStatus())
	return nil
}

func resourceAppgateReplicationSourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Deleting Replication Source")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.ReplicationSourceApi
	ctx = context.WithValue(ctx, openapi.ContextAccessToken, token)
	if _, response, err := api.ReplicationSourceDelete(ctx).Execute(); err != nil && response != nil && response.StatusCode != http.StatusPreconditionFailed {
		return diag.FromErr(fmt.Errorf("could not delete Replication Source %w", prettyPrintAPIError(err)))
	}
	d.SetId(id)
	return nil
}
