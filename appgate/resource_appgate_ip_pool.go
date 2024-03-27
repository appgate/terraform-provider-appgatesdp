package appgate

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/appgate/sdp-api-client-go/api/v20/openapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAppgateIPPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppgateIPPoolCreate,
		Read:   resourceAppgateIPPoolRead,
		Update: resourceAppgateIPPoolUpdate,
		Delete: resourceAppgateIPPoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{

			"ip_pool_id": resourceUUID(),

			"name": {
				Type:        schema.TypeString,
				Description: "Name of the object.",
				Required:    true,
			},

			"notes": {
				Type:        schema.TypeString,
				Description: "Notes for the object. Used for documentation purposes.",
				Default:     DefaultDescription,
				Optional:    true,
			},

			"tags": tagsSchema(),

			"ip_version6": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"ranges":          ipPoolRange(),
			"excluded_ranges": ipPoolRange(),

			"lease_time_days": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceAppgateIPPoolCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating Ip pool: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.IPPoolsApi
	currentVersion := meta.(*Client).ApplianceVersion
	args := openapi.IpPool{}
	if v, ok := d.GetOk("ip_pool_id"); ok {
		args.SetId(v.(string))
	}
	args.SetName(d.Get("name").(string))
	args.SetNotes(d.Get("notes").(string))

	if v, ok := d.GetOk("ip_version6"); ok {
		args.SetIpVersion6(v.(bool))
	}

	if v, ok := d.GetOk("lease_time_days"); ok {
		args.SetLeaseTimeDays(int32(v.(int)))
	}
	if v, ok := d.GetOk("ranges"); ok {
		ranges, err := readIPPoolRangesFromConfig(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("Failed to read ip pool ranges %w", err)
		}
		args.SetRanges(ranges)
	}

	if currentVersion.GreaterThanOrEqual(Appliance61Version) {
		if v, ok := d.GetOk("excluded_ranges"); ok {
			excludedRanges, err := readIPPoolRangesFromConfig(v.([]interface{}))
			if err != nil {
				return fmt.Errorf("Failed to read ip pool excluded ranges %w", err)
			}
			args.SetExcludedRanges(excludedRanges)
		}
	}

	args.SetTags(schemaExtractTags(d))

	request := api.IpPoolsPost(context.TODO())
	request = request.IpPool(args)

	IPPool, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not create Ip pool %w", prettyPrintAPIError(err))
	}

	d.SetId(IPPool.GetId())
	d.Set("ip_pool_id", IPPool.GetId())

	return resourceAppgateIPPoolRead(d, meta)
}

func readIPPoolRangesFromConfig(ranges []interface{}) ([]openapi.IpPoolRangeInner, error) {
	result := make([]openapi.IpPoolRangeInner, 0)
	for _, ipRange := range ranges {
		if ipRange == nil {
			continue
		}
		r := openapi.IpPoolRangeInner{}
		raw := ipRange.(map[string]interface{})
		if v, ok := raw["first"]; ok {
			r.SetFirst(v.(string))
		}

		if v, ok := raw["last"]; ok {
			r.SetLast(v.(string))
		}
		result = append(result, r)
	}
	return result, nil
}

func resourceAppgateIPPoolRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading Ip pool id: %+v", d.Id())
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.IPPoolsApi
	ctx := context.TODO()
	request := api.IpPoolsIdGet(ctx, d.Id())
	IPPool, res, err := request.Authorization(token).Execute()
	if err != nil {
		d.SetId("")
		if res != nil && res.StatusCode == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("Failed to read Ip pool, %w", err)
	}
	d.SetId(IPPool.GetId())
	d.Set("ip_pool_id", IPPool.GetId())
	d.Set("name", IPPool.GetName())
	d.Set("notes", IPPool.GetNotes())
	d.Set("tags", IPPool.GetTags())
	d.Set("ip_version6", IPPool.IpVersion6)
	d.Set("lease_time_days", IPPool.LeaseTimeDays)
	if ranges, ok := IPPool.GetRangesOk(); ok {
		if err = d.Set("ranges", flattenIPPoolRanges(ranges)); err != nil {
			return fmt.Errorf("Failed to read ip pool ranges %w", err)
		}
	}
	if ranges, ok := IPPool.GetExcludedRangesOk(); ok {
		if err = d.Set("excluded_ranges", flattenIPPoolRanges(ranges)); err != nil {
			return fmt.Errorf("Failed to read ip pool excluded ranges %w", err)
		}
	}

	return nil
}

func flattenIPPoolRanges(in []openapi.IpPoolRangeInner) []map[string]interface{} {
	var out = make([]map[string]interface{}, len(in), len(in))
	for i, v := range in {
		m := make(map[string]interface{})

		if val, ok := v.GetFirstOk(); ok {
			m["first"] = *val
		}
		if val, ok := v.GetLastOk(); ok {
			m["last"] = *val
		}

		out[i] = m
	}
	return out
}

func resourceAppgateIPPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Updating Ip pool: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.IPPoolsApi
	ctx := context.TODO()
	request := api.IpPoolsIdGet(ctx, d.Id())
	originalIPPool, _, err := request.Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Failed to read Ip pool while updating, %w", err)
	}

	if d.HasChange("name") {
		originalIPPool.SetName(d.Get("name").(string))
	}

	if d.HasChange("notes") {
		originalIPPool.SetNotes(d.Get("notes").(string))
	}

	if d.HasChange("tags") {
		originalIPPool.SetTags(schemaExtractTags(d))
	}

	if d.HasChange("ip_version6") {
		originalIPPool.SetIpVersion6(d.Get("ip_version6").(bool))
	}

	if d.HasChange("lease_time_days") {
		originalIPPool.SetLeaseTimeDays(int32(d.Get("lease_time_days").(int)))
	}

	if d.HasChange("ranges") {
		_, n := d.GetChange("ranges")
		ranges, err := readIPPoolRangesFromConfig(n.([]interface{}))
		if err != nil {
			return fmt.Errorf("Failed to read ip pool ranges %w", err)
		}
		originalIPPool.SetRanges(ranges)
	}

	if d.HasChange("excluded_ranges") {
		_, n := d.GetChange("excluded_ranges")
		ranges, err := readIPPoolRangesFromConfig(n.([]interface{}))
		if err != nil {
			return fmt.Errorf("Failed to read ip pool excluded ranges %w", err)
		}
		originalIPPool.SetExcludedRanges(ranges)
	}

	req := api.IpPoolsIdPut(ctx, d.Id())
	_, _, err = req.IpPool(*originalIPPool).Authorization(token).Execute()
	if err != nil {
		return fmt.Errorf("Could not update Ip pool %w", prettyPrintAPIError(err))
	}

	return resourceAppgateIPPoolRead(d, meta)
}

func resourceAppgateIPPoolDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Delete Ip pool: %s", d.Get("name").(string))
	token, err := meta.(*Client).GetToken()
	if err != nil {
		return err
	}
	api := meta.(*Client).API.IPPoolsApi

	if _, err := api.IpPoolsIdDelete(context.TODO(), d.Id()).Authorization(token).Execute(); err != nil {
		return fmt.Errorf("Could not delete Ip pool %w", prettyPrintAPIError(err))
	}
	d.SetId("")
	return nil
}
