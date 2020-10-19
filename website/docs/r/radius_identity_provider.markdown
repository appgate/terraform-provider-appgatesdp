---
layout: "appgate"
page_title: "APPGATE: appgate_radius_identity_provider"
sidebar_current: "docs-appgate-resource-radius_identity_provider"
description: |-
   Create a new Radius Identity Provider.
---

# appgate_radius_identity_provider

Create a new Radius Identity Provider.

## Example Usage

```hcl

data "appgate_ip_pool" "ip_sex_pool" {
  ip_pool_name = "default pool v6"
}

data "appgate_ip_pool" "ip_four_pool" {
  ip_pool_name = "default pool v4"
}

data "appgate_mfa_provider" "fido" {
  mfa_provider_name = "Default FIDO2 Provider"
}

resource "appgate_radius_identity_provider" "radius" {
  name = "the-radius"
  hostnames = [
    "radius.company.com"
  ]
  admin_provider = true
  port           = 1812
  shared_secret  = "hunter2"
  ip_pool_v4     = data.appgate_ip_pool.ip_four_pool.id
  ip_pool_v6     = data.appgate_ip_pool.ip_sex_pool.id
  dns_servers = [
    "172.17.18.19",
    "192.100.111.31"
  ]
  dns_search_domains = [
    "internal.company.com"
  ]
  block_local_dns_requests = true
  on_boarding_two_factor {
    mfa_provider_id       = data.appgate_mfa_provider.fido.id
    device_limit_per_user = 6
    message               = "welcome"
  }
  tags = [
    "terraform",
    "api-created"
  ]
  claim_mappings {
    attribute_name = "objectGUID"
    claim_name     = "userId"
    encrypted      = false
    list           = false
  }
  claim_mappings {
    attribute_name = "sAMAccountName"
    claim_name     = "username"
    encrypted      = false
    list           = false
  }
  claim_mappings {
    attribute_name = "givenName"
    claim_name     = "firstName"
    encrypted      = false
    list           = false
  }
  claim_mappings {
    attribute_name = "sn"
    claim_name     = "lastName"
    encrypted      = false
    list           = false
  }
  claim_mappings {
    attribute_name = "mail"
    claim_name     = "emails"
    encrypted      = false
    list           = true
  }
  claim_mappings {
    attribute_name = "memberOf"
    claim_name     = "groups"
    encrypted      = false
    list           = true
  }

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


* `type`: (Required) The type of the Identity Provider.
* `display_name`: (Optional) The name displayed to the user. "name" field is used for Distinguished Name generation. Deprecated as of 5.1 since the Client does not have the option to choose Identity Provider anymore.
* `default`: (Optional) Whether the provider will be chosen by default in the Client UI. If enabled, it will remove the default flag of the current default Identity Provider.
* `client_provider`: (Optional) Whether the provider will be listed in the Client UI or not. Deprecated as of 5.1 since the Client does not have the option to choose Identity Provider anymore.
* `admin_provider`: (Optional) Whether the provider will be listed in the Admin UI or not.
* `on_boarding_two_factor`: (Optional) On-boarding two-factor authentication settings. Leave it empty keep it disabled.
* `on_boarding_type`: (Optional) Client on-boarding type. Deprecated as of 5.0. Use onBoarding2FA object instead.
* `on_boarding_otp_provider`: (Optional) On-boarding MFA Provider ID if "onBoardingType" is Require2FA.  Deprecated as of 5.0. Use onBoarding2FA object instead.
* `on_boarding_otp_message`: (Optional) On-boarding MFA message to be displayed on the Client UI if "onBoardingType" is Require2FA. Deprecated as of 5.0. Use onBoarding2FA object instead.
* `inactivity_timeout_minutes`: (Optional) (Desktop) clients will sign out automatically after the user has been inactive on the device for the configured duration. Set it to 0 to disable.
* `ip_pool_v4`: (Optional) The IPv4 Pool ID the users in this Identity Provider are going to use to allocate IP addresses for the tunnels.
* `ip_pool_v6`: (Optional) The IPv6 Pool ID the users in this Identity Provider are going to use to allocate IP addresses for the tunnels.
* `dns_servers`: (Optional) The dns servers to be assigned to the Clients of the users in this Identity Provider.
* `dns_search_domains`: (Optional) The dns search domains to be assigned to Clients of the users in this Identity Provider.
* `block_local_dns_requests`: (Optional) Whether the Windows Client will block local DNS requests or not.
* `claim_mappings`: (Optional) The mapping of Identity Provider attributes to claims.
* `on_demand_claim_mappings`: (Optional) The mapping of Identity Provider on demand attributes to claims.
* `id`: (Required) ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.
* `hostnames`: (Required) Hostnames/IP addresses to connect.
* `port`: (Optional) Port to connect.
* `ssl_enabled`: (Optional) Whether to use LDAPS protocol or not.
* `admin_distinguished_name`: (Optional) The Distinguished Name to login to LDAP and query users with.
* `admin_password`: (Optional) The password to login to LDAP and query users with. Required on creation.
* `base_dn`: (Optional) The subset of the LDAP server to search users from. If not set, root of the server is used.
* `object_class`: (Optional) The object class of the users to be authenticated and queried.
* `username_attribute`: (Optional) The name of the attribute to get the exact username from the LDAP server.
* `membership_filter`: (Optional) The filter to use while querying users' nested groups.
* `membership_base_dn`: (Optional) The subset of the LDAP server to search groups from. If not set, "baseDn" is used.
* `password_warning`: (Optional) Password warning configuration for Active Directory. If enabled, the client will display the configured message before the password expiration.
* `shared_secret`: (Required) Radius shared secret to authenticate to the server.
* `authentication_protocol`: (Optional) Radius protocol to use while authenticating users.


### on_boarding_two_factor
On-boarding two-factor authentication settings. Leave it empty keep it disabled.

* `mfa_provider_id`: (Required) MFA provider ID to use for the authentication.
* `message`:  (Optional) On-boarding MFA message to be displayed on the Client UI during the second-factor authentication. Example: Please use your multi factor authentication device to on-board..
* `device_limit_per_user`:  (Optional)  default value `100` The device limit per user. The existing on-boarded devices will still be able to sign in even if the limit is exceeded.
### dns_servers
The dns servers to be assigned to the Clients of the users in this Identity Provider.

### dns_search_domains
The dns search domains to be assigned to Clients of the users in this Identity Provider.

### claim_mappings
The mapping of Identity Provider attributes to claims.

* `attribute_name`:  (Optional) The name of the attribute coming from the Identity Provider.
* `claim_name`:  (Optional) The name of the user claim to be used in Appgate SDP.
* `list`:  (Optional)  default value `false` Whether the claim is expected to be a list and have multiple values or not.
* `encrypt`:  (Optional)  default value `false` Whether the claim should be encrypted or not.
### on_demand_claim_mappings
The mapping of Identity Provider on demand attributes to claims.

* `command`:  (Optional)  Enum values: `fileSize,fileExists,fileCreated,fileUpdated,fileVersion,fileSha512,processRunning,processList,serviceRunning,serviceList,regExists,regQuery,runScript`The name of the command.
* `claim_name`:  (Optional) The name of the device claim to be used in Appgate SDP.
* `parameters`:  (Optional) Depending on the command type, extra parameters to pass to the on-demand claim.
* `platform`:  (Optional)  Enum values: `desktop.windows.all,desktop.macos.all,desktop.linux.all,desktop.all,mobile.android.all,mobile.ios.all,mobile.all,all`The platform(s) to run the on-demand claim.
### tags
Array of tags.

### hostnames
Hostnames/IP addresses to connect.

### password_warning
Password warning configuration for Active Directory. If enabled, the client will display the configured message before the password expiration.

* `enabled`:  (Optional) Whether to check and warn the users for password expiration.
* `threshold_days`:  (Optional)  default value `5` How many days before the password warning to be displayed to the user.
* `message`:  (Optional) The given message will be displayed to the user. Use this field to guide the users on how to change their passwords. The expiration time will displayed on the client on a separate section. Example: Your password is about to expire. Please change it..



## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgate_radius_identity_provider d3131f83-10d1-4abc-ac0b-7349538e8300
```
