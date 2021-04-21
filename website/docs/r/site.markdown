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

}

```

## Argument Reference

The following arguments are supported:


* `short_name`: (Optional) A short 4 letter name for the Site to be displayed on the Client.
* `description`: (Optional) Description of the Site to be displayed on the Client.
* `network_subnets`: (Optional) Network subnets in CIDR format to define the Site's boundaries. They are added as routes by the Client.
* `ip_pool_mappings`: (Optional) List of IP Pool mappings for this specific Site. When IPs are allocated this Site, they will be mapped to a new one using this setting.
* `default_gateway`: (Optional) Default Gateway configuration.
* `entitlement_based_routing`: (Optional) When enabled, the routes are sent to the Client by the Gateways according to the user's Entitlements "networkSubnets" should be left be empty if it's enabled.
* `vpn`: (Optional) VPN configuration for this Site.
* `name_resolution`: (Optional) Settings for asset name resolution.
* `id`: (Required) ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.


### network_subnets
Network subnets in CIDR format to define the Site's boundaries. They are added as routes by the Client.

### ip_pool_mappings
List of IP Pool mappings for this specific Site. When IPs are allocated this Site, they will be mapped to a new one using this setting.

* `from`:  (Optional) IP Pool ID to map from. If a user is authorizing with this IP Pool via Identity Provider assignment and has access to this Site, mapping will occur for that user.
* `to`:  (Optional) IP Pool ID to map to.
### default_gateway
Default Gateway configuration.

* `enabled_v4`:  (Optional)  default value `false` When enabled, the Client uses this Site as the Default Default for all IPV4 traffic.
* `enabled_v6`:  (Optional)  default value `false` When enabled, the Client uses this Site as the Default Default for all IPv6 traffic.
* `excluded_subnets`:  (Optional) Network subnets to exclude when Default Gateway is enabled. The traffic for these subnets will not go through the Gateway in this Site.
#### excluded_subnets
Network subnets to exclude when Default Gateway is enabled. The traffic for these subnets will not go through the Gateway in this Site.
### vpn
VPN configuration for this Site.

* `state_sharing`:  (Optional)  default value `false` Configuration for keeping track of states.
* `snat`:  (Optional)  default value `false` Source NAT.
* `tls`:  (Optional) VPN over TLS protocol configuration.
* `dtls`:  (Optional) VPN over DTLS protocol configuration.
* `route_via`:  (Optional) Override routing for tunnel traffic.
* `web_proxy_enabled`:  (Optional) Flag for manipulating web proxy p12 file. Setting this false will delete the existing p12 file from database.
* `web_proxy_key_store`:  (Optional) The PKCS12 package to be used for web proxy. The file must be with no password and must include the full certificate chain and a private key. In Base64 format.
* `web_proxy_verify_upstream_certificate`:  (Optional)  default value `true` Gateway will verify the certificate of the endpoints.
* `web_proxy_certificate_subject_name`:  (Optional) The subject name of the certificate with private key in the PKCS12 file for web proxy assigned to this site.
* `ip_access_log_interval_seconds`:  (Optional)  default value `120` Frequency configuration for generating IP Access audit logs for a connection.
#### tls
VPN over TLS protocol configuration.
* `enabled`: (Optional)
#### dtls
VPN over DTLS protocol configuration.
* `enabled`: (Optional)
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
#### dns_resolvers
Resolver to resolve hostnames using DNS servers.
* `name`: (Required) Identifier name. Has no functional effect. Example: DNS Resolver 1.
* `update_interval`: (Optional) How often will the resolver poll the server. In seconds.
* `servers`: (Required) DNS Server addresses that will be used to resolve hostnames within the Site.
* `search_domains`: (Optional) DNS search domains that will be used to resolve hostnames within the Site.
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
* `subscription_id`: (Required) Azure subscription id, visible with the azure cli command `azure account show`.
* `tenant_id`: (Required) Azure tenant id, visible with the azure cli command `azure account show`.
* `client_id`: (Required) Azure client id, also called app id. Visible for a given application using the azure cli command `azure ad app show`.
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
### tags
Array of tags.




## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_site d3131f83-10d1-4abc-ac0b-7349538e8300
```
