---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_site"
sidebar_current: "docs-appgate-resource-site"
description: |-
   Create a new Site.
---

# appgatesdp_site
Create a new Site.

## Example Usage
```hcl
resource "appgatesdp_site" "gbg_site" {
    name       = "Gothenburg site"
    short_name = "gbg"
    tags = [
        "developer",
        "api-created"
    ]
    notes = "This object has been created for test purposes."
    network_subnets = [
        "10.0.0.0/16"
    ]
    default_gateway {
        enabled_v4       = false
        enabled_v6       = false
        excluded_subnets = []
    }
    ip_pool_mappings {
        from = "64cbd0ca-688d-4b55-a8a4-4069d8be6ce5"
        to   = "699d9d61-67bb-4cc3-9365-7666d2ddd0a8"
        type = "Allocation
    }
    default_gateway {
        enabled_v4       = false
        enabled_v6       = false
        excluded_subnets = []
    }
    entitlement_based_routing = true
    vpn {
        state_sharing                  = false
        snat                           = false
        ip_access_log_interval_seconds = 120
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
	name_resolution {
	    use_hosts_file = false
	    dns_resolvers {
	        name                = "DNS Resolver 1"
	        update_interval     = ""
	        query_aaaa          = true
	        default_ttl_seconds = 5
	        servers = [
	            "server1"
	        ]
	        search_domains = [
	            "example.com"
	        ]
	    }
	    aws_resolvers {
	        name               = "AWS Resolver 1"
	        update_interval    = 5
	        vpc_auto_discovery = true
	        regions            = "us-west-1"
	        use_iam_role       = true
	        access_key_id      = "access_key_id1"
	        secret_access_key  = "secret_access_key1"
	        https_proxy        = "proxy.example.com"
	        resolve_with_master_credentiasl = true
	        vpcs = [
	            "vpc1"
	        ]
	        assumed_roles {
	            account_id  = "account_id1
	            role_name   = "role_1"
	            external_id = "external_id1"
	            regions = [
	                "us-west-1
                ]
	        }
	    }
	    azure_resolvers {
	        name                   = "Azure Resolver 1"
	        update_interval        = 5
	        use_managed_identities = true
	        subscription_id        = "subscription1"
	        tenant_id              = "tenant1"
	        client_id              = "client_id1"
	        secret                 = "secret1"
	    }
	    esx_resolvers {
	        name            = "ESX Resolver 1"
	        update_interval = 5
	        hostname        = "esx.example.com"
	        username        = "admin"
	        password        = "password"
	    }
	    gcp_resolvers {
	        name            = "GCP Resolver 1"
	        update_interval = 5
	        project_filter  = "project_filter1"
	        instance_filter = "instance_filter1"
	    }
	    dns_forwarding {
            site_ipv4 = ""
            site_ipv6 = ""
            dns_servers = [
                "dns_server1"
            ]
            allow_destinations {
                address = "192.168.1.1"
                netmask = 32
            }
        }
        illumio_resolvers {
            name            = "Illumio Resolver 1"
            hostname        = "illumio.example.com"
            update_interval = 5
            port            = 65530
            username        = "admin"
            password        = "adminadmin"
        }
    }
}
```

## Argument Reference
The following arguments are supported:
* `site_id`: (Optional) Computed if empty -  ID of the object.
* `name`: (Required) Name of the object.
* `description`: (Optional) Description of the Site to be displayed on the Client.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.
* `short_name`: (Optional) A short 4 letter name for the Site to be displayed on the Client.
* `network_subnets`: (Optional) Network subnets in CIDR format to define the Site's boundaries. They are added as routes by the Client.
* `ip_pool_mappings`: (Optional) List of IP Pool mappings for this specific Site. When IPs are allocated this Site, they will be mapped to a new one using this setting.
* `default_gateway`: (Optional) Default Gateway configuration.
* `entitlement_based_routing`: (Optional) When enabled, the routes are sent to the Client by the Gateways according to the user's Entitlements "networkSubnets" should be left be empty if it's enabled.
* `vpn`: (Optional) VPN configuration for this Site.
* `name_resolution`: (Optional) Settings for asset name resolution.

