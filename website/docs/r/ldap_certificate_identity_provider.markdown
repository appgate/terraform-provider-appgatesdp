---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_ldap_certificate_identity_provider"
sidebar_current: "docs-appgate-resource-ldap_certificate_identity_provider"
description: |-
   Create a new LdapCertificate Identity Provider.
---

# appgatesdp_ldap_certificate_identity_provider

Create a new Identity Provider.

~> **NOTE:**  The resource documentation is based on the latest available appgate sdp appliance version, which currently is 5.5.0
Some attributes may not be available if you are running an older version, if you try to use an attribute block that is not permitted in your current version, you will be prompted by an error message.


## Example Usage

```hcl


data "appgatesdp_ip_pool" "ip_four_pool" {
  ip_pool_name = "default pool v4"
}

data "appgatesdp_ip_pool" "ip_v6_pool" {
  ip_pool_name = "default pool v6"
}
data "appgatesdp_mfa_provider" "fido" {
  mfa_provider_name = "Default FIDO2 Provider"
}
resource "appgatesdp_ldap_certificate_identity_provider" "ldap_cert_test_resource" {
  name                       = "%s"
  port                       = 389
  admin_distinguished_name   = "CN=admin,OU=Users,DC=company,DC=com"
  hostnames                  = ["dc.ad.company.com"]
  ssl_enabled                = true
  base_dn                    = "OU=Users,DC=company,DC=com"
  object_class               = "user"
  username_attribute         = "sAMAccountName"
  membership_filter          = "(objectCategory=group)"
  membership_base_dn         = "OU=Groups,DC=company,DC=com"
  default                    = false
  inactivity_timeout_minutes = 28
  ip_pool_v4                 = data.appgatesdp_ip_pool.ip_four_pool.id
  ip_pool_v6                 = data.appgatesdp_ip_pool.ip_v6_pool.id
  admin_password             = "helloworld"
  dns_servers = [
    "172.17.18.19",
    "192.100.111.31"
  ]
  dns_search_domains = [
    "internal.company.com"
  ]
  block_local_dns_requests = true
  on_boarding_two_factor {
    mfa_provider_id       = data.appgatesdp_mfa_provider.fido.id
    device_limit_per_user = 6
    message               = "welcome"
  }
  certificate_user_attribute = "blabla"
  ca_certificates = [
    <<-EOF
-----BEGIN CERTIFICATE-----
...
...
...
-----END CERTIFICATE-----
EOF
  ]
  tags = [
    "terraform",
    "api-created"
  ]
  on_demand_claim_mappings {
    command    = "fileSize"
    claim_name = "antiVirusIsRunning"
    parameters {
      path = "/usr/bin/python3"
    }
    platform = "desktop.windows.all"
  }
}


```


## Argument Reference

The following arguments are supported:


* `hostnames`: (Required) Hostnames/IP addresses to connect.
* `port`: (Required) Port to connect.
* `ssl_enabled`: (Optional) Whether to use LDAPS protocol or not.
* `admin_distinguished_name`: (Required) The Distinguished Name to login to LDAP and query users with.
* `admin_password`: (Optional) The password to login to LDAP and query users with. Required on creation.
* `base_dn`: (Optional) The subset of the LDAP server to search users from. If not set, root of the server is used.
* `object_class`: (Optional) The object class of the users to be authenticated and queried.
* `username_attribute`: (Optional) The name of the attribute to get the exact username from the LDAP server.
* `membership_filter`: (Optional) The filter to use while querying users' nested groups.
* `membership_base_dn`: (Optional) The subset of the LDAP server to search groups from. If not set, "baseDn" is used.
* `password_warning`: (Optional) Password warning configuration for Active Directory. If enabled, the client will display the configured message before the password expiration.
* `ca_certificates`: (Required) CA certificates to verify the Client certificates. In PEM format.
* `certificate_user_attribute`: (Optional) The LDAP attribute to compare the Client certificate's Subject Alternative Name.
* `certificate_attribute`: (Optional) The LDAP attribute to compare the Client certificate binary. Leave it null to skip this comparison.
* `skip_x509_external_checks`: (Optional) By default, Controller contacts the endpoints on the certificate extensions in order to verify revocation status and pull the intermediate CA certificates. Set this flag in order to skip them.
* `certificate_priorities`: (Optional) Client will order the available certificates according to the given priority list.


### hostnames
Hostnames/IP addresses to connect.

### password_warning
Password warning configuration for Active Directory. If enabled, the client will display the configured message before the password expiration.

* `enabled`:  (Optional) Whether to check and warn the users for password expiration.
* `threshold_days`:  (Optional)  default value `5` How many days before the password warning to be displayed to the user.
* `message`:  (Optional) The given message will be displayed to the user. Use this field to guide the users on how to change their passwords. The expiration time will displayed on the client on a separate section. Example: Your password is about to expire. Please change it..
### ca_certificates
CA certificates to verify the Client certificates. In PEM format.

### certificate_priorities
Client will order the available certificates according to the given priority list.

* `type`:  (Optional)  default value `Template`  Enum values: `Template,Issuer`Which attribute to compare.
* `value`: (Required) The value to compare to the certificate attribute.



## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_ldap_certificate_identity_provider d3131f83-10d1-4abc-ac0b-7349538e8300
```
