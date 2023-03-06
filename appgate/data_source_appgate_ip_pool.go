package appgate

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAppgateIPPool() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppgateIPPoolRead,
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

func dataSourceAppgateIPPoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return diag.FromErr(err)
	}
	api := meta.(*Client).API.IPPoolsApi
	ippool, diags := ResolveIpPoolFromResourceData(ctx, d, api, token)
	if diags != nil {
		return diags
	}

	d.SetId(ippool.GetId())
	d.Set("ip_pool_name", ippool.GetName())
	d.Set("lease_time_days", ippool.GetLeaseTimeDays())
	d.Set("total", ippool.GetTotal().String())
	d.Set("currently_used", strconv.FormatInt(ippool.GetCurrentlyUsed(), 10))
	d.Set("reserved", strconv.FormatInt(ippool.GetReserved(), 10))

	return nil
}
