package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccApplianceBasicController(t *testing.T) {
	resourceName := "appgate_appliance.test_controller"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplianceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApplianceBasicController(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplianceExists(resourceName),
					// testAccCheckExampleWidgetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "controller-test"),
					resource.TestCheckResourceAttr(resourceName, "hostname", "envy-10-97-168-1337.devops"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),

					resource.TestCheckResourceAttr(resourceName, "controller.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "controller.2223495429.enabled", "true"),

					resource.TestCheckResourceAttr(resourceName, "gateway.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "log_server.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "client_interface.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.3426664317.allow_sources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.3426664317.dtls_port", "445"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.3426664317.https_port", "444"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.3426664317.proxy_protocol", "true"),

					resource.TestCheckResourceAttr(resourceName, "networking.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.dns_domains.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.dns_domains.112524683", "aa.com"),

					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.dns_servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.dns_servers.251826590", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.dns_servers.2609393598", "8.8.8.8"),

					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.nics.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.nics.0.name", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.nics.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.nics.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.nics.0.ipv4.3519857096.dhcp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.nics.0.ipv4.3519857096.dhcp.2319808068.dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.nics.0.ipv4.3519857096.dhcp.2319808068.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.nics.0.ipv4.3519857096.dhcp.2319808068.ntp", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.nics.0.ipv4.3519857096.dhcp.2319808068.routers", "true"),

					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.routes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.routes.0.address", "0.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.routes.0.gateway", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.routes.0.netmask", "24"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.routes.0.nic", "eth0"),

					resource.TestCheckResourceAttr(resourceName, "snmp_server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.10952821.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.10952821.snmpd_conf", "foo"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.10952821.tcp_port", "161"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.10952821.udp_port", "161"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.10952821.allow_sources.0.address", "1.3.3.7"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.10952821.allow_sources.0.netmask", "0"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.10952821.allow_sources.0.nic", "eth0"),

					resource.TestCheckResourceAttr(resourceName, "healthcheck_server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "healthcheck_server.1825715455.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "healthcheck_server.1825715455.port", "5555"),
					resource.TestCheckResourceAttr(resourceName, "healthcheck_server.1825715455.allow_sources.0.address", "1.3.3.7"),
					resource.TestCheckResourceAttr(resourceName, "healthcheck_server.1825715455.allow_sources.0.netmask", "0"),
					resource.TestCheckResourceAttr(resourceName, "healthcheck_server.1825715455.allow_sources.0.nic", "eth0"),

					resource.TestCheckResourceAttr(resourceName, "prometheus_exporter.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "prometheus_exporter.420734837.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "prometheus_exporter.420734837.port", "1234"),
					resource.TestCheckResourceAttr(resourceName, "prometheus_exporter.420734837.allow_sources.0.address", "1.3.3.7"),
					resource.TestCheckResourceAttr(resourceName, "prometheus_exporter.420734837.allow_sources.0.netmask", "0"),
					resource.TestCheckResourceAttr(resourceName, "prometheus_exporter.420734837.allow_sources.0.nic", "eth0"),

					resource.TestCheckResourceAttr(resourceName, "ping.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ping.1092164059.allow_sources.0.address", "1.3.3.7"),
					resource.TestCheckResourceAttr(resourceName, "ping.1092164059.allow_sources.0.netmask", "0"),
					resource.TestCheckResourceAttr(resourceName, "ping.1092164059.allow_sources.0.nic", "eth0"),

					resource.TestCheckResourceAttr(resourceName, "log_forwarder.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.3792789776.elasticsearch.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.3792789776.elasticsearch.1193666681.aws_id", "string"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.3792789776.elasticsearch.1193666681.aws_region", "eu-west-2"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.3792789776.elasticsearch.1193666681.aws_secret", ""),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.3792789776.elasticsearch.1193666681.retention_days", "3"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.3792789776.elasticsearch.1193666681.url", "https://aws.com/elasticsearch/instance/asdaxllkmda64"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.3792789776.elasticsearch.1193666681.use_instance_credentials", "true"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.3792789776.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.3792789776.sites.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.3792789776.sites.2312403857", "8a4add9e-0e99-4bb1-949c-c9faf9a49ad4"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.3792789776.tcp_clients.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.3792789776.tcp_clients.0.format", "json"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.3792789776.tcp_clients.0.host", "siem.company.com"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.3792789776.tcp_clients.0.name", "Company SIEM"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.3792789776.tcp_clients.0.port", "8888"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.3792789776.tcp_clients.0.use_tls", "true"),

					resource.TestCheckResourceAttr(resourceName, "ntp.3633490038.servers.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "ntp.3633490038.servers.0.hostname", "0.ubuntu.pool.ntp.org"),
					resource.TestCheckResourceAttr(resourceName, "ntp.3633490038.servers.0.key", ""),
					resource.TestCheckResourceAttr(resourceName, "ntp.3633490038.servers.0.key_type", ""),
					resource.TestCheckResourceAttr(resourceName, "ntp.3633490038.servers.1.hostname", "1.ubuntu.pool.ntp.org"),
					resource.TestCheckResourceAttr(resourceName, "ntp.3633490038.servers.1.key", ""),
					resource.TestCheckResourceAttr(resourceName, "ntp.3633490038.servers.1.key_type", ""),
					resource.TestCheckResourceAttr(resourceName, "ntp.3633490038.servers.2.hostname", "2.ubuntu.pool.ntp.org"),
					resource.TestCheckResourceAttr(resourceName, "ntp.3633490038.servers.2.key", ""),
					resource.TestCheckResourceAttr(resourceName, "ntp.3633490038.servers.2.key_type", ""),
					resource.TestCheckResourceAttr(resourceName, "ntp.3633490038.servers.3.hostname", "3.ubuntu.pool.ntp.org"),
					resource.TestCheckResourceAttr(resourceName, "ntp.3633490038.servers.3.key", ""),
					resource.TestCheckResourceAttr(resourceName, "ntp.3633490038.servers.3.key_type", ""),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				// ImportStateVerify: true,
				ImportStateCheck: testAccApplianceImportStateCheckFunc(1),
			},
			{
				Config: testAccCheckApplianceBasicControllerUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "networking.0.dns_domains.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.dns_servers.#", "0"),
				),
			},
		},
	})
}

func testAccCheckExampleWidgetExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Widget ID is not set")
		}
		fmt.Println("=====")
		fmt.Println("=====")
		fmt.Println("=====")
		fmt.Printf("\n resourceName: %+v\n", resourceName)
		fmt.Println("=====")
		fmt.Printf("\nstate: %+v\n", s)
		fmt.Println("=====")
		fmt.Printf("\nresource: %+v\n", rs)
		fmt.Println("=====")

		return nil
	}
}
func testAccApplianceImportStateCheckFunc(expectedStates int) resource.ImportStateCheckFunc {
	return func(s []*terraform.InstanceState) error {
		if len(s) != expectedStates {
			return fmt.Errorf("expected %d states, got %d: %+v", expectedStates, len(s), s)
		}
		return nil
	}
}

func testAccCheckApplianceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "appgate_appliance" {
			continue
		}

		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.AppliancesApi

		_, _, err := api.AppliancesIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err == nil {
			return fmt.Errorf("Appliance still exists, %+v", err)
		}
	}
	return nil
}

func testAccCheckApplianceBasicController() string {
	return fmt.Sprintf(`
data "appgate_site" "default_site" {
	   site_name = "Default site"
}

resource "appgate_appliance" "test_controller" {
	name     = "controller-test"
	hostname = "envy-10-97-168-1337.devops"

	client_interface {
        hostname = "envy-10-97-168-1337.devops"
        proxy_protocol = true
        https_port     = 444
        dtls_port      = 445
        allow_sources {
          address = "1.3.3.7"
          netmask = 0
          nic     = "eth0"
        }
        override_spa_mode = "TCP"
	}

	peer_interface {
		hostname   = "envy-10-97-168-1337.devops"
		https_port = "1337"
	}
	tags = [
		"terraform",
		"api-test-created"
	]
	ntp {
		servers {
			hostname = "0.ubuntu.pool.ntp.org"
		}
		servers {
			hostname = "1.ubuntu.pool.ntp.org"
		}
		servers {
			hostname = "2.ubuntu.pool.ntp.org"
		}
		servers {
			hostname = "3.ubuntu.pool.ntp.org"
		}
	}
	networking {

		hosts {
		  hostname = "bla"
		  address  = "0.0.0.0"
		}

		nics {
		  enabled = true
		  name    = "eth0"
		  ipv4 {
			dhcp {
			  enabled = true
			  dns     = true
			  routers = true
			  ntp     = true
			}
		  }
		}
		dns_servers = [
		  "8.8.8.8",
		  "1.1.1.1",
		]
		dns_domains = [
		  "aa.com"
        ]
        routes {
            address = "0.0.0.0"
            netmask = 24
            gateway = "1.2.3.4"
            nic = "eth0"
        }
	}
	controller {
		enabled = true
    }
    snmp_server {
        enabled    = true
        tcp_port   = 161
        udp_port   = 161
        snmpd_conf = "foo"
        allow_sources {
          address = "1.3.3.7"
          netmask = 0
          nic     = "eth0"
        }
    }
    healthcheck_server {
        enabled = true
        port    = 5555
        allow_sources {
          address = "1.3.3.7"
          netmask = 0
          nic     = "eth0"
        }
    }
    prometheus_exporter {
        enabled = true
        port    = 1234
        allow_sources {
          address = "1.3.3.7"
          netmask = 0
          nic     = "eth0"
        }
    }
    ping {
        allow_sources {
          address = "1.3.3.7"
          netmask = 0
          nic     = "eth0"
        }
    }
    log_forwarder {
        enabled = true
        elasticsearch {
            url = "https://aws.com/elasticsearch/instance/asdaxllkmda64"
            aws_id = "string"
            aws_region = "eu-west-2"
            use_instance_credentials = true
            retention_days = 3
        }

        tcp_clients {
            name = "Company SIEM"
            host = "siem.company.com"
            port = 8888
            format = "json"
            use_tls = true
        }
        sites = [
            data.appgate_site.default_site.id
        ]
    }
}
`)
}
func testAccCheckApplianceBasicControllerUpdate() string {
	return fmt.Sprintf(`
data "appgate_site" "default_site" {
	   site_name = "Default site"
}

resource "appgate_appliance" "test_controller" {
	name     = "controller-test"
	hostname = "envy-10-97-168-1337.devops"

	client_interface {
        hostname = "envy-10-97-168-1337.devops"
        proxy_protocol = true
        https_port     = 444
        dtls_port      = 445
        allow_sources {
          address = "1.3.3.7"
          netmask = 0
          nic     = "eth0"
        }
        override_spa_mode = "TCP"
	}

	peer_interface {
		hostname   = "envy-10-97-168-1337.devops"
		https_port = "1337"
	}
	tags = [
		"terraform",
		"api-test-created"
	]
	networking {

		hosts {
		  hostname = "bla"
		  address  = "0.0.0.0"
		}

		nics {
		  enabled = true
		  name    = "eth0"
		  ipv4 {
			dhcp {
			  enabled = true
			  dns     = true
			  routers = true
			  ntp     = true
			}
		  }
		}
	}
	controller {
		enabled = true
    }
}
`)
}

