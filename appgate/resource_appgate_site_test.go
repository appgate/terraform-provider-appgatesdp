package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSiteBasic(t *testing.T) {
	resourceName := "appgate_site.test_site"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSite(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v4", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v6", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.excluded_subnets.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "entitlement_based_routing", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_mappings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.access_key_id", "string1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.https_proxy", ""),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.name", "AWS Resolver 1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.regions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.regions.0", "eu-central-1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.regions.1", "eu-west-1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.resolve_with_master_credentials", "true"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.secret_access_key", "string2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.update_interval", "59"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.use_iam_role", "true"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.vpc_auto_discovery", "true"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.vpcs.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.0.client_id", "string3"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.0.name", "Azure Resolver 1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.0.secret", "string4"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.0.subscription_id", "string1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.0.tenant_id", "string2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.0.update_interval", "30"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.name", "DNS Resolver 1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.search_domains.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.search_domains.0", "hostname.dns"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.search_domains.1", "foo.bar"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.servers.0", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.servers.1", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.update_interval", "13"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.esx_resolvers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.esx_resolvers.0.hostname", "string1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.esx_resolvers.0.name", "ESX Resolver 1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.esx_resolvers.0.password", "secret_password"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.esx_resolvers.0.update_interval", "120"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.esx_resolvers.0.username", "string2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.gcp_resolvers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.gcp_resolvers.0.instance_filter", "string2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.gcp_resolvers.0.name", "GCP Resolver 1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.gcp_resolvers.0.project_filter", "string1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.gcp_resolvers.0.update_interval", "360"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.use_hosts_file", "false"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "notes", "This object has been created for test purposes."),
					resource.TestCheckResourceAttr(resourceName, "short_name", "tst"),

					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "developer"),
					resource.TestCheckResourceAttr(resourceName, "vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.ip_access_log_interval_seconds", "120"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.snat", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.state_sharing", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.web_proxy_enabled", "false"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccCriteriaScripImportStateCheckFunc(1),
			},
			{
				Config: testAccCheckSiteUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "The test site"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.1", "10.20.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "developer"),
					resource.TestCheckResourceAttr(resourceName, "tags.2", "updated-tag"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccCriteriaScripImportStateCheckFunc(1),
			},
			{
				Config: testAccCheckSiteNetworkDelete(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "The test site"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.0", "10.20.0.0/24"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccSiteImportStateCheckFunc(1),
			},
		},
	})
}

func testAccSiteImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %#v", expectedStates, len(s), s)
		}
		return nil
	}
}

func testAccCheckSiteDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgate_site" {
			continue
		}

		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.SitesApi

		_, _, err := api.SitesIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err == nil {
			return fmt.Errorf("Site still exists, %+v", err)
		}
	}
	return nil
}

func testAccCheckSite(rName string) string {
	return fmt.Sprintf(`
resource "appgate_site" "test_site" {
    name       = "%s"
    short_name = "tst"
    tags = [
        "developer",
        "api-created"
    ]

    notes = "This object has been created for test purposes."
    entitlement_based_routing = false
    network_subnets = [
        "10.0.0.0/16"
    ]

    default_gateway {
        enabled_v4       = false
        enabled_v6       = false
        excluded_subnets = []
    }

    name_resolution {

        dns_resolvers {
            name            = "DNS Resolver 1"
            update_interval = 13
            servers = [
                "8.8.8.8",
                "1.1.1.1"
            ]
            search_domains = [
                "hostname.dns",
                "foo.bar"
            ]
        }

        aws_resolvers {
            name = "AWS Resolver 1"
            regions = [
                "eu-central-1",
                "eu-west-1"
            ]
            update_interval    = 59
            vpcs               = []
            vpc_auto_discovery = true
            use_iam_role       = true
            access_key_id      = "string1"
            secret_access_key  = "string2"
            resolve_with_master_credentials = true
        }

        azure_resolvers {
            name            = "Azure Resolver 1"
            update_interval = 30
            subscription_id = "string1"
            tenant_id       = "string2"
            client_id       = "string3"
            secret          = "string4"
        }

        esx_resolvers {
            name            = "ESX Resolver 1"
            update_interval = 120
            hostname        = "string1"
            username        = "string2"
            password        = "secret_password"
        }

        gcp_resolvers {
            name            = "GCP Resolver 1"
            update_interval = 360
            project_filter  = "string1"
            instance_filter = "string2"
        }
    }
}
`, rName)
}

