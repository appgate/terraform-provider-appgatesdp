package appgate

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestSiteNearestSiteFieldsRoundTrip(t *testing.T) {
	config := map[string]interface{}{
		"name": "test",
		"geolocation": []interface{}{
			map[string]interface{}{"latitude": 59.3293, "longitude": 18.0686},
		},
		"local_site_detection": []interface{}{
			map[string]interface{}{"enabled": true, "public_ips": []interface{}{"192.168.1.1/32"}},
		},
		"use_for_nearest_site_selection": true,
		"fallback_site":                  "5c1a1b8e-7d32-4c5b-9a8e-2b1d3c4e5f60",
	}
	d := schema.TestResourceDataRaw(t, resourceAppgateSite().Schema, config)

	geolocation, err := readGeolocation(d.Get("geolocation").(*schema.Set).List())
	if err != nil {
		t.Fatal(err)
	}
	localSiteDetection, err := readLocalSiteDetection(d.Get("local_site_detection").(*schema.Set).List())
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(map[string]interface{}{"geolocation": geolocation, "localSiteDetection": localSiteDetection})
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"geolocation":{"latitude":59.3293,"longitude":18.0686},"localSiteDetection":{"enabled":true,"publicIps":["192.168.1.1/32"]}}`; string(encoded) != want {
		t.Fatalf("payload = %s, want %s", encoded, want)
	}

	read := schema.TestResourceDataRaw(t, resourceAppgateSite().Schema, map[string]interface{}{"name": "test"})
	if err := read.Set("geolocation", flattenSiteGeolocation(geolocation)); err != nil {
		t.Fatal(err)
	}
	if err := read.Set("local_site_detection", flattenSiteLocalSiteDetection(localSiteDetection)); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"geolocation", "local_site_detection"} {
		if !d.Get(key).(*schema.Set).Equal(read.Get(key)) {
			t.Errorf("%s read back as %v, want %v", key, read.Get(key).(*schema.Set).List(), d.Get(key).(*schema.Set).List())
		}
	}
}