func TestAccApplianceIoTConnector(t *testing.T) {
	resourceName := "appgate_appliance.iot_connector"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplianceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApplianceBasicIotConnector(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplianceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "iot-connector-test"),
					resource.TestCheckResourceAttr(resourceName, "hostname", "envy-10-97-168-1234.devops"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),

					resource.TestCheckResourceAttr(resourceName, "controller.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "gateway.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "log_server.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "client_interface.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.3170980607.allow_sources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.3170980607.dtls_port", "445"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.3170980607.https_port", "444"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.3170980607.proxy_protocol", "true"),

					resource.TestCheckResourceAttr(resourceName, "networking.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.dns_domains.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.dns_domains.112524683", "aa.com"),

					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.dns_servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.dns_servers.251826590", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.dns_servers.2609393598", "8.8.8.8"),

					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.nics.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.nics.0.name", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.nics.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.nics.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.nics.0.ipv4.3519857096.dhcp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.nics.0.ipv4.3519857096.dhcp.2319808068.dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.nics.0.ipv4.3519857096.dhcp.2319808068.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.nics.0.ipv4.3519857096.dhcp.2319808068.ntp", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.nics.0.ipv4.3519857096.dhcp.2319808068.routers", "true"),

					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.routes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.routes.0.address", "0.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.routes.0.gateway", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.routes.0.netmask", "24"),
					resource.TestCheckResourceAttr(resourceName, "networking.1914549515.routes.0.nic", "eth0"),

					resource.TestCheckResourceAttr(resourceName, "iot_connector.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "iot_connector.1446797058.clients.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "iot_connector.1446797058.clients.0.device_id", "12699e27-b584-464a-81ee-5b4784b6d425"),
					resource.TestCheckResourceAttr(resourceName, "iot_connector.1446797058.clients.0.name", "Printers"),
					resource.TestCheckResourceAttr(resourceName, "iot_connector.1446797058.clients.0.snat", "true"),
					resource.TestCheckResourceAttr(resourceName, "iot_connector.1446797058.clients.0.sources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "iot_connector.1446797058.clients.0.sources.0.address", "1.3.3.7"),
					resource.TestCheckResourceAttr(resourceName, "iot_connector.1446797058.clients.0.sources.0.netmask", "24"),
					resource.TestCheckResourceAttr(resourceName, "iot_connector.1446797058.clients.0.sources.0.nic", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "iot_connector.1446797058.enabled", "true"),

					resource.TestCheckResourceAttr(resourceName, "rsyslog_destinations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "rsyslog_destinations.0.destination", "10.10.10.2"),
					resource.TestCheckResourceAttr(resourceName, "rsyslog_destinations.0.selector", "*.*"),
					resource.TestCheckResourceAttr(resourceName, "rsyslog_destinations.0.template", "hostname"),

					resource.TestCheckResourceAttr(resourceName, "hostname_aliases.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "hostname_aliases.0", "appgatealias.company.com"),
					resource.TestCheckResourceAttr(resourceName, "hostname_aliases.1", "alias2.appgate.company.com"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				// ImportStateVerify: true,
				ImportStateCheck: testAccApplianceImportStateCheckFunc(1),
			},
		},
	})
}
func testAccCheckApplianceBasicIotConnector() string {
	return fmt.Sprintf(`
data "appgate_site" "default_site" {
	   site_name = "Default site"
}

resource "appgate_appliance" "iot_connector" {
	name     = "iot-connector-test"
	hostname = "envy-10-97-168-1234.devops"

	client_interface {
        hostname = "envy-10-97-168-1234.devops"
        proxy_protocol = true
        https_port     = 444
        dtls_port      = 445
        allow_sources {
          address = "1.3.3.7"
          netmask = 0
          nic     = "eth0"
        }
        override_spa_mode = "TCP"
	}

	peer_interface {
		hostname   = "envy-10-97-168-1234.devops"
		https_port = "1337"
	}
	tags = [
		"terraform",
		"api-test-created"
	]
	networking {

		hosts {
		  hostname = "bla"
		  address  = "0.0.0.0"
		}

		nics {
		  enabled = true
		  name    = "eth0"
		  ipv4 {
			dhcp {
			  enabled = true
			  dns     = true
			  routers = true
			  ntp     = true
			}
		  }
		}
		dns_servers = [
		  "8.8.8.8",
		  "1.1.1.1",
		]
		dns_domains = [
		  "aa.com"
        ]
        routes {
            address = "0.0.0.0"
            netmask = 24
            gateway = "1.2.3.4"
            nic = "eth0"
        }
	}
    iot_connector {
        enabled = true
        clients {
          name      = "Printers"
          device_id = "12699e27-b584-464a-81ee-5b4784b6d425"

          sources {
            address = "1.3.3.7"
            netmask = 24
            nic     = "eth0"
          }
          snat = true
        }
    }

    rsyslog_destinations {
        selector    = "*.*"
        template    = "hostname"
        destination = "10.10.10.2"
    }
    hostname_aliases = [
        "appgatealias.company.com",
        "alias2.appgate.company.com"
    ]
}
`)
}

func testAccCheckApplianceExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		token := testAccProvider.Meta().(*Client).Token
		api := testAccProvider.Meta().(*Client).API.AppliancesApi

		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		_, _, err := api.AppliancesIdGet(context.Background(), rs.Primary.ID).Authorization(token).Execute()
		if err != nil {
			return fmt.Errorf("error fetching appliance with resource %s. %s", resource, err)
		}
		return nil
	}
}

func TestAccApplianceBasicGateway(t *testing.T) {
	resourceName := "appgate_appliance.test_gateway"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplianceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApplianceBasicGateway(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplianceExists(resourceName),
					// testAccCheckExampleWidgetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "gateway-test"),
					resource.TestCheckResourceAttr(resourceName, "hostname", "envy-10-97-168-1338.devops"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),

					resource.TestCheckResourceAttr(resourceName, "controller.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "gateway.1260058071.vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "gateway.1260058071.vpn.2866194231.allow_destinations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "gateway.1260058071.vpn.2866194231.allow_destinations.0.address", "0.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "gateway.1260058071.vpn.2866194231.allow_destinations.0.netmask", "0"),
					resource.TestCheckResourceAttr(resourceName, "gateway.1260058071.vpn.2866194231.allow_destinations.0.nic", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "gateway.1260058071.vpn.2866194231.weight", "100"),

					resource.TestCheckResourceAttr(resourceName, "log_server.#", "0"),
					// TOOD; get site ID
					// resource.TestCheckResourceAttr(resourceName, "site", "8a4add9e-0e99-4bb1-949c-c9faf9a49ad4"),

					resource.TestCheckResourceAttr(resourceName, "client_interface.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.1281356043.allow_sources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.1281356043.dtls_port", "445"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.1281356043.https_port", "444"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.1281356043.proxy_protocol", "true"),

					resource.TestCheckResourceAttr(resourceName, "networking.#", "1"),

					resource.TestCheckResourceAttr(resourceName, "networking.1374199749.nics.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.1374199749.nics.0.name", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "networking.1374199749.nics.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.1374199749.nics.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.1374199749.nics.0.ipv4.3519857096.dhcp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.1374199749.nics.0.ipv4.3519857096.dhcp.2319808068.dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.1374199749.nics.0.ipv4.3519857096.dhcp.2319808068.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.1374199749.nics.0.ipv4.3519857096.dhcp.2319808068.ntp", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.1374199749.nics.0.ipv4.3519857096.dhcp.2319808068.routers", "true"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				// ImportStateVerify: true,
				ImportStateCheck: testAccApplianceImportStateCheckFunc(1),
			},
		},
	})
}