### ip_pool_mappings
List of IP Pool mappings for this specific Site. When IPs are allocated this Site, they will be mapped to a new one using this setting.
* `from`:  (Required) IP Pool ID to map from. If a user is authorizing with this IP Pool via Identity Provider assignment and has access to this Site, mapping will occur for that user.
* `to`:  (Required) IP Pool ID to map to.
* `type`:  (Optional) Mapping type.  Enum values: `Translation,Allocation`

### default_gateway
Default Gateway configuration.
* `enabled_v4`:  (Required - See Note below) default value `false` When enabled, the Client uses this Site as the Default for all IPV4 traffic. At least one of enabled_v4 or enabled_v6b 
* `enabled_v6`:  (Required - See Note below) default value `false` When enabled, the Client uses this Site as the Default for all IPv6 traffic.
* `excluded_subnets`:  (Optional) Network subnets to exclude when Default Gateway is enabled. The traffic for these subnets will not go through the Gateway in this Site. Deprecated as of 6.0. Use action type 'exclude' in Entitlements instead.
> Note: At least one of `enabled_v4` or `enabled_v6` must be set

### vpn
VPN configuration for this Site.
* `state_sharing`:  (Required)  default value `false` Configuration for keeping track of states.
* `snat`:  (Required)  default value `false` Source NAT.
* `tls`:  (Optional) VPN over TLS protocol configuration.
* `dtls`:  (Optional) VPN over DTLS protocol configuration.
* `route_via`:  (Optional) Override routing for tunnel traffic.
* `ip_access_log_interval_seconds`:  (Optional)  default value `120` Frequency configuration for generating IP Access audit logs for a connection.

#### tls
VPN over TLS protocol configuration.
* `enabled`: (Required) Whether to enable tls

#### dtls
VPN over DTLS protocol configuration.
* `enabled`: (Required) Whether to enable dtls

#### route_via
Override routing for tunnel traffic.
* `ipv4`: (Optional) IPv4 address for routing tunnel traffic. Example: 10.0.0.2.
* `ipv6`: (Optional) IPv6 address for routing tunnel traffic. Example: 2001:db8:0:0:0:ff00:42:8329.

### name_resolution
Settings for asset name resolution.
* `use_hosts_file`:  (Optional)  default value `false` Name resolution to use Appliance's /etc/hosts file.
* `dns_resolvers`:  (Optional) Resolver to resolve hostnames using DNS servers.
* `aws_resolvers`:  (Optional) Resolvers to resolve Amazon machines by querying Amazon Web Services.
* `azure_resolvers`:  (Optional) Resolvers to resolve Azure machines by querying Azure App Service.
* `esx_resolvers`:  (Optional) Resolvers to resolve VMware vSphere machines by querying the vCenter.
* `gcp_resolvers`:  (Optional) Resolvers to resolve GCP machine by querying Google web services.
* `dns_forwarding`:  (Optional) Enable DNS Forwarding feature.

#### dns_resolvers
Resolver to resolve hostnames using DNS servers.
* `name`: (Required) Identifier name. Has no functional effect. Example: DNS Resolver 1.
* `update_interval`: (Optional) How often will the resolver poll the server. In seconds.
* `query_aaaa`: (Optional) Perform AAAA lookups.
* `default_ttl_seconds`: (Optional) This will apply whenever Gateway gets a DNS response which has no TTL set.
* `servers`: (Required) DNS Server addresses that will be used to resolve hostnames within the Site.
* `search_domains`: (Optional) DNS search domains that will be used to resolve hostnames within the Site.
* `match_domains`: (Optional) The DNS resolver will only attempt to resolve names matching the match domains. If match domains are not specified the DNS resolver will attempt to resolve all hostnames.

