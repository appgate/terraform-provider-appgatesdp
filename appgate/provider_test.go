package appgate

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"appgate": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("APPGATE_ADDRESS"); v == "" {
		t.Fatal("APPGATE_ADDRESS must be set for acceptance tests")
	}
	if v := os.Getenv("APPGATE_USERNAME"); v == "" {
		t.Fatal("APPGATE_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("APPGATE_PASSWORD"); v == "" {
		t.Fatal("APPGATE_PASSWORD must be set for acceptance tests")
	}
}
