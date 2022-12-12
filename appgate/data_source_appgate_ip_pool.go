package appgate

import (
	"context"
	"fmt"
	"strconv"

	"github.com/appgate/sdp-api-client-go/api/v18/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateIPPool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppgateIPPoolRead,
		Schema: map[string]*schema.Schema{
			"ip_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip_pool_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"lease_time_days": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"currently_used": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reserved": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAppgateIPPoolRead(d *schema.ResourceData, meta interface{}) error {
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.IPPoolsApi

	ippoolID, iok := d.GetOk("ip_pool_id")
	ippoolName, nok := d.GetOk("ip_pool_name")

	if !iok && !nok {
		return fmt.Errorf("please provide one of ip_pool_id or ip_pool_name attributes")
	}
	var reqErr error
	var ippool *openapi.IpPool
	if iok {
		ippool, reqErr = findIPPoolByUUID(api, ippoolID.(string), token)
	} else {
		ippool, reqErr = findIPPoolByName(api, ippoolName.(string), token)
	}
	if reqErr != nil {
		return reqErr
	}

	d.SetId(ippool.GetId())
	d.Set("ip_pool_name", ippool.GetName())
	d.Set("lease_time_days", ippool.GetLeaseTimeDays())
	d.Set("total", ippool.GetTotal().String())
	d.Set("currently_used", strconv.FormatInt(ippool.GetCurrentlyUsed(), 10))
	d.Set("reserved", strconv.FormatInt(ippool.GetReserved(), 10))

	return nil
}

func findIPPoolByUUID(api *openapi.IPPoolsApiService, id string, token string) (*openapi.IpPool, error) {
	ippool, _, err := api.IpPoolsIdGet(context.Background(), id).Authorization(token).Execute()
	if err != nil {
		return nil, err
	}
	return ippool, nil
}

func findIPPoolByName(api *openapi.IPPoolsApiService, name string, token string) (*openapi.IpPool, error) {
	request := api.IpPoolsGet(context.Background())

	ippool, _, err := request.Query(name).OrderBy("name").Range_("0-1").Authorization(token).Execute()
	if err != nil {
		return nil, err
	}

	for _, c := range ippool.GetData() {
		return &c, nil
	}
	return nil, fmt.Errorf("Failed to find ippool %s", name)
}
