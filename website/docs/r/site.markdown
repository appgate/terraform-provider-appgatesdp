---
layout: "appgate"
page_title: "APPGATE: appgate_site"
sidebar_current: "docs-appgate-resource-site"
description: |-
   Create a new Site.
---

# appgate_site

Create a new Site..

## Example Usage

```hcl

resource "appgate_site" "test_site" {

}

```

## Argument Reference

The following arguments are supported:


* `short_name`: (Optional) A short 4 letter name for the site
* `network_subnets`: (Optional) Network subnets in CIDR format to define the Site&#39;s boundaries. They are added as routes by the Client.
* `ip_pool_mappings`: (Optional) List of IP Pool mappings for this specific Site. When IPs are allocated this Site, they will be mapped to a new one using this setting.
* `default_gateway`: (Optional) Default Gateway configuration.
* `entitlement_based_routing`: (Optional) When enabled, the routes are sent to the Client by the Gateways according to the user&#39;s Entitlements &quot;networkSubnets&quot; should be left be empty if it&#39;s enabled.
* `vpn`: (Optional) VPN configuration for this Site.
* `name_resolution`: (Optional) Settings for asset name resolution.


### default_gateway
Default Gateway configuration.

* `enabled_v4`: When enabled, the Client uses this Site as the Default Default for all IPV4 traffic.
* `enabled_v6`: When enabled, the Client uses this Site as the Default Default for all IPv6 traffic.
* `excluded_subnets`: Network subnets to exclude when Default Gateway is enabled. The traffic for these subnets will not go through the Gateway in this Site.

#### excluded_subnets
Network subnets to exclude when Default Gateway is enabled. The traffic for these subnets will not go through the Gateway in this Site.
### vpn
VPN configuration for this Site.

* `state_sharing`: Configuration for keeping track of states.
* `snat`: Source NAT.
* `tls`: VPN over TLS protocol configuration.
* `dtls`: VPN over DTLS protocol configuration.
* `route_via`: Override routing for tunnel traffic.
* `web_proxy_enabled`: Flag for manipulating web proxy p12 file. Setting this false will delete the existing p12 file from database.
* `web_proxy_key_store`: The PKCS12 package to be used for web proxy. The file must be with no password and must include the full certificate chain and a private key. In Base64 format.
* `web_proxy_certificate_subject_name`: The subject name of the certificate with private key in the PKCS12 file for web proxy assigned to this site.
* `ip_access_log_interval_seconds`: Frequency configuration for generating IP Access audit logs for a connection.

### name_resolution
Settings for asset name resolution.

* `use_hosts_file`: Name resolution to use Appliance&#39;s &#x2F;etc&#x2F;hosts file.
* `dns_resolvers`: Resolver to resolve hostnames using DNS servers.
* `aws_resolvers`: Resolvers to resolve Amazon machines by querying Amazon Web Services.
* `azure_resolvers`: Resolvers to resolve Azure machines by querying Azure App Service.
* `esx_resolvers`: Resolvers to resolve VMware vSphere machines by querying the vCenter.
* `gcp_resolvers`: Resolvers to resolve GCP machine by querying Google web services.

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
* `use_i_a_m_role`: (Optional) Uses the built-in IAM role in AWS instances to authenticate against the API.
* `access_key_id`: (Optional) ID of the access key.
* `secret_access_key`: (Optional) Secret access key for accessKeyId.
* `https_proxy`: (Optional) Proxy address to use while communicating with AWS. format: username:password@ip&#x2F;hostname:port
* `resolve_with_master_credentials`: (Optional) Use master credentials to resolve names in addition to any assumed roles.
* `assumed_roles`: (Optional) Roles to be assumed to perform AWS name resolution.











#### azure_resolvers
Resolvers to resolve Azure machines by querying Azure App Service.
* `name`: (Required) Identifier name. Has no functional effect.
* `update_interval`: (Optional) How often will the resolver poll the server. In seconds.
* `subscription_id`: (Required) Azure subscription id, visible with the azure cli command &#x60;azure account show&#x60;.
* `tenant_id`: (Required) Azure tenant id, visible with the azure cli command &#x60;azure account show&#x60;.
* `client_id`: (Required) Azure client id, also called app id. Visible for a given application using the azure cli command &#x60;azure ad app show&#x60;.
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







## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgate_site d3131f83-10d1-4abc-ac0b-7349538e8300
```
