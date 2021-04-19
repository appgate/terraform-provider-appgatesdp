package appgate

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccApplianceBasicController(t *testing.T) {
	resourceName := "appgate_appliance.test_controller"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplianceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckApplianceBasicController(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplianceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "hostname", "envy-10-97-168-1337.devops"),
					resource.TestCheckResourceAttr(resourceName, "connect_to_peers_using_client_port_with_spa", "true"),

					resource.TestCheckResourceAttr(resourceName, "client_interface.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.allow_sources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.allow_sources.0.address", "1.3.3.7"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.allow_sources.0.netmask", "0"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.allow_sources.0.nic", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.dtls_port", "445"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.hostname", "envy-10-97-168-1337.devops"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.https_port", "444"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.override_spa_mode", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.proxy_protocol", "true"),

					resource.TestCheckResourceAttr(resourceName, "controller.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "controller.0.enabled", "true"),

					resource.TestCheckResourceAttr(resourceName, "healthcheck_server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "healthcheck_server.0.allow_sources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "healthcheck_server.0.allow_sources.0.address", "1.3.3.7"),
					resource.TestCheckResourceAttr(resourceName, "healthcheck_server.0.allow_sources.0.netmask", "0"),
					resource.TestCheckResourceAttr(resourceName, "healthcheck_server.0.allow_sources.0.nic", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "healthcheck_server.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "healthcheck_server.0.port", "5555"),

					resource.TestCheckResourceAttr(resourceName, "log_forwarder.#", "1"),

					resource.TestCheckResourceAttr(resourceName, "log_forwarder.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.0.sites.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.0.tcp_clients.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.0.tcp_clients.0.format", "json"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.0.tcp_clients.0.host", "siem.company.com"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.0.tcp_clients.0.name", "Company SIEM"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.0.tcp_clients.0.port", "8888"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.0.tcp_clients.0.use_tls", "true"),

					resource.TestCheckResourceAttr(resourceName, "networking.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.dns_domains.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.dns_domains.0", "aa.com"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.dns_servers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.dns_servers.0", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.dns_servers.1", "8.8.8.8"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.hosts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.hosts.0.address", "0.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.hosts.0.hostname", "bla"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.0.dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.0.ntp", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.0.routers", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.static.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.static.0.address", "10.10.10.1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.static.0.hostname", "appgate.company.com"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.static.0.netmask", "24"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.static.0.snat", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.static.1.address", "20.20.20.1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.static.1.hostname", "test.company.com"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.static.1.netmask", "32"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.static.1.snat", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.0.dhcp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.0.dhcp.0.dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.0.dhcp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.0.dhcp.0.ntp", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.0.static.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.name", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv4.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv4.0.dhcp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv4.0.dhcp.0.dns", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv4.0.dhcp.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv4.0.dhcp.0.ntp", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv4.0.dhcp.0.routers", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv4.0.static.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv6.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv6.0.dhcp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv6.0.dhcp.0.dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv6.0.dhcp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv6.0.dhcp.0.ntp", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv6.0.static.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.name", "eth1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.routes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.routes.0.address", "0.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.routes.0.gateway", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.routes.0.netmask", "24"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.routes.0.nic", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "ntp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.0.hostname", "0.ubuntu.pool.ntp.org"),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.0.key", ""),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.0.key_type", ""),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.1.hostname", "1.ubuntu.pool.ntp.org"),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.1.key", ""),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.1.key_type", ""),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.2.hostname", "2.ubuntu.pool.ntp.org"),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.2.key", ""),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.2.key_type", ""),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.3.hostname", "3.ubuntu.pool.ntp.org"),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.3.key", ""),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.3.key_type", ""),

					resource.TestCheckResourceAttr(resourceName, "peer_interface.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "peer_interface.0.allow_sources.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "peer_interface.0.hostname", "envy-10-97-168-1337.devops"),
					resource.TestCheckResourceAttr(resourceName, "peer_interface.0.https_port", "1337"),

					resource.TestCheckResourceAttr(resourceName, "ping.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ping.0.allow_sources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ping.0.allow_sources.0.address", "1.3.3.7"),
					resource.TestCheckResourceAttr(resourceName, "ping.0.allow_sources.0.netmask", "0"),
					resource.TestCheckResourceAttr(resourceName, "ping.0.allow_sources.0.nic", "eth0"),

					resource.TestCheckResourceAttr(resourceName, "prometheus_exporter.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "prometheus_exporter.0.allow_sources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "prometheus_exporter.0.allow_sources.0.address", "1.3.3.7"),
					resource.TestCheckResourceAttr(resourceName, "prometheus_exporter.0.allow_sources.0.netmask", "0"),
					resource.TestCheckResourceAttr(resourceName, "prometheus_exporter.0.allow_sources.0.nic", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "prometheus_exporter.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "prometheus_exporter.0.port", "1234"),

					resource.TestCheckResourceAttr(resourceName, "rsyslog_destinations.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "snmp_server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.0.allow_sources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.0.allow_sources.0.address", "1.3.3.7"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.0.allow_sources.0.netmask", "0"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.0.allow_sources.0.nic", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.0.snmpd_conf", "foo"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.0.tcp_port", "161"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.0.udp_port", "161"),

					resource.TestCheckResourceAttr(resourceName, "ssh_server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ssh_server.0.allow_sources.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ssh_server.0.allow_sources.0.address", "1.3.3.7"),
					resource.TestCheckResourceAttr(resourceName, "ssh_server.0.allow_sources.0.netmask", "0"),
					resource.TestCheckResourceAttr(resourceName, "ssh_server.0.allow_sources.0.nic", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "ssh_server.0.allow_sources.1.address", "0.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "ssh_server.0.allow_sources.1.netmask", "0"),
					resource.TestCheckResourceAttr(resourceName, "ssh_server.0.allow_sources.1.nic", "eth1"),
					resource.TestCheckResourceAttr(resourceName, "ssh_server.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssh_server.0.password_authentication", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssh_server.0.port", "2222"),

					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-test-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateCheck:        testAccApplianceImportStateCheckFunc(1),
				ImportStateVerifyIgnore: []string{"site", "seed_file"},
			},
			{
				Config: testAccCheckApplianceBasicControllerUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplianceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName+" updated"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "hostname", "updated-envy-10-97-168-1337.devops"),
					resource.TestCheckResourceAttr(resourceName, "connect_to_peers_using_client_port_with_spa", "true"),

					resource.TestCheckResourceAttr(resourceName, "client_interface.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.allow_sources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.allow_sources.0.address", "1.3.3.6"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.allow_sources.0.netmask", "32"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.allow_sources.0.nic", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.dtls_port", "4454"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.hostname", "updated-envy-10-97-168-1337.devops"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.https_port", "4444"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.override_spa_mode", "UDP-TCP"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.proxy_protocol", "true"),
					resource.TestCheckResourceAttr(resourceName, "peer_interface.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "peer_interface.0.allow_sources.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "peer_interface.0.hostname", "updated-envy-10-97-168-1337.devops"),
					resource.TestCheckResourceAttr(resourceName, "peer_interface.0.https_port", "13371"),

					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-test-created-updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "new"),
					resource.TestCheckResourceAttr(resourceName, "tags.2", "terraform"),

					resource.TestCheckResourceAttr(resourceName, "ntp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.0.hostname", "0.ubuntu.pool.ntp.org"),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.0.key", ""),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.0.key_type", ""),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.1.hostname", "1.ubuntu.pool.ntp.org"),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.1.key", ""),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.1.key_type", ""),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.2.hostname", "23.ubuntu.pool.ntp.org"),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.2.key", ""),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.2.key_type", ""),

					resource.TestCheckResourceAttr(resourceName, "networking.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.hosts.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.hosts.0.address", "1.0.1.0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.hosts.0.hostname", "bla-updated"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.hosts.1.address", "10.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.hosts.1.hostname", "foobar"),

					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.0.dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.0.ntp", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.0.routers", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.static.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.static.0.address", "10.10.10.10"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.static.0.hostname", "appgate-updated.company.com"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.static.0.netmask", "32"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.static.0.snat", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.static.1.address", "20.20.20.1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.static.1.hostname", "test.company.com"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.static.1.netmask", "32"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.static.1.snat", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.0.dhcp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.0.dhcp.0.dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.0.dhcp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.0.dhcp.0.ntp", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.0.static.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.name", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv4.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv4.0.dhcp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv4.0.dhcp.0.dns", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv4.0.dhcp.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv4.0.dhcp.0.ntp", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv4.0.dhcp.0.routers", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv4.0.static.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv6.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv6.0.dhcp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv6.0.dhcp.0.dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv6.0.dhcp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv6.0.dhcp.0.ntp", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv6.0.static.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.name", "eth1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.2.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.2.ipv4.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.2.ipv4.0.dhcp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.2.ipv4.0.dhcp.0.dns", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.2.ipv4.0.dhcp.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.2.ipv4.0.dhcp.0.ntp", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.2.ipv4.0.dhcp.0.routers", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.2.ipv4.0.static.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.2.ipv6.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.2.ipv6.0.dhcp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.2.ipv6.0.dhcp.0.dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.2.ipv6.0.dhcp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.2.ipv6.0.dhcp.0.ntp", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.2.ipv6.0.static.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.2.name", "eth2"),

					resource.TestCheckResourceAttr(resourceName, "networking.0.dns_domains.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.dns_domains.0", "aa-updated.com"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.dns_domains.1", "new.com"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.dns_servers.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.dns_servers.0", "1.1.1.1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.dns_servers.1", "2.2.2.2"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.dns_servers.2", "8.8.4.4"),

					resource.TestCheckResourceAttr(resourceName, "networking.0.routes.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.routes.0.address", "1.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.routes.0.gateway", "1.2.3.5"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.routes.0.netmask", "16"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.routes.0.nic", "eth1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.routes.1.address", "10.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.routes.1.gateway", "10.0.0.1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.routes.1.netmask", "24"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.routes.1.nic", "eth0"),

					resource.TestCheckResourceAttr(resourceName, "controller.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "controller.0.enabled", "true"),

					resource.TestCheckResourceAttr(resourceName, "snmp_server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.0.allow_sources.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.0.allow_sources.0.address", "1.3.3.3"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.0.allow_sources.0.netmask", "24"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.0.allow_sources.0.nic", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.0.allow_sources.1.address", "1.3.3.7"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.0.allow_sources.1.netmask", "24"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.0.allow_sources.1.nic", "eth1"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.0.snmpd_conf", "bar"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.0.tcp_port", "160"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.0.udp_port", "160"),

					resource.TestCheckResourceAttr(resourceName, "ssh_server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ssh_server.0.allow_sources.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "ssh_server.0.allow_sources.0.address", "1.3.3.8"),
					resource.TestCheckResourceAttr(resourceName, "ssh_server.0.allow_sources.0.netmask", "0"),
					resource.TestCheckResourceAttr(resourceName, "ssh_server.0.allow_sources.0.nic", "eth1"),
					resource.TestCheckResourceAttr(resourceName, "ssh_server.0.allow_sources.1.address", "10.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "ssh_server.0.allow_sources.1.netmask", "32"),
					resource.TestCheckResourceAttr(resourceName, "ssh_server.0.allow_sources.1.nic", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "ssh_server.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "ssh_server.0.password_authentication", "false"),
					resource.TestCheckResourceAttr(resourceName, "ssh_server.0.port", "2223"),

					resource.TestCheckResourceAttr(resourceName, "prometheus_exporter.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "prometheus_exporter.0.allow_sources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "prometheus_exporter.0.allow_sources.0.address", "1.3.3.8"),
					resource.TestCheckResourceAttr(resourceName, "prometheus_exporter.0.allow_sources.0.netmask", "0"),
					resource.TestCheckResourceAttr(resourceName, "prometheus_exporter.0.allow_sources.0.nic", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "prometheus_exporter.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "prometheus_exporter.0.port", "1235"),

					resource.TestCheckResourceAttr(resourceName, "healthcheck_server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "healthcheck_server.0.allow_sources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "healthcheck_server.0.allow_sources.0.address", "1.3.3.8"),
					resource.TestCheckResourceAttr(resourceName, "healthcheck_server.0.allow_sources.0.netmask", "0"),
					resource.TestCheckResourceAttr(resourceName, "healthcheck_server.0.allow_sources.0.nic", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "healthcheck_server.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "healthcheck_server.0.port", "5556"),

					resource.TestCheckResourceAttr(resourceName, "ping.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ping.0.allow_sources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ping.0.allow_sources.0.address", "1.3.3.8"),
					resource.TestCheckResourceAttr(resourceName, "ping.0.allow_sources.0.netmask", "0"),
					resource.TestCheckResourceAttr(resourceName, "ping.0.allow_sources.0.nic", "eth0"),

					resource.TestCheckResourceAttr(resourceName, "log_forwarder.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.0.sites.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.0.tcp_clients.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.0.tcp_clients.0.format", "json"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.0.tcp_clients.0.host", "siem-updated.company.com"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.0.tcp_clients.0.name", "Company SIEM updated"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.0.tcp_clients.0.port", "8887"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.0.tcp_clients.0.use_tls", "false"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateCheck:        testAccApplianceImportStateCheckFunc(1),
				ImportStateVerifyIgnore: []string{"site", "seed_file"},
			},
			{
				Config: testAccCheckApplianceBasicControllerDisableDelete(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplianceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName+" disable"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),
					resource.TestCheckResourceAttr(resourceName, "hostname", "updated-envy-10-97-168-1337.devops"),
					resource.TestCheckResourceAttr(resourceName, "connect_to_peers_using_client_port_with_spa", "true"),

					resource.TestCheckResourceAttr(resourceName, "client_interface.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.allow_sources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.allow_sources.0.address", "1.3.3.6"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.allow_sources.0.netmask", "32"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.allow_sources.0.nic", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.dtls_port", "4454"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.hostname", "updated-envy-10-97-168-1337.devops"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.https_port", "4444"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.override_spa_mode", "UDP-TCP"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.proxy_protocol", "true"),
					resource.TestCheckResourceAttr(resourceName, "peer_interface.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "peer_interface.0.allow_sources.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "peer_interface.0.hostname", "updated-envy-10-97-168-1337.devops"),
					resource.TestCheckResourceAttr(resourceName, "peer_interface.0.https_port", "13371"),

					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-test-created-updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),

					resource.TestCheckResourceAttr(resourceName, "ntp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.0.hostname", "0.ubuntu.pool.ntp.org"),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.0.key", ""),
					resource.TestCheckResourceAttr(resourceName, "ntp.0.servers.0.key_type", ""),

					resource.TestCheckResourceAttr(resourceName, "networking.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.hosts.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.0.dns", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.0.ntp", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.0.routers", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.static.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.0.dhcp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.0.dhcp.0.dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.0.dhcp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.0.dhcp.0.ntp", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.0.static.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.name", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv4.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv4.0.dhcp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv4.0.dhcp.0.dns", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv4.0.dhcp.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv4.0.dhcp.0.ntp", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv4.0.dhcp.0.routers", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv4.0.static.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv6.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv6.0.dhcp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv6.0.dhcp.0.dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv6.0.dhcp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv6.0.dhcp.0.ntp", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.ipv6.0.static.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.1.name", "eth1"),

					resource.TestCheckResourceAttr(resourceName, "networking.0.dns_domains.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.dns_servers.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.routes.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "controller.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "controller.0.enabled", "false"),

					resource.TestCheckResourceAttr(resourceName, "snmp_server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "snmp_server.0.enabled", "false"),

					resource.TestCheckResourceAttr(resourceName, "ssh_server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ssh_server.0.enabled", "false"),

					resource.TestCheckResourceAttr(resourceName, "prometheus_exporter.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "prometheus_exporter.0.enabled", "false"),

					resource.TestCheckResourceAttr(resourceName, "healthcheck_server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "healthcheck_server.0.enabled", "false"),

					resource.TestCheckResourceAttr(resourceName, "ping.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "ping.0.allow_sources.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "log_forwarder.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "log_forwarder.0.enabled", "false"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateCheck:        testAccApplianceImportStateCheckFunc(1),
				ImportStateVerifyIgnore: []string{"site", "seed_file"},
			},
		},
	})
}

//lint:file-ignore U1000 Debug function, used during development
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

		fmt.Println(rs.Primary)
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

func testAccCheckApplianceBasicController(rName string) string {
	return fmt.Sprintf(`
data "appgate_site" "default_site" {
  site_name = "Default Site"
}

resource "appgate_appliance" "test_controller" {
  name     = "%s"
  hostname = "envy-10-97-168-1337.devops"
  connect_to_peers_using_client_port_with_spa = true
  client_interface {
    hostname       = "envy-10-97-168-1337.devops"
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
          enabled = false
          dns     = true
          routers = true
          ntp     = true
        }

        static {
          address  = "10.10.10.1"
          netmask  = 24
          hostname = "appgate.company.com"
          snat     = true
        }

        static {
          address  = "20.20.20.1"
          netmask  = 32
          hostname = "test.company.com"
          snat     = false
        }
      }
      ipv6 {
        dhcp {
          enabled = false
          dns     = true
          ntp     = false
        }
      }

    }

    nics {
      enabled = true
      name    = "eth1"
      ipv4 {
        dhcp {
          enabled = true
          dns     = false
          routers = false
          ntp     = false
        }
      }
      ipv6 {
        dhcp {
          enabled = false
          dns     = true
          ntp     = false
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
      nic     = "eth0"
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
  ssh_server {
    enabled                 = true
    port                    = 2222
    password_authentication = true
    allow_sources {
      address = "1.3.3.7"
      netmask = 0
      nic     = "eth0"
    }
    allow_sources {
      address = "0.0.0.0"
      netmask = 0
      nic     = "eth1"
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
  healthcheck_server {
    enabled = true
    port    = 5555
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
     tcp_clients {
      name    = "Company SIEM"
      host    = "siem.company.com"
      port    = 8888
      format  = "json"
      use_tls = true
      filter  = "log.client_ip=='10.0.23.523'"
    }
    sites = [
      data.appgate_site.default_site.id
    ]
  }
}
`, rName)
}

func testAccCheckApplianceBasicControllerUpdate(rName string) string {
	return fmt.Sprintf(`
data "appgate_site" "default_site" {
  site_name = "Default Site"
}

resource "appgate_appliance" "test_controller" {
  name     = "%s updated"
  hostname = "updated-envy-10-97-168-1337.devops"

  client_interface {
    hostname       = "updated-envy-10-97-168-1337.devops"
    proxy_protocol = true
    https_port     = 4444
    dtls_port      = 4454
    allow_sources {
      address = "1.3.3.6"
      netmask = 32
      nic     = "eth0"
    }
    override_spa_mode = "UDP-TCP"
  }

  peer_interface {
    hostname   = "updated-envy-10-97-168-1337.devops"
    https_port = "13371"
  }
  tags = [
    "terraform",
    "api-test-created-updated",
    "new"
  ]
  ntp {
    servers {
      hostname = "0.ubuntu.pool.ntp.org"
    }
    servers {
      hostname = "1.ubuntu.pool.ntp.org"
    }
    servers {
      hostname = "23.ubuntu.pool.ntp.org"
    }
  }
  networking {

    hosts {
      hostname = "bla-updated"
      address  = "1.0.1.0"
    }
    hosts {
      hostname = "foobar"
      address  = "10.0.0.0"
    }

    nics {
      enabled = true
      name    = "eth0"

      ipv4 {
        dhcp {
          enabled = false
          dns     = true
          routers = true
          ntp     = true
        }

        static {
          address  = "10.10.10.10"
          netmask  = 32
          hostname = "appgate-updated.company.com"
          snat     = false
        }

        static {
          address  = "20.20.20.1"
          netmask  = 32
          hostname = "test.company.com"
          snat     = false
        }
      }
      ipv6 {
        dhcp {
          enabled = false
          dns     = true
          ntp     = false
        }
      }

    }

    nics {
      enabled = true
      name    = "eth1"
      ipv4 {
        dhcp {
          enabled = true
          dns     = false
          routers = false
          ntp     = false
        }
      }
      ipv6 {
        dhcp {
          enabled = false
          dns     = true
          ntp     = false
        }
      }
    }

    nics {
      enabled = true
      name    = "eth2"
      ipv4 {
        dhcp {
          enabled = true
          dns     = false
          routers = false
          ntp     = false
        }
      }
      ipv6 {
        dhcp {
          enabled = false
          dns     = true
          ntp     = false
        }
      }
    }

    dns_servers = [
      "8.8.4.4",
      "1.1.1.1",
      "2.2.2.2"
    ]
    dns_domains = [
      "aa-updated.com",
      "new.com"
    ]

    routes {
      address = "1.0.0.0"
      netmask = 16
      gateway = "1.2.3.5"
      nic     = "eth1"
    }
    routes {
      address = "10.0.0.0"
      netmask = 24
      gateway = "10.0.0.1"
      nic     = "eth0"
    }
  }
  controller {
    enabled = true
  }

  snmp_server {
    enabled    = true
    tcp_port   = 160
    udp_port   = 160
    snmpd_conf = "bar"
    allow_sources {
      address = "1.3.3.3"
      netmask = 24
      nic     = "eth0"
    }
    allow_sources {
      address = "1.3.3.7"
      netmask = 24
      nic     = "eth1"
    }
  }

  ssh_server {
    enabled                 = true
    port                    = 2223
    password_authentication = false
    allow_sources {
      address = "1.3.3.8"
      netmask = 0
      nic     = "eth1"
    }
    allow_sources {
      address = "10.0.0.0"
      netmask = 32
      nic     = "eth0"
    }
  }
  prometheus_exporter {
    enabled = true
    port    = 1235
    allow_sources {
      address = "1.3.3.8"
      netmask = 0
      nic     = "eth0"
    }
  }
  healthcheck_server {
    enabled = true
    port    = 5556
    allow_sources {
      address = "1.3.3.8"
      netmask = 0
      nic     = "eth0"
    }
  }

  ping {
    allow_sources {
      address = "1.3.3.8"
      netmask = 0
      nic     = "eth0"
    }
  }

  log_forwarder {
    enabled = true
     tcp_clients {
      name    = "Company SIEM updated"
      host    = "siem-updated.company.com"
      port    = 8887
      format  = "json"
      use_tls = false
      filter  = "log.client_ip=='10.0.24.523'"
    }
    sites = [
      data.appgate_site.default_site.id
    ]
  }
}
`, rName)
}

func testAccCheckApplianceBasicControllerDisableDelete(rName string) string {
	return fmt.Sprintf(`
data "appgate_site" "default_site" {
  site_name = "Default site"
}

resource "appgate_appliance" "test_controller" {
  name     = "%s disable"
  hostname = "updated-envy-10-97-168-1337.devops"

  client_interface {
    hostname       = "updated-envy-10-97-168-1337.devops"
    proxy_protocol = true
    https_port     = 4444
    dtls_port      = 4454
    allow_sources {
      address = "1.3.3.6"
      netmask = 32
      nic     = "eth0"
    }
    override_spa_mode = "UDP-TCP"
  }

  peer_interface {
    hostname   = "updated-envy-10-97-168-1337.devops"
    https_port = "13371"
  }

  tags = [
    "terraform",
    "api-test-created-updated",
  ]

  ntp {
    servers {
      hostname = "0.ubuntu.pool.ntp.org"
    }
  }

  networking {

    nics {
      enabled = true
      name    = "eth0"
      ipv4 {
        dhcp {
          enabled = true
          dns     = false
          routers = false
          ntp     = false
        }
      }
      ipv6 {
        dhcp {
          enabled = false
          dns     = true
          ntp     = false
        }
      }
    }

    nics {
      enabled = true
      name    = "eth1"
      ipv4 {
        dhcp {
          enabled = true
          dns     = false
          routers = false
          ntp     = false
        }
      }
      ipv6 {
        dhcp {
          enabled = false
          dns     = true
          ntp     = false
        }
      }
    }

  }

  controller {
    enabled = false
  }

  snmp_server {
    enabled    = false
    tcp_port   = 160
    udp_port   = 160
    snmpd_conf = "bar"
  }

  ssh_server {
    enabled                 = false
    port                    = 2223
    password_authentication = false
  }

  prometheus_exporter {
    enabled = false
    port    = 1235
  }

  healthcheck_server {
    enabled = false
    port    = 5556
  }

  ping {
  }

  log_forwarder {
    enabled = false
  }
}
`, rName)
}

func TestAccApplianceConnector(t *testing.T) {
	resourceName := "appgate_appliance.connector"
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplianceDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckApplianceBasicConnector(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplianceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "client_interface.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.allow_sources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.allow_sources.0.address", "1.3.3.7"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.allow_sources.0.netmask", "0"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.allow_sources.0.nic", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.dtls_port", "445"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.hostname", "envy-10-97-168-1234.devops"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.https_port", "444"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.override_spa_mode", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.proxy_protocol", "true"),

					resource.TestCheckResourceAttr(resourceName, "hostname", "envy-10-97-168-1234.devops"),

					resource.TestCheckResourceAttr(resourceName, "connector.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "connector.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "connector.0.express_clients.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "connector.0.express_clients.0.allow_resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "connector.0.express_clients.0.allow_resources.0.address", "0.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "connector.0.express_clients.0.allow_resources.0.netmask", "32"),
					resource.TestCheckResourceAttr(resourceName, "connector.0.express_clients.0.device_id", "12699e27-b584-464a-81ee-5b4784b6d425"),
					resource.TestCheckResourceAttr(resourceName, "connector.0.express_clients.0.name", "Printers"),
					resource.TestCheckResourceAttr(resourceName, "connector.0.express_clients.0.snat_to_resources", "true"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),

					resource.TestCheckResourceAttr(resourceName, "networking.#", "1"),

					resource.TestCheckResourceAttr(resourceName, "networking.0.hosts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.hosts.0.address", "0.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.hosts.0.hostname", "bla"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.0.dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.0.ntp", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.0.routers", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.static.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.0.dhcp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.0.dhcp.0.dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.0.dhcp.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.0.dhcp.0.ntp", "false"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv6.0.static.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.name", "eth0"),

					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),

					resource.TestCheckResourceAttr(resourceName, "peer_interface.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "peer_interface.0.allow_sources.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "peer_interface.0.hostname", "envy-10-97-168-1234.devops"),
					resource.TestCheckResourceAttr(resourceName, "peer_interface.0.https_port", "1337"),

					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", "api-test-created"),
					resource.TestCheckResourceAttr(resourceName, "tags.1", "terraform"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck:  testAccApplianceImportStateCheckFunc(1),
				ImportStateVerifyIgnore: []string{
					"site",
					"seed_file",
				},
			},
		},
	})
}

func testAccCheckApplianceBasicConnector(rName string) string {
	return fmt.Sprintf(`
data "appgate_site" "default_site" {
    site_name = "Default site"
}
resource "appgate_appliance" "connector" {
    name     = "%s"
    hostname = "envy-10-97-168-1234.devops"
    site     = data.appgate_site.default_site.id

    client_interface {
      hostname       = "envy-10-97-168-1234.devops"
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
        ipv6 {
          dhcp {
            enabled = false
            dns     = true
            ntp     = false
          }
        }
      }
    }
    connector {
      enabled = true
      express_clients {
        name      = "Printers"
        device_id = "12699e27-b584-464a-81ee-5b4784b6d425"

        allow_resources {
          address = "0.0.0.0"
          netmask = 32
        }
        snat_to_resources = true
      }
    }
}
`, rName)
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
	rName := RandStringFromCharSet(10, CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckApplianceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApplianceBasicGateway(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckApplianceExists(resourceName),

					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "hostname", "envy-10-97-168-1338.devops"),
					resource.TestCheckResourceAttr(resourceName, "notes", "Managed by terraform"),

					resource.TestCheckResourceAttr(resourceName, "controller.#", "0"),

					resource.TestCheckResourceAttr(resourceName, "gateway.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "gateway.0.vpn.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "gateway.0.vpn.0.allow_destinations.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "gateway.0.vpn.0.allow_destinations.0.address", "0.0.0.0"),
					resource.TestCheckResourceAttr(resourceName, "gateway.0.vpn.0.allow_destinations.0.netmask", "0"),
					resource.TestCheckResourceAttr(resourceName, "gateway.0.vpn.0.allow_destinations.0.nic", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "gateway.0.vpn.0.weight", "100"),

					resource.TestCheckResourceAttr(resourceName, "log_server.#", "0"),
					// TODO; get site ID
					// resource.TestCheckResourceAttr(resourceName, "site", "8a4add9e-0e99-4bb1-949c-c9faf9a49ad4"),

					resource.TestCheckResourceAttr(resourceName, "client_interface.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.allow_sources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.dtls_port", "445"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.https_port", "444"),
					resource.TestCheckResourceAttr(resourceName, "client_interface.0.proxy_protocol", "true"),

					resource.TestCheckResourceAttr(resourceName, "networking.#", "1"),

					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.name", "eth0"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.0.dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.0.enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.0.ntp", "true"),
					resource.TestCheckResourceAttr(resourceName, "networking.0.nics.0.ipv4.0.dhcp.0.routers", "true"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateCheck:        testAccApplianceImportStateCheckFunc(1),
				ImportStateVerifyIgnore: []string{"site", "seed_file"},
			},
		},
	})
}

func testAccCheckApplianceBasicGateway(rName string) string {
	return fmt.Sprintf(`
data "appgate_site" "default_site" {
  site_name = "Default site"
}

resource "appgate_appliance" "test_gateway" {
  name     = "%s"
  hostname = "envy-10-97-168-1338.devops"
  site     = data.appgate_site.default_site.id
  client_interface {
    hostname       = "envy-10-97-168-1338.devops"
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
      ipv6 {
        dhcp {
          enabled = false
          dns     = true
          ntp     = false
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
`, rName)
}
