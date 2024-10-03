package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSiteBasic(t *testing.T) {
	resourceName := "appgatesdp_site.test_site"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					// TODO
					// test includes
					// name_resolution.0.dns_resolvers.0.search_domains.#
					// which has been removed in 6.2
					applianceConstraintCheck(t, ">= 5.5, < 6.2")
				},
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
					resource.TestCheckResourceAttr(resourceName, "short_name", "ts0"),

					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "developer"),
					resource.TestCheckResourceAttr(resourceName, "vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.ip_access_log_interval_seconds", "120"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.snat", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.enabled", "true"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccSiteImportStateCheckFunc(1),
			},
			{
				Config: testAccCheckSiteUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "The test site"),
					resource.TestCheckResourceAttr(resourceName, "short_name", "ts1"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.1", "10.20.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "developer"),
					resource.TestCheckResourceAttr(resourceName, "tags.2", "updated-tag"),
					resource.TestCheckResourceAttr(resourceName, "notes", "note updated"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccSiteImportStateCheckFunc(1),
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
			{
				Config: testAccCheckSiteTagsDelete(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccSiteImportStateCheckFunc(1),
			},
			{
				Config: testAccCheckSiteTagsAdd(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),

					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "qwerty"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccSiteImportStateCheckFunc(1),
			},
			{
				Config: testAccCheckSiteTagsDelete(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
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
		if rs.Type != "appgatesdp_site" {
			continue
		}

		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.SitesApi

		if _, _, err := api.SitesIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute(); err == nil {
			return fmt.Errorf("Site still exists, %+v", err)
		}
	}
	return nil
}

func testAccCheckSite(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_site" "test_site" {
    name       = "%s"
    short_name = "ts0"
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
            name                   = "Azure Resolver 1"
            update_interval        = 30
            use_managed_identities = true
            tenant_id              = "string2"
            client_id              = "string3"
            secret                 = "string4"
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
resource "appgatesdp_site" "test_site" {
    name       = "The test site"
    short_name = "ts1"
    tags = [
        "developer",
        "api-created",
	    "updated-tag"
    ]
    notes = "note updated"
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
            name                   = "Azure Resolver 1"
            update_interval        = 30
            use_managed_identities = true
            tenant_id              = "string2"
            client_id              = "string3"
            secret                 = "string4"
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
resource "appgatesdp_site" "test_site" {
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
            name                   = "Azure Resolver 1"
            update_interval        = 30
            tenant_id              = "string2"
            client_id              = "string3"
            secret                 = "string4"
            use_managed_identities = true
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

func testAccCheckSiteTagsDelete() string {
	return `
resource "appgatesdp_site" "test_site" {
    name       = "The test site"
    short_name = "tst"
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
            name                   = "Azure Resolver 1"
            update_interval        = 30
            tenant_id              = "string2"
            client_id              = "string3"
            secret                 = "string4"
            use_managed_identities = true
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

func testAccCheckSiteTagsAdd() string {
	return `
resource "appgatesdp_site" "test_site" {
    name       = "The test site"
    short_name = "tst"
	tags = ["qwerty"]
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
            name                   = "Azure Resolver 1"
            update_interval        = 30
            tenant_id              = "string2"
            client_id              = "string3"
            secret                 = "string4"
			use_managed_identities = true
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
		token, err := testAccProvider.Meta().(*Client).GetToken()
		if err != nil {
			return err
		}
		api := testAccProvider.Meta().(*Client).API.SitesApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		if _, _, err := api.SitesIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute(); err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func TestAccSiteBasicAwsResolverWithoutSecret(t *testing.T) {
	resourceName := "appgatesdp_site.test_site"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	context := map[string]interface{}{
		"name": rName,
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteBasicAwsResolverConfig(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v4", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v6", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.excluded_subnets.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "entitlement_based_routing", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_mappings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", context["name"].(string)),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.%", "12"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.access_key_id", ""),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.https_proxy", ""),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.name", "AWS Resolver 10"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.regions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.regions.0", "eu-central-1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.regions.1", "eu-west-1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.resolve_with_master_credentials", "true"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.secret_access_key", ""),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.update_interval", "59"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.use_iam_role", "true"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.vpc_auto_discovery", "true"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.vpcs.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.esx_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.gcp_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.use_hosts_file", "false"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "short_name", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "developer"),
					resource.TestCheckResourceAttr(resourceName, "vpn.#", "1"),
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

func testAccSiteBasicAwsResolverConfig(context map[string]interface{}) string {
	return Nprintf(`
    resource "appgatesdp_site" "test_site" {
        name       = "%{name}"
        tags = [
          "developer",
          "api-created"
        ]
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
          aws_resolvers {
            name = "AWS Resolver 10"
            regions = [
              "eu-central-1",
              "eu-west-1"
            ]
            update_interval                 = 59
            vpcs                            = []
            vpc_auto_discovery              = true
            use_iam_role                    = true
            resolve_with_master_credentials = true
          }
        }
      }
    `, context)
}

// https://github.com/appgate/terraform-provider-appgatesdp/issues/229
func TestAccSiteBasicAwsResolverresolveWithMasterCredentials(t *testing.T) {
	resourceName := "appgatesdp_site.test_site"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	context := map[string]interface{}{
		"name": rName,
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteBasicAwsResolverConfiWithMasterCredentials(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v4", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v6", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.excluded_subnets.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "entitlement_based_routing", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_mappings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", context["name"].(string)),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.%", "12"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.access_key_id", "string1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.0.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.0.account_id", "abc123"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.0.external_id", "3a569552-48a4-4589-990d-5f1d8e3a6a18"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.0.regions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.0.regions.0", "eu-central-1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.0.role_name", "The role"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.https_proxy", ""),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.name", "AWS Resolver 1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.regions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.regions.0", "eu-central-1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.regions.1", "eu-west-1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.resolve_with_master_credentials", "false"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.secret_access_key", "string2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.update_interval", "59"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.use_iam_role", "true"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.vpc_auto_discovery", "true"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.vpcs.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.esx_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.gcp_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.use_hosts_file", "false"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "short_name", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "developer"),
					resource.TestCheckResourceAttr(resourceName, "vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.ip_access_log_interval_seconds", "120"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.snat", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.enabled", "true"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccSiteImportStateCheckFunc(1),
			},
			{
				Config: testAccSiteBasicAwsResolverConfiWithMasterCredentialsUpdated(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v4", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v6", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.excluded_subnets.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "entitlement_based_routing", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_mappings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", context["name"].(string)),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.%", "12"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.access_key_id", "string1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.0.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.0.account_id", "abc123"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.0.external_id", "3a569552-48a4-4589-990d-5f1d8e3a6a18"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.0.regions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.0.regions.0", "eu-central-1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.0.role_name", "The role"),
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
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.esx_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.gcp_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.use_hosts_file", "false"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "short_name", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "developer"),
					resource.TestCheckResourceAttr(resourceName, "vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.ip_access_log_interval_seconds", "120"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.snat", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.enabled", "true"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccSiteImportStateCheckFunc(1),
			},

			{
				Config: testAccSiteBasicAwsResolverConfiWithMasterCredentials(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.%", "12"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.access_key_id", "string1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.0.%", "4"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.0.account_id", "abc123"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.0.external_id", "3a569552-48a4-4589-990d-5f1d8e3a6a18"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.0.regions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.0.regions.0", "eu-central-1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.assumed_roles.0.role_name", "The role"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.https_proxy", ""),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.name", "AWS Resolver 1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.regions.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.regions.0", "eu-central-1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.regions.1", "eu-west-1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.resolve_with_master_credentials", "false"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.secret_access_key", "string2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.update_interval", "59"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.use_iam_role", "true"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.vpc_auto_discovery", "true"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.0.vpcs.#", "0"),
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

func testAccSiteBasicAwsResolverConfiWithMasterCredentials(context map[string]interface{}) string {
	return Nprintf(`
    resource "appgatesdp_site" "test_site" {
        name       = "%{name}"
        tags = [
          "developer",
          "api-created"
        ]
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
			aws_resolvers {
			  name = "AWS Resolver 1"
			  regions = [
				"eu-central-1",
				"eu-west-1"
			  ]
			  update_interval                 = 59
			  vpcs                            = []
			  vpc_auto_discovery              = true
			  use_iam_role                    = true
			  access_key_id                   = "string1"
			  secret_access_key               = "string2"
			  resolve_with_master_credentials = false
			  assumed_roles {
				account_id  = "abc123"
				role_name   = "The role"
				external_id = "3a569552-48a4-4589-990d-5f1d8e3a6a18"
				regions = [
				  "eu-central-1",
				]
			  }
			}
		  }
      }
    `, context)
}

func testAccSiteBasicAwsResolverConfiWithMasterCredentialsUpdated(context map[string]interface{}) string {
	return Nprintf(`
    resource "appgatesdp_site" "test_site" {
        name       = "%{name}"
        tags = [
          "developer",
          "api-created"
        ]
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
			aws_resolvers {
			  name = "AWS Resolver 1"
			  regions = [
				"eu-central-1",
				"eu-west-1"
			  ]
			  update_interval                 = 59
			  vpcs                            = []
			  vpc_auto_discovery              = true
			  use_iam_role                    = true
			  access_key_id                   = "string1"
			  secret_access_key               = "string2"
			  resolve_with_master_credentials = true # updated
			  assumed_roles {
				account_id  = "abc123"
				role_name   = "The role"
				external_id = "3a569552-48a4-4589-990d-5f1d8e3a6a18"
				regions = [
				  "eu-central-1",
				]
			  }
			}
		  }
      }
    `, context)
}

// Test for
// https://github.com/appgate/terraform-provider-appgatesdp/pull/201
// https://github.com/appgate/terraform-provider-appgatesdp/issues/203
func TestAccSiteVPNRouteVia(t *testing.T) {
	resourceName := "appgatesdp_site.default_test_site"
	rName := RandStringFromCharSet(11, CharSetAlphaNum)
	context := map[string]interface{}{
		"name": rName,
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSiteDestroy,
		Steps: []resource.TestStep{
			{

				PreConfig: func() {
					// TODO
					// test includes
					// name_resolution.0.dns_resolvers.0.search_domains.#
					// which has been removed in 6.2
					applianceConstraintCheck(t, ">= 5.5, < 6.2")
				},
				Config: testAccSiteVPNRouteVia(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v4", "true"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v6", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.excluded_subnets.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "entitlement_based_routing", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_mappings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", context["name"].(string)),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.name", "test"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.search_domains.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.search_domains.0", "example.com"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.servers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.servers.0", "10.10.10.10"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.update_interval", "60"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.esx_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.gcp_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.use_hosts_file", "false"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "short_name", "DT"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "default_test_site"),
					resource.TestCheckResourceAttr(resourceName, "tags.2", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.ip_access_log_interval_seconds", "120"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.0.ipv4", "10.10.10.10"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.0.ipv6", "fdf8:f53b:82e4::53"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.snat", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.enabled", "true"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccSiteImportStateCheckFunc(1),
			},
			{

				Config: testAccSiteVPNRouteViaUpdatedV4Route(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v4", "true"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v6", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.excluded_subnets.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "entitlement_based_routing", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_mappings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", context["name"].(string)),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.name", "test"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.search_domains.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.search_domains.0", "example.com"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.servers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.servers.0", "10.10.10.10"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.update_interval", "60"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.esx_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.gcp_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.use_hosts_file", "false"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "short_name", "DT"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "default_test_site"),
					resource.TestCheckResourceAttr(resourceName, "tags.2", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.ip_access_log_interval_seconds", "120"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.0.ipv4", "20.20.20.20"), // updated
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.0.ipv6", "fdf8:f53b:82e4::53"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.snat", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.enabled", "true"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccSiteImportStateCheckFunc(1),
			},
			{

				Config: testAccSiteVPNRouteViaDeleted(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v4", "true"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v6", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.excluded_subnets.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "entitlement_based_routing", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_mappings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", context["name"].(string)),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.name", "test"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.search_domains.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.search_domains.0", "example.com"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.servers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.servers.0", "10.10.10.10"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.update_interval", "60"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.esx_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.gcp_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.use_hosts_file", "false"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "short_name", "DT"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "default_test_site"),
					resource.TestCheckResourceAttr(resourceName, "tags.2", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.ip_access_log_interval_seconds", "120"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.#", "0"), // deleted
					resource.TestCheckResourceAttr(resourceName, "vpn.0.snat", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.enabled", "true"),
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

// https://github.com/appgate/terraform-provider-appgatesdp/pull/201#issuecomment-1001642495
func testAccSiteVPNRouteVia(context map[string]interface{}) string {
	return Nprintf(`
    resource "appgatesdp_site" "default_test_site" {
        name                      = "%{name}"
        short_name                = "DT"
        entitlement_based_routing = false
        notes                     = "Managed by terraform"
        default_gateway {
          enabled_v4 = true
          enabled_v6 = false
        }
        name_resolution {
          use_hosts_file = false
          dns_resolvers {
            name            = "test"
            update_interval = 60
            servers         = ["10.10.10.10"]
            search_domains  = ["example.com"]
          }
        }
        vpn {
          ip_access_log_interval_seconds = 120
          snat                           = false
          tls {
            enabled = true
          }
          dtls {
            enabled = false
          }
          route_via {
            ipv4 = "10.10.10.10"
            ipv6 = "fdf8:f53b:82e4::53"
          }
        }
        tags = ["terraform", "api-created", "default_test_site"]
      }
`, context)
}

// https://github.com/appgate/terraform-provider-appgatesdp/pull/201#issuecomment-1001642495
func testAccSiteVPNRouteViaUpdatedV4Route(context map[string]interface{}) string {
	return Nprintf(`
    resource "appgatesdp_site" "default_test_site" {
        name                      = "%{name}"
        short_name                = "DT"
        entitlement_based_routing = false
        notes                     = "Managed by terraform"
        default_gateway {
          enabled_v4 = true
          enabled_v6 = false
        }
        name_resolution {
          use_hosts_file = false
          dns_resolvers {
            name            = "test"
            update_interval = 60
            servers         = ["10.10.10.10"]
            search_domains  = ["example.com"]
          }
        }
        vpn {
          ip_access_log_interval_seconds = 120
          snat                           = false
          tls {
            enabled = true
          }
          dtls {
            enabled = false
          }
          route_via {
            ipv4 = "20.20.20.20"
            ipv6 = "fdf8:f53b:82e4::53"
          }
        }
        tags = ["terraform", "api-created", "default_test_site"]
      }
`, context)
}

// https://github.com/appgate/terraform-provider-appgatesdp/pull/201#issuecomment-1001642495
func testAccSiteVPNRouteViaDeleted(context map[string]interface{}) string {
	return Nprintf(`
    resource "appgatesdp_site" "default_test_site" {
        name                      = "%{name}"
        short_name                = "DT"
        entitlement_based_routing = false
        notes                     = "Managed by terraform"
        default_gateway {
          enabled_v4 = true
          enabled_v6 = false
        }
        name_resolution {
          use_hosts_file = false
          dns_resolvers {
            name            = "test"
            update_interval = 60
            servers         = ["10.10.10.10"]
            search_domains  = ["example.com"]
          }
        }
        vpn {
          ip_access_log_interval_seconds = 120
          snat                           = false
          tls {
            enabled = true
          }
          dtls {
            enabled = false
          }
        }
        tags = ["terraform", "api-created", "default_test_site"]
      }
`, context)
}

// Test for
// https://github.com/appgate/terraform-provider-appgatesdp/issues/204
func TestAccSiteVPNRouteViaIpv4Only(t *testing.T) {
	resourceName := "appgatesdp_site.d_test_site"
	rName := RandStringFromCharSet(11, CharSetAlphaNum)
	context := map[string]interface{}{
		"name": rName,
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSiteDestroy,
		Steps: []resource.TestStep{
			{

				Config: testAccSiteVPNRouteViaIpv4Only(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v4", "true"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v6", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.excluded_subnets.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "entitlement_based_routing", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_mappings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", context["name"].(string)),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.esx_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.gcp_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.use_hosts_file", "false"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "short_name", "DT"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "default_test"),
					resource.TestCheckResourceAttr(resourceName, "tags.2", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.ip_access_log_interval_seconds", "120"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.0.ipv4", "10.10.10.10"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.0.ipv6", ""),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.snat", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.enabled", "true"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccSiteImportStateCheckFunc(1),
			},
			{

				Config: testAccSiteVPNRouteViaIpv4OnlyUpdated(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v4", "true"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v6", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.excluded_subnets.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "entitlement_based_routing", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_mappings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", context["name"].(string)),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.esx_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.gcp_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.use_hosts_file", "false"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "short_name", "DT"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "default_test"),
					resource.TestCheckResourceAttr(resourceName, "tags.2", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.ip_access_log_interval_seconds", "120"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.0.ipv4", "10.20.10.20"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.0.ipv6", ""),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.snat", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.enabled", "true"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccSiteImportStateCheckFunc(1),
			},
			{

				Config: testAccSiteVPNRouteViaIpv4OnlyUpdatedWithIpv6(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v4", "true"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v6", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.excluded_subnets.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "entitlement_based_routing", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_mappings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", context["name"].(string)),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.esx_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.gcp_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.use_hosts_file", "false"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "short_name", "DT"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "default_test"),
					resource.TestCheckResourceAttr(resourceName, "tags.2", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.ip_access_log_interval_seconds", "120"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.0.ipv6", "fdf8:f53b:82e4::53"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.snat", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.enabled", "true"),
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

// https://github.com/appgate/terraform-provider-appgatesdp/issues/204
func testAccSiteVPNRouteViaIpv4Only(context map[string]interface{}) string {
	return Nprintf(`
resource "appgatesdp_site" "d_test_site" {
	name                      = "%{name}"
	short_name                = "DT"
	entitlement_based_routing = false
	notes                     = "Managed by terraform"

	default_gateway {
	  enabled_v4 = true
	  enabled_v6 = false
	}
	vpn {
	  ip_access_log_interval_seconds = 120
	  snat                           = false
	  tls {
		enabled = true
	  }
	  dtls {
		enabled = false
	  }
	  route_via {
		ipv4 = "10.10.10.10"
	  }
	}
	tags = ["terraform", "api-created", "default_test"]
}
`, context)
}

// https://github.com/appgate/terraform-provider-appgatesdp/issues/204
func testAccSiteVPNRouteViaIpv4OnlyUpdated(context map[string]interface{}) string {
	return Nprintf(`
resource "appgatesdp_site" "d_test_site" {
	name                      = "%{name}"
	short_name                = "DT"
	entitlement_based_routing = false
	notes                     = "Managed by terraform"

	default_gateway {
	  enabled_v4 = true
	  enabled_v6 = false
	}
	vpn {
	  ip_access_log_interval_seconds = 120
	  snat                           = false
	  tls {
		enabled = true
	  }
	  dtls {
		enabled = false
	  }
	  route_via {
		ipv4 = "10.20.10.20"
	  }
	}
	tags = ["terraform", "api-created", "default_test"]
}
`, context)
}

// https://github.com/appgate/terraform-provider-appgatesdp/issues/204
func testAccSiteVPNRouteViaIpv4OnlyUpdatedWithIpv6(context map[string]interface{}) string {
	return Nprintf(`
resource "appgatesdp_site" "d_test_site" {
	name                      = "%{name}"
	short_name                = "DT"
	entitlement_based_routing = false
	notes                     = "Managed by terraform"

	default_gateway {
	  enabled_v4 = true
	  enabled_v6 = false
	}
	vpn {
	  ip_access_log_interval_seconds = 120
	  snat                           = false
	  tls {
		enabled = true
	  }
	  dtls {
		enabled = false
	  }
	  route_via {
		ipv6 = "fdf8:f53b:82e4::53"
	  }
	}
	tags = ["terraform", "api-created", "default_test"]
}
`, context)
}

func TestAccSiteNameResolver6(t *testing.T) {
	resourceName := "appgatesdp_site.test_site"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					// test includes
					// name_resolution.dns_resolvers.search_domains
					// which has been removed in 6.2
					applianceConstraintCheck(t, ">= 6.0, < 6.2")
				},
				Config: testAccSiteNameResolver6(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),

					resource.TestCheckResourceAttr(resourceName, "default_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v4", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v6", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.excluded_subnets.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "entitlement_based_routing", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_mappings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.0.address", "0.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.0.netmask", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.1.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.1.address", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.1.netmask", "32"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.2.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.2.address", "::"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.2.netmask", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.default_ttl_seconds", "15"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.dns_servers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.dns_servers.0", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.site_ipv4", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.site_ipv6", "2001:0db8:85a3:0000:0000:8a2e:0370:7334"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.%", "6"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.default_ttl_seconds", "99"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.name", "DNS Resolver 1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.query_aaaa", "true"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.search_domains.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.search_domains.0", "hostname.dns"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.search_domains.1", "foo.bar"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.servers.0", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.servers.1", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.update_interval", "13"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.esx_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.gcp_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.use_hosts_file", "false"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "notes", "This object has been created for test purposes."),
					resource.TestCheckResourceAttr(resourceName, "short_name", "ts0"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "developer"),
					resource.TestCheckResourceAttr(resourceName, "vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.ip_access_log_interval_seconds", "120"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.snat", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.enabled", "true"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccSiteImportStateCheckFunc(1),
			},
			{
				Config: testAccSiteNameResolver6Updated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.0.address", "0.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.0.netmask", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.1.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.1.address", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.1.netmask", "32"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.2.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.2.address", "::"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.2.netmask", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.default_ttl_seconds", "43199"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.dns_servers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.dns_servers.0", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.site_ipv4", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.site_ipv6", "2001:0db8:85a3:0000:0000:8a2e:0370:7334"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.%", "6"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.default_ttl_seconds", "5"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.name", "DNS Resolver 1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.query_aaaa", "true"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.search_domains.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.search_domains.0", "hostname.dns"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.search_domains.1", "foo.bar"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.servers.0", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.servers.1", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.0.update_interval", "13"),
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

func testAccSiteNameResolver6(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_site" "test_site" {
	name       = "%s"
	short_name = "ts0"
	tags = [
	  "developer",
	  "api-created"
	]

	notes                     = "This object has been created for test purposes."
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
		name                = "DNS Resolver 1"
		update_interval     = 13
		query_aaaa          = true
		default_ttl_seconds = 99
		servers = [
		  "8.8.8.8",
		  "1.1.1.1"
		]
		search_domains = [
		  "hostname.dns",
		  "foo.bar"
		]
	  }

	  dns_forwarding {
		site_ipv4 = "1.2.3.4"
		site_ipv6 = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
		dns_servers = [
		  "1.1.1.1"
		]
		default_ttl_seconds = 15
		allow_destinations {
		  address = "1.1.1.1"
		  netmask = 32
		}
		allow_destinations {
		  address = "0.0.0.0"
		  netmask = 0
		}
		allow_destinations {
		  address = "::"
		  netmask = 0
		}
	  }
	}
}`, rName)
}

func testAccSiteNameResolver6Updated(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_site" "test_site" {
	name       = "%s"
	short_name = "ts0"
	tags = [
	  "developer",
	  "api-created"
	]

	notes                     = "This object has been created for test purposes."
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
		name                = "DNS Resolver 1"
		update_interval     = 13
		query_aaaa          = true
		default_ttl_seconds = 5
		servers = [
		  "8.8.8.8",
		  "1.1.1.1"
		]
		search_domains = [
		  "hostname.dns",
		  "foo.bar"
		]
	  }

	  dns_forwarding {
		site_ipv4 = "1.2.3.4"
		site_ipv6 = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
		dns_servers = [
		  "1.1.1.1"
		]
		default_ttl_seconds = 43199
		allow_destinations {
		  address = "1.1.1.1"
		  netmask = 32
		}
		allow_destinations {
		  address = "0.0.0.0"
		  netmask = 0
		}
		allow_destinations {
		  address = "::"
		  netmask = 0
		}
	  }
	}
}`, rName)
}

func TestAccSiteNameResolverIllumio61(t *testing.T) {
	resourceName := "appgatesdp_site.illumio_site"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					applianceConstraintCheck(t, ">= 6.1, < 6.2")
				},
				Config: testAccSiteNameResolverIllumio(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v4", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v6", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.excluded_subnets.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "entitlement_based_routing", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_mappings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.esx_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.gcp_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.%", "7"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.hostname", "illumio.acme.com"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.name", "Illumio Resolver 1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.password", "adminadmin"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.port", "65530"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.update_interval", "10"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.username", "admin"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.use_hosts_file", "false"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "short_name", "ts0"),
					resource.TestCheckResourceAttr(resourceName, "vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.ip_access_log_interval_seconds", "120"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.snat", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.enabled", "true"),
				),
				ImportStateVerifyIgnore: []string{
					"name_resolution.0.illumio_resolvers.0.org_id",
				},
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccSiteImportStateCheckFunc(1),
				ImportStateVerifyIgnore: []string{
					"name_resolution.0.illumio_resolvers.0.org_id",
				},
			},
			{
				Config: testAccSiteNameResolverIllumioUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v4", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v6", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.excluded_subnets.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "entitlement_based_routing", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_mappings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.esx_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.gcp_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.%", "7"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.hostname", "illumio.acme.com"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.name", "Illumio Resolver 99"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.password", "adminadmin"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.port", "1337"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.update_interval", "50"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.username", "acme"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.use_hosts_file", "false"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "short_name", "ts0"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.ip_access_log_interval_seconds", "120"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.snat", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.enabled", "true"),
				),
				ImportStateVerifyIgnore: []string{
					"name_resolution.0.illumio_resolvers.0.org_id",
				},
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccSiteImportStateCheckFunc(1),
			},
			{
				Config: testAccSiteNameResolverIllumioRemoved(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v4", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v6", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.excluded_subnets.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "entitlement_based_routing", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_mappings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.esx_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.gcp_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.use_hosts_file", "false"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "short_name", "ts0"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.ip_access_log_interval_seconds", "120"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.snat", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.enabled", "true"),
				),
				ImportStateVerifyIgnore: []string{
					"name_resolution.0.illumio_resolvers.0.org_id",
				},
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccSiteImportStateCheckFunc(1),
			},
		},
	})
}

func testAccSiteNameResolverIllumio(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_site" "illumio_site" {
	name                      = "%s"
	short_name                = "ts0"
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
	  illumio_resolvers {
		name     = "Illumio Resolver 1"
		hostname = "illumio.acme.com"
		update_interval = 10
		port     = 65530
		username = "admin"
		password = "adminadmin"
	  }
	}
}`, rName)
}

func testAccSiteNameResolverIllumioUpdated(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_site" "illumio_site" {
	name                      = "%s"
	short_name                = "ts0"
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
	  illumio_resolvers {
		name     = "Illumio Resolver 99"
		hostname = "illumio.acme.com"
		update_interval = 50
		port     = 1337
		username = "acme"
		password = "adminadmin"
	  }
	}
}`, rName)
}

func testAccSiteNameResolverIllumioRemoved(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_site" "illumio_site" {
	name                      = "%s"
	short_name                = "ts0"
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
	}
}`, rName)
}

func TestAccSiteNameResolverIllumio62(t *testing.T) {
	resourceName := "appgatesdp_site.illumio_site"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testFor62AndAbove(t)
				},
				Config: testAccSiteNameResolverIllumio62(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v4", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v6", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.excluded_subnets.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "entitlement_based_routing", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_mappings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.esx_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.gcp_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.%", "7"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.hostname", "illumio.acme.com"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.name", "Illumio Resolver 1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.password", "adminadmin"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.port", "65530"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.update_interval", "10"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.username", "admin"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.org_id", "org12345"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.use_hosts_file", "false"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "short_name", "ts0"),
					resource.TestCheckResourceAttr(resourceName, "vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.ip_access_log_interval_seconds", "120"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.snat", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.enabled", "true"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccSiteImportStateCheckFunc(1),
			},
			{
				Config: testAccSiteNameResolverIllumioUpdated62(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v4", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v6", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.excluded_subnets.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "entitlement_based_routing", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_mappings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.esx_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.gcp_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.%", "7"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.hostname", "illumio.acme.com"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.name", "Illumio Resolver 99"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.password", "adminadmin"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.port", "1337"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.update_interval", "50"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.username", "acme"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.0.org_id", "org12345"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.use_hosts_file", "false"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "short_name", "ts0"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.ip_access_log_interval_seconds", "120"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.snat", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.enabled", "true"),
				),
			},
			{
				ResourceName:     resourceName,
				ImportState:      true,
				ImportStateCheck: testAccSiteImportStateCheckFunc(1),
			},
			{
				Config: testAccSiteNameResolverIllumioRemoved(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v4", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v6", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.excluded_subnets.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "entitlement_based_routing", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_mappings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.aws_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.azure_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.esx_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.gcp_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.illumio_resolvers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.use_hosts_file", "false"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "short_name", "ts0"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.%", "5"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.dtls.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.ip_access_log_interval_seconds", "120"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.route_via.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.snat", "false"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.tls.0.enabled", "true"),
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

func testAccSiteNameResolverIllumio62(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_site" "illumio_site" {
	name                      = "%s"
	short_name                = "ts0"
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
	  illumio_resolvers {
		name     = "Illumio Resolver 1"
		hostname = "illumio.acme.com"
		update_interval = 10
		port     = 65530
		username = "admin"
		password = "adminadmin"
		org_id = "org12345"
	  }
	}
}`, rName)
}

func testAccSiteNameResolverIllumioUpdated62(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_site" "illumio_site" {
	name                      = "%s"
	short_name                = "ts0"
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
	  illumio_resolvers {
		name     = "Illumio Resolver 99"
		hostname = "illumio.acme.com"
		update_interval = 50
		port     = 1337
		username = "acme"
		password = "adminadmin"
		org_id = "org12345"
	  }
	}
}`, rName)
}

func TestAccSiteBasic2(t *testing.T) {
	resourceName := "appgatesdp_site.test_site"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSite2(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v4", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v6", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.excluded_subnets.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "entitlement_based_routing", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_mappings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "notes", "This object has been created for test purposes."),
					resource.TestCheckResourceAttr(resourceName, "vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.snat", "false"),
				),
			},
			{
				Config: testAccCheckSiteUpdate2(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "The test site"),
					resource.TestCheckResourceAttr(resourceName, "short_name", "ts1"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.1", "10.20.0.0/24"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "developer"),
					resource.TestCheckResourceAttr(resourceName, "tags.2", "updated-tag"),
					resource.TestCheckResourceAttr(resourceName, "notes", "note updated"),
					resource.TestCheckResourceAttr(resourceName, "vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.snat", "true"),
				),
			},
		},
	})
}

func testAccCheckSite2(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_site" "test_site" {
    name       = "%s"
    network_subnets = [
        "10.0.0.0/16"
    ]
    notes = "This object has been created for test purposes."
	vpn {
	  snat          = false
	}
}
`, rName)
}

func testAccCheckSiteUpdate2() string {
	return `
resource "appgatesdp_site" "test_site" {
    name       = "The test site"
    short_name = "ts1"
    tags = [
        "developer",
        "api-created",
	    "updated-tag"
    ]
    notes = "note updated"
    entitlement_based_routing = false
	network_subnets = [
        "10.0.0.0/16",
        "10.20.0.0/24",
    ]
	vpn {
		snat          = true
	}
}
`
}

func TestAccSiteBasic3(t *testing.T) {
	resourceName := "appgatesdp_site.test_site"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSite3(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v4", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v6", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.excluded_subnets.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "entitlement_based_routing", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_mappings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "notes", "This object has been created for test purposes."),
					resource.TestCheckResourceAttr(resourceName, "vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.snat", "false"),
				),
			},
			{
				Config: testAccCheckSite3Updated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v4", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.enabled_v6", "false"),
					resource.TestCheckResourceAttr(resourceName, "default_gateway.0.excluded_subnets.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "entitlement_based_routing", "false"),
					resource.TestCheckResourceAttr(resourceName, "ip_pool_mappings.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "network_subnets.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "notes", "This object has been created for test purposes."),
					resource.TestCheckResourceAttr(resourceName, "vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "vpn.0.snat", "false"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.0.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.0.address", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.allow_destinations.0.netmask", "32"),
					resource.TestCheckResourceAttr(resourceName, "name_resolution.0.dns_forwarding.0.site_ipv4", "192.168.1.1"),
				),
			},
		},
	})
}

func testAccCheckSite3(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_site" "test_site" {
    name       = "%s"
    network_subnets = [
        "10.0.0.0/16"
    ]
    notes = "This object has been created for test purposes."
	vpn {
	  snat          = false
	}
}
`, rName)
}

func testAccCheckSite3Updated(rName string) string {
	return fmt.Sprintf(`
resource "appgatesdp_site" "test_site" {
    name       = "%s"
    network_subnets = [
        "10.0.0.0/16"
    ]
    notes = "This object has been created for test purposes."
	vpn {
	    snat          = false
	}
	name_resolution {
		dns_forwarding {
            default_ttl_seconds = 300
			site_ipv4 = "192.168.1.1"
			dns_servers = [
			    "1.1.1.1"
			]
			allow_destinations {
			    address = "1.1.1.1"
			    netmask = 32
			}
		}
	}
}
`, rName)
}
