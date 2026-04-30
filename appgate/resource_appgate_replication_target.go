package appgate

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v24/openapi"
	"github.com/appgate/terraform-provider-appgatesdp/appgate/hashcode"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateReplicationTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppgateReplicationTargetCreate,
		ReadContext:   resourceAppgateReplicationTargetRead,
		UpdateContext: resourceAppgateReplicationTargetUpdate,
		DeleteContext: resourceAppgateReplicationTargetDelete,
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
			return mergeSchemaMaps(baseEntitySchema(), map[string]*schema.Schema{
				"replication_target_id": resourceUUID(),
				"replication_tags": {
					Type:        schema.TypeSet,
					Description: "Array of tags.",
					Optional:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					StateFunc: func(val interface{}) string {
						return strings.ToLower(val.(string))
					},
					Set: func(v interface{}) int {
						var buf bytes.Buffer
						str := v.(string)
						buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(str)))
						return hashcode.String(buf.String())
					},
				},
				"registration_token": {
					Type:      schema.TypeString,
					Sensitive: true,
					Computed:  true,
				},
			})
		}(),
	}
}

func resourceAppgateReplicationTargetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating Replication Target: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.ReplicationTargetsApi
	args := openapi.ReplicationTarget{}
	if v, ok := d.GetOk("replication_target_id"); ok {
		args.SetId(v.(string))
	}
	args.SetName(d.Get("name").(string))
	args.SetNotes(d.Get("notes").(string))

	args.SetReplicationTags(schemaExtractReplicationTags(d))

	ctx = context.WithValue(ctx, openapi.ContextAccessToken, token)
	replTarget, _, err := api.ReplicationTargetsPost(ctx).ReplicationTarget(args).Execute()
	if err != nil {
		return diag.FromErr(prettyPrintAPIError(err))
	}

	d.SetId(replTarget.GetId())
	d.Set("replication_target_id", replTarget.GetId())

	return resourceAppgateReplicationTargetRead(ctx, d, meta)
}

func schemaExtractReplicationTags(d *schema.ResourceData) []string {
	rawtags := d.Get("replication_tags").(*schema.Set).List()
	tags := make([]string, 0)
	for _, raw := range rawtags {
		tags = append(tags, strings.ToLower(raw.(string)))
	}
	return tags
}

func resourceAppgateReplicationTargetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading Replication Target id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.ReplicationTargetsApi
	ctx = context.WithValue(ctx, openapi.ContextAccessToken, token)
	replTarget, response, err := api.ReplicationTargetsIdGet(ctx, d.Id()).Execute()
	if err != nil {
		d.SetId("")
		if response != nil && response.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.FromErr(prettyPrintAPIError(err))
	}
	replToken, response, err := api.ReplicationTargetsIdExportGet(ctx, replTarget.GetId()).Execute()
	if err != nil && response.StatusCode != http.StatusPreconditionFailed {
		if response != nil && response.StatusCode == http.StatusNotFound {
			log.Printf("[DEBUG] Replication token not found for Replication Target id: %+v", d.Id())
		} else {
			return diag.FromErr(fmt.Errorf("could not retrieve replication token for Replication Target %w", prettyPrintAPIError(err)))
		}
	}

	d.SetId(replTarget.GetId())
	d.Set("replication_target_id", replTarget.GetId())
	d.Set("name", replTarget.GetName())
	d.Set("notes", replTarget.GetNotes())
	d.Set("replication_tags", replTarget.GetReplicationTags())
	if replToken != nil {
		d.Set("registration_token", replToken.GetToken())
	} else {
		log.Printf("[DEBUG] No replication token available for Replication Target id: %+v", d.Id())
	}

	return nil
}

func resourceAppgateReplicationTargetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Updating Replication Target: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.ReplicationTargetsApi
	ctx = context.WithValue(ctx, openapi.ContextAccessToken, token)
	replTarget, _, err := api.ReplicationTargetsIdGet(ctx, d.Id()).Execute()
	if err != nil {
		return diag.FromErr(prettyPrintAPIError(err))
	}

	if d.HasChange("name") {
		replTarget.SetName(d.Get("name").(string))
	}
	if d.HasChange("notes") {
		replTarget.SetNotes(d.Get("notes").(string))
	}
	if d.HasChange("replication_tags") {
		replTarget.SetReplicationTags(schemaExtractReplicationTags(d))
	}

	ctx = context.WithValue(ctx, openapi.ContextAccessToken, token)
	_, _, err = api.ReplicationTargetsIdPut(ctx, d.Id()).ReplicationTarget(*replTarget).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not update Replication Target %w", prettyPrintAPIError(err)))
	}
	return resourceAppgateReplicationTargetRead(ctx, d, meta)
}

func resourceAppgateReplicationTargetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading Replication Target id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.ReplicationTargetsApi
	ctx = context.WithValue(ctx, openapi.ContextAccessToken, token)
	if _, err := api.ReplicationTargetsIdDelete(ctx, d.Id()).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf("could not delete Replication Target %w", prettyPrintAPIError(err)))
	}
	d.SetId("")
	return nil
}
