package appgate

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v22/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	iso8601Format = "2006-01-02T15:04:05Z0700"
)

func resourceAppgateLocalUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppgateLocalUserCreate,
		ReadContext:   resourceAppgateLocalUserRead,
		UpdateContext: resourceAppgateLocalUserUpdate,
		DeleteContext: resourceAppgateLocalUserDelete,
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
				"local_user_id": resourceUUID(),
				"first_name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"last_name": {
					Type:     schema.TypeString,
					Required: true,
				},
				"password": {
					Type:      schema.TypeString,
					Optional:  true,
					Sensitive: true,
				},
				"email": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"phone": {
					Type:     schema.TypeString,
					Optional: true,
				},
				"failed_login_attempts": {
					Type:     schema.TypeInt,
					Optional: true,
					Computed: true,
				},
				"lock_start": {
					Type:     schema.TypeString,
					Computed: true,
					Optional: true,
				},
			})
		}(),
	}
}

func resourceAppgateLocalUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating Local user: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.LocalUsersApi
	args := openapi.LocalUsersGetRequest{}
	if v, ok := d.GetOk("local_user_id"); ok {
		args.SetId(v.(string))
	}
	args.SetName(d.Get("name").(string))
	args.SetNotes(d.Get("notes").(string))

	args.SetTags(schemaExtractTags(d))

	if v, ok := d.GetOk("first_name"); ok {
		args.SetFirstName(v.(string))
	}

	if v, ok := d.GetOk("last_name"); ok {
		args.SetLastName(v.(string))
	}
	if v, ok := d.GetOk("password"); ok {
		args.SetPassword(v.(string))
	}
	if v, ok := d.GetOk("email"); ok {
		args.SetEmail(v.(string))
	}
	if v, ok := d.GetOk("phone"); ok {
		args.SetPhone(v.(string))
	}
	if v, ok := d.GetOk("failed_login_attempts"); ok {
		args.SetFailedLoginAttempts(float32(v.(int)))
	}
	if v, ok := d.GetOk("lock_start"); ok {
		t, err := parseDateTimeString(v.(string))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Failed to read lock start timestamp %w", err))
		}
		args.SetLockStart(t)
	}

	ctx = context.WithValue(ctx, openapi.ContextAccessToken, token)
	localUser, _, err := api.LocalUsersPost(ctx).LocalUsersGetRequest(args).Execute()
	if err != nil {
		return diag.FromErr(prettyPrintAPIError(err))
	}

	d.SetId(localUser.GetId())
	d.Set("local_user_id", localUser.GetId())

	return resourceAppgateLocalUserRead(ctx, d, meta)
}

func parseDateTimeString(input string) (time.Time, error) {
	return time.Parse(iso8601Format, input)
}

func resourceAppgateLocalUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading Local user id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.LocalUsersApi
	ctx = context.WithValue(ctx, openapi.ContextAccessToken, token)
	localUser, response, err := api.LocalUsersIdGet(ctx, d.Id()).Execute()
	if err != nil {
		d.SetId("")
		if response != nil && response.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.FromErr(prettyPrintAPIError(err))
	}
	d.SetId(localUser.GetId())
	d.Set("local_user_id", localUser.GetId())
	d.Set("name", localUser.GetName())
	d.Set("notes", localUser.GetNotes())
	d.Set("tags", localUser.GetTags())
	d.Set("first_name", localUser.GetFirstName())
	d.Set("last_name", localUser.GetLastName())
	d.Set("email", localUser.GetEmail())
	d.Set("phone", localUser.GetPhone())
	d.Set("failed_login_attempts", localUser.GetFailedLoginAttempts())

	if v, ok := d.GetOk("lock_start"); ok {
		d.Set("lock_start", v.(string))
	}

	return nil
}

func resourceAppgateLocalUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Updating Local user: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.LocalUsersApi
	ctx = context.WithValue(ctx, openapi.ContextAccessToken, token)
	user, _, err := api.LocalUsersIdGet(ctx, d.Id()).Execute()
	if err != nil {
		return diag.FromErr(prettyPrintAPIError(err))
	}

	if d.HasChange("name") {
		user.SetName(d.Get("name").(string))
	}
	if d.HasChange("notes") {
		user.SetNotes(d.Get("notes").(string))
	}
	if d.HasChange("tags") {
		user.SetTags(schemaExtractTags(d))
	}
	if d.HasChange("first_name") {
		user.SetFirstName(d.Get("first_name").(string))
	}
	if d.HasChange("last_name") {
		user.SetLastName(d.Get("last_name").(string))
	}
	if d.HasChange("email") {
		user.SetEmail(d.Get("email").(string))
	}

	if d.HasChange("phone") {
		user.SetPhone(d.Get("phone").(string))
	}
	if d.HasChange("failed_login_attempts") {
		user.SetFailedLoginAttempts(float32(d.Get("failed_login_attempts").(int)))
	}
	if d.HasChange("lock_start") {
		raw := d.Get("lock_start").(string)
		if len(raw) > 0 {
			t, err := parseDateTimeString(raw)
			if err != nil {
				return diag.FromErr(fmt.Errorf("Failed to read lock start timestamp %w", err))
			}
			user.SetLockStart(t)
		}
	}

	ctx = context.WithValue(ctx, openapi.ContextAccessToken, token)
	_, _, err = api.LocalUsersIdPut(ctx, d.Id()).LocalUser(*user).Execute()
	if err != nil {
		return diag.FromErr(fmt.Errorf("could not update Local user %w", prettyPrintAPIError(err)))
	}
	return resourceAppgateLocalUserRead(ctx, d, meta)
}

func resourceAppgateLocalUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading Local user id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.LocalUsersApi
	ctx = context.WithValue(ctx, openapi.ContextAccessToken, token)
	if _, err := api.LocalUsersIdDelete(ctx, d.Id()).Execute(); err != nil {
		return diag.FromErr(fmt.Errorf("could not delete Local user %w", prettyPrintAPIError(err)))
	}
	d.SetId("")
	return nil
}