func testAccCheckApplianceBasicGateway() string {
	return fmt.Sprintf(`
data "appgate_site" "default_site" {
	   site_name = "Default site"
}

resource "appgate_appliance" "test_gateway" {
	name     = "gateway-test"
	hostname = "envy-10-97-168-1338.devops"
    site  = data.appgate_site.default_site.id
	client_interface {
        hostname = "envy-10-97-168-1338.devops"
        proxy_protocol = true
        https_port     = 444
        dtls_port      = 445
        allow_sources {
          address = "1.3.3.7"
          netmask = 0
          nic     = "eth0"
        }
        override_spa_mode = "TCP"
	}

	peer_interface {
		hostname   = "envy-10-97-168-1338.devops"
		https_port = "1337"
	}
	tags = [
		"terraform",
		"api-test-created"
	]
	networking {

		hosts {
		  hostname = "bla"
		  address  = "0.0.0.0"
		}

		nics {
		  enabled = true
		  name    = "eth0"
		  ipv4 {
			dhcp {
			  enabled = true
			  dns     = true
			  routers = true
			  ntp     = true
			}
		  }
		}
	}
    gateway {
        enabled = true
        vpn {
          weight = 100
          allow_destinations {
            address = "0.0.0.0"
            nic     = "eth0"
          }
        }
    }
}
`)
}
