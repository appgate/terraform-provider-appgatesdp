package appgate

import (
	"context"
	"fmt"
	"log"

	"github.com/appgate/sdp-api-client-go/api/v16/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateLocalUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppgateLocalUserRead,
		Schema: map[string]*schema.Schema{
			"local_user_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"local_user_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceAppgateLocalUserRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Data source local user")
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.LocalUsersApi

	userID, iok := d.GetOk("local_user_id")
	userName, nok := d.GetOk("local_user_name")

	if !iok && !nok {
		return fmt.Errorf("please provide one of local_user_id or local_user_name attributes")
	}
	var reqErr error
	var localUser *openapi.LocalUser
	if iok {
		localUser, reqErr = findLocalUserByUUID(api, userID.(string), token)
	} else {
		localUser, reqErr = findLocalUserByName(api, userName.(string), token)
	}
	if reqErr != nil {
		return reqErr
	}
	log.Printf("[DEBUG] Got local user: %+v", localUser)

	d.SetId(localUser.Id)
	d.Set("local_user_name", localUser.Name)
	d.Set("local_user_id", localUser.Id)
	return nil
}

func findLocalUserByUUID(api *openapi.LocalUsersApiService, id string, token string) (*openapi.LocalUser, error) {
	log.Printf("[DEBUG] Data source local user get by UUID %s", id)
	localUser, _, err := api.LocalUsersIdGet(context.Background(), id).Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	return &localUser, nil
}

func findLocalUserByName(api *openapi.LocalUsersApiService, name string, token string) (*openapi.LocalUser, error) {
	log.Printf("[DEBUG] Data local user get by name %s", name)
	request := api.LocalUsersGet(context.Background())

	localUser, _, err := request.Query(name).OrderBy("name").Range_("0-1").Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	for _, s := range localUser.GetData() {
		return &s, nil
	}
	return nil, fmt.Errorf("Failed to find local user %s", name)
}