#### aws_resolvers
Resolvers to resolve Amazon machines by querying Amazon Web Services.
* `name`: (Required) Identifier name. Has no functional effect. Example: AWS Resolver 1.
* `update_interval`: (Optional) How often will the resolver poll the server. In seconds.
* `vpcs`: (Optional) VPC IDs to resolve names.
* `vpc_auto_discovery`: (Optional) Use VPC auto discovery.
* `regions`: (Optional) Amazon regions.
* `use_iamrole`: (Optional) Uses the built-in IAM role in AWS instances to authenticate against the API.
* `access_key_id`: (Optional) ID of the access key.
* `secret_access_key`: (Optional) Secret access key for accessKeyId.
* `https_proxy`: (Optional) Proxy address to use while communicating with AWS. format: username:password@ip/hostname:port
* `resolve_with_master_credentials`: (Optional) Use master credentials to resolve names in addition to any assumed roles.
* `assumed_roles`: (Optional) Roles to be assumed to perform AWS name resolution.

#### azure_resolvers
Resolvers to resolve Azure machines by querying Azure App Service.
* `name`: (Required) Identifier name. Has no functional effect.
* `update_interval`: (Optional) How often will the resolver poll the server. In seconds.
* `use_managed_identities`: (Optional) Uses the built-in Managed Identities in Azure instances to authenticate against the API.
* `subscription_id`: (Optional) Azure subscription id, visible with the azure cli command `azure account show`.
* `tenant_id`: (Optional) Azure tenant id, visible with the azure cli command `azure account show`.
* `client_id`: (Optional) Azure client id, also called app id. Visible for a given application using the azure cli command `azure ad app show`.
* `secret`: (Optional) Azure client secret. For Azure AD Apps this is done by creating a key for the app.

#### esx_resolvers
Resolvers to resolve VMware vSphere machines by querying the vCenter.
* `name`: (Required) Identifier name. Has no functional effect.
* `update_interval`: (Optional) How often will the resolver poll the server. In seconds.
* `hostname`: (Required) Hostname of the vCenter.
* `username`: (Required) Username with admin access to the vCenter.
* `password`: (Optional) Password for the username.

#### gcp_resolvers
Resolvers to resolve GCP machine by querying Google web services.
* `name`: (Required) Identifier name. Has no functional effect.
* `update_interval`: (Optional) How often will the resolver poll the server. In seconds.
* `project_filter`: (Optional) GCP project filter.
* `instance_filter`: (Optional) GCP instance filter.

#### dns_forwarding
Enable DNS Forwarding feature.
* `site_ipv4`: (Optional) DNS Forwarder Site IPv4 address. Example: 100.110.0.0.
* `site_ipv6`: (Optional) DNS Forwarder Site IPv6 address. Example: 2001:db8:0:0:0:ff00:42:8329.
* `dns_servers`: (Required) DNS Servers to use for resolving endpoints. Example: 172.17.18.19,192.100.111.31.
* `allow_destinations`: (Required) A list of subnets to allow access.
* `default_ttl_seconds`: (Optional) This will apply whenever Gateway gets a DNS response which has no TTL set.

#### illumio_resolvers
Resolvers to resolve names by querying Appgate Illumio Resolver.
* `name`: (Required) Identifier name. Has no functional effect.
* `update_interval`: (Optional) How often will the resolver poll the server. In seconds.
* `hostname`: (Required) Hostname of the Illumio Resolver.
* `port`: (Required) Port number of the Illumio Resolver.
* `username`: (Required) Username with access to the Illumio Resolver.
* `password`: (Optional) Password for the username.
* `org_id`: (Optional) Organization ID of the Illumio Resolver

## Import
Instances can be imported using the `id`, e.g.
```
$ terraform import appgatesdp_site.example d3131f83-10d1-4abc-ac0b-7349538e8300
```
