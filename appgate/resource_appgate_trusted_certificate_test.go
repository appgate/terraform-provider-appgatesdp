package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTrustedCertificateBasic(t *testing.T) {
	resourceName := "appgate_trusted_certificate.test_trusted_certificate"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTrustedCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTrustedCertificateBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTrustedCertificateExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "cli"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.2876187004", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.535570215", "terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "pem"),
					resource.TestCheckResourceAttrSet(resourceName, "trusted_certificate_id"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccTrustedCertificateImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckTrustedCertificateBasic() string {
	return fmt.Sprintf(`
resource "appgate_trusted_certificate" "test_trusted_certificate" {
  name = "cli"
  tags = [
    "terraform",
    "api-created"
  ]
  pem = <<-EOF
-----BEGIN CERTIFICATE-----
MIICZjCCAc+gAwIBAgIUT0AsBLRI7aKjaMTnH1N9J6eS+7EwDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMDA5MjIxNDQ5MTZaFw0yMTA5
MjIxNDQ5MTZaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwgZ8wDQYJKoZIhvcNAQEB
BQADgY0AMIGJAoGBAOWp5CnfLvNpjeESzTg/B/1kG1BRdXtM00q59WPj7adZ5gq+
+Hr0mWEQ5GldgmXRE3HsXfv7hiq4RwX9h+qtRinwhSvtLquM54/Fpw+TYZl5N27m
ov8a04qqlo8c3BqXR5Vp+ohPVcXs2I21k5bUTh5XwHj4uiv8uxmKzk42WETbAgMB
AAGjUzBRMB0GA1UdDgQWBBSpc1YN7rgPiBrVPn0roGV+1B4ETDAfBgNVHSMEGDAW
gBSpc1YN7rgPiBrVPn0roGV+1B4ETDAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3
DQEBCwUAA4GBAMgxxBlfgH98ME7Es9xlV3HrurwG1p2gBvrrEACMtFNgtZE1vgck
jmhbc3t+Af9Dv9KBkaI6ZDl16uiptdpAv59wLgbVFgEPUJboRjhIaw5mPcMCeSDE
eIE/AV/qHWNEiLIMP5JO2FUbjpDCYtHkCOFDmv01e6rs86L3MQ8zF76T
-----END CERTIFICATE-----
EOF
}
`)
}

func testAccCheckTrustedCertificateExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.TrustedCertificatesApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		_, _, err := api.TrustedCertificatesIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err != nil {
			return fmt.Errorf("error fetching trusted certificate with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckTrustedCertificateDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgate_trusted_certificate" {
			continue
		}

		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.TrustedCertificatesApi

		_, _, err := api.TrustedCertificatesIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err == nil {
			return fmt.Errorf("trusted certificate still exists, %+v", err)
		}
	}
	return nil
}

func testAccTrustedCertificateImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}