func testAccCheckSiteUpdate() string {
	return `
resource "appgate_site" "test_site" {
    name       = "The test site"
    short_name = "tst"
    tags = [
        "developer",
        "api-created",
	    "updated-tag"
    ]
    notes = "This object has been created for test purposes."
    entitlement_based_routing = false
	network_subnets = [
        "10.0.0.0/16",
        "10.20.0.0/24",
    ]

    default_gateway {
        enabled_v4       = false
        enabled_v6       = false
        excluded_subnets = []
    }

    name_resolution {
        dns_resolvers {
            name            = "DNS Resolver 1"
            update_interval = 13
            servers = [
                "8.8.8.8",
                "1.1.1.1"
            ]
            search_domains = [
                "hostname.dns",
                "foo.bar"
            ]
        }

        aws_resolvers {
            name = "AWS Resolver 1"
            regions = [
                "eu-central-1",
                "eu-west-1"
            ]
            update_interval    = 59
            vpcs               = []
            vpc_auto_discovery = true
            use_iam_role       = true
            access_key_id      = "string1"
            secret_access_key  = "string2"
            resolve_with_master_credentials = true
        }

        azure_resolvers {
            name            = "Azure Resolver 1"
            update_interval = 30
            subscription_id = "string1"
            tenant_id       = "string2"
            client_id       = "string3"
            secret          = "string4"
        }

        esx_resolvers {
            name            = "ESX Resolver 1"
            update_interval = 120
            hostname        = "string1"
            username        = "string2"
            password        = "secret_password"
        }

        gcp_resolvers {
            name            = "GCP Resolver 1"
            update_interval = 360
            project_filter  = "string1"
            instance_filter = "string2"
        }
    }
}
`
}

func testAccCheckSiteNetworkDelete() string {
	return `
resource "appgate_site" "test_site" {
    name       = "The test site"
    short_name = "tst"
    tags = [
        "developer",
        "api-created",
    	"updated-tag"
    ]
    notes = "This object has been created for test purposes."
    entitlement_based_routing = false
    network_subnets = [
        "10.20.0.0/24"
    ]

    default_gateway {
        enabled_v4       = false
        enabled_v6       = false
        excluded_subnets = []
    }

    name_resolution {
        dns_resolvers {
            name            = "DNS Resolver 1"
            update_interval = 13
            servers = [
                "8.8.8.8",
                "1.1.1.1"
            ]
            search_domains = [
                "hostname.dns",
                "foo.bar"
            ]
        }

        aws_resolvers {
            name = "AWS Resolver 1"
            update_interval    = 59
            vpcs               = ["test"]
            vpc_auto_discovery = true
            use_iam_role       = true
            access_key_id      = "string1"
            secret_access_key  = "string2"
            resolve_with_master_credentials = true
        }

        azure_resolvers {
            name            = "Azure Resolver 1"
            update_interval = 30
            subscription_id = "string1"
            tenant_id       = "string2"
            client_id       = "string3"
            secret          = "string4"
        }

        esx_resolvers {
            name            = "ESX Resolver 1"
            update_interval = 120
            hostname        = "string1"
            username        = "string2"
            password        = "secret_password"
        }

        gcp_resolvers {
            name            = "GCP Resolver 1"
            update_interval = 360
            project_filter  = "string1"
            instance_filter = "string2"
        }
    }
}
`
}

func testAccCheckSiteExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.SitesApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		_, _, err := api.SitesIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}
