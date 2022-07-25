package appgate

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v17/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	iso8601Format = "2006-01-02T15:04:05Z0700"
)

func resourceAppgateLocalUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateLocalUserCreate,
		Read:   resourceAppgateLocalUserRead,
		Update: resourceAppgateLocalUserUpdate,
		Delete: resourceAppgateLocalUserDelete,
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
					Required:  true,
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

func resourceAppgateLocalUserCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating Local user: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.LocalUsersApi
	args := openapi.NewLocalUserWithDefaults()
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
			return fmt.Errorf("Failed to read lock start timestamp %w", err)
		}
		args.SetLockStart(*t)
	}
	request := api.LocalUsersPost(context.TODO())
	request = request.LocalUser(*args)

	localUser, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create Local user %w", prettyPrintAPIError(err))
	}

	d.SetId(localUser.GetId())
	d.Set("local_user_id", localUser.GetId())

	return resourceAppgateLocalUserRead(d, meta)
}

func parseDateTimeString(input string) (*time.Time, error) {
	t, err := time.Parse(iso8601Format, input)
	if err != nil {
		return nil, fmt.Errorf("Failed resolve datetime %w", err)
	}
	return &t, nil
}

func resourceAppgateLocalUserRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Local user id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.LocalUsersApi
	ctx := context.TODO()
	request := api.LocalUsersIdGet(ctx, d.Id())
	localUser, _, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Failed to read Local user, %w", err)
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

	// TODO: determine if Casting lock_start to a string is enough
	if v, ok := d.GetOk("lock_start"); ok {
		d.Set("lock_start", v.(string))
	}

	return nil
}

func resourceAppgateLocalUserUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating Local user: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.LocalUsersApi
	ctx := context.TODO()
	request := api.LocalUsersIdGet(ctx, d.Id())
	originalLocalUser, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to read Local user while updating, %w", err)
	}
	updatedLocalUser := openapi.NewLocalUserWithDefaults()
	updatedLocalUser.SetPassword(d.Get("password").(string))
	updatedLocalUser.SetId(d.Id())
	if d.HasChange("name") {
		updatedLocalUser.SetName(d.Get("name").(string))
	} else {
		updatedLocalUser.SetName(originalLocalUser.GetName())
	}
	if d.HasChange("notes") {
		updatedLocalUser.SetNotes(d.Get("notes").(string))
	} else {
		updatedLocalUser.SetNotes(originalLocalUser.GetNotes())
	}
	if d.HasChange("tags") {
		updatedLocalUser.SetTags(schemaExtractTags(d))
	} else {
		updatedLocalUser.SetTags(originalLocalUser.GetTags())
	}
	if d.HasChange("first_name") {
		updatedLocalUser.SetFirstName(d.Get("first_name").(string))
	} else {
		updatedLocalUser.SetFirstName(originalLocalUser.GetFirstName())
	}
	if d.HasChange("last_name") {
		updatedLocalUser.SetLastName(d.Get("last_name").(string))
	} else {
		updatedLocalUser.SetLastName(originalLocalUser.GetLastName())
	}

	if d.HasChange("email") {
		updatedLocalUser.SetEmail(d.Get("email").(string))
	} else {
		updatedLocalUser.SetEmail(originalLocalUser.GetEmail())
	}

	if d.HasChange("phone") {
		updatedLocalUser.SetPhone(d.Get("phone").(string))
	} else {
		updatedLocalUser.SetPhone(originalLocalUser.GetPhone())
	}
	if d.HasChange("failed_login_attempts") {
		updatedLocalUser.SetFailedLoginAttempts(float32(d.Get("failed_login_attempts").(int)))
	} else {
		updatedLocalUser.SetFailedLoginAttempts(originalLocalUser.GetFailedLoginAttempts())
	}
	if d.HasChange("lock_start") {
		raw := d.Get("lock_start").(string)
		if len(raw) > 0 {
			t, err := parseDateTimeString(raw)
			if err != nil {
				return fmt.Errorf("Failed to read lock start timestamp %w", err)
			}
			updatedLocalUser.SetLockStart(*t)
		}
	}

	req := api.LocalUsersIdPut(ctx, d.Id())
	req = req.LocalUser(*updatedLocalUser)
	_, _, err = req.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update Local user %w", prettyPrintAPIError(err))
	}
	return resourceAppgateLocalUserRead(d, meta)
}

func resourceAppgateLocalUserDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Local user id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.LocalUsersApi

	if _, err := api.LocalUsersIdDelete(context.TODO(), d.Id()).Authorization(token).Execute(); err != nil {
		return fmt.Errorf("Could not delete Local user %w", prettyPrintAPIError(err))
	}
	d.SetId("")
	return nil
}
