---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_oidc_identity_provider"
sidebar_current: "docs-appgate-resource-oidc_identity_provider"
description: |-
   Create a new oidc Identity Provider.
---

# appgatesdp_oidc_identity_provider

Create a new Identity Provider.

~> **NOTE:**  The resource documentation is based on the latest available appgate sdp appliance version, which currently is 6.0.0
Some attributes may not be available if you are running an older version, if you try to use an attribute block that is not permitted in your current version, you will be prompted by an error message.


## Example Usage

```hcl


data "appgatesdp_ip_pool" "ip_v6_pool" {
  ip_pool_name = "default pool v6"
}

data "appgatesdp_ip_pool" "ip_v4_pool" {
  ip_pool_name = "default pool v4"
}

data "appgatesdp_mfa_provider" "fido" {
  mfa_provider_name = "Default FIDO2 Provider"
}

resource "appgatesdp_oidc_identity_provider" "oidc" {
  issuer = "https://example.com/oidc/issuer"
  audience = "oidc_test_audience"
  scope = "oidc_test_scope"
  google { 
    enabled = true
    client_secret = "oidc_test_client_secret"
    refresh_token = true
  }

  name = "the-oidc"
  admin_provider = true
  ip_pool_v4     = data.appgatesdp_ip_pool.ip_v4_pool.id
  ip_pool_v6     = data.appgatesdp_ip_pool.ip_v6_pool.id
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
  tags = [
    "terraform",
    "api-created"
  ]
  claim_mappings {
    attribute_name = "objectGUID"
    claim_name     = "userId"
    encrypt      = false
    list           = false
  }
  claim_mappings {
    attribute_name = "sAMAccountName"
    claim_name     = "username"
    encrypt      = false
    list           = false
  }
  claim_mappings {
    attribute_name = "givenName"
    claim_name     = "firstName"
    encrypt      = false
    list           = false
  }
  claim_mappings {
    attribute_name = "sn"
    claim_name     = "lastName"
    encrypt      = false
    list           = false
  }
  claim_mappings {
    attribute_name = "mail"
    claim_name     = "emails"
    encrypt      = false
    list           = true
  }
  claim_mappings {
    attribute_name = "memberOf"
    claim_name     = "groups"
    encrypt      = false
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
### Base Identity Provider Arguments
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.
* `type`: (Computed) The type of the Identity Provider.
* `ip_pool_v4`: (Optional) The IPv4 Pool ID the users in this Identity Provider are going to use to allocate IP addresses for the tunnels.
* `ip_pool_v6`: (Optional) The IPv6 Pool ID the users in this Identity Provider are going to use to allocate IP addresses for the tunnels.
* `claim_mappings`: (Optional) The mapping of Identity Provider attributes to claims.
  * `attribute_name`: (Required) The name of the attribute coming from the Identity Provider.
  * `claim_name`: (Required) The name of the user claim to be used in Appgate SDP.
  * `list`:  (Optional)  default value `false` Whether the claim is expected to be a list and have multiple values or not.
  * `encrypt`:  (Optional)  default value `false` Whether the claim should be encrypt or not.
* `user_scripts`: (Optional) ID of the User Claim Scripts to run during authorization.

### Configurable Identity Provider Arguments
`admin_provider`: (Optional) Whether the provider will be listed in the Admin UI or not.

* `dns_servers`: (Optional) The DNS servers to be assigned to the Clients of the users in this Identity Provider.
* `dns_search_domains`: (Optional) The DNS search domains to be assigned to Clients of the users in this Identity Provider.
* `device_limit_per_user`:  (Optional) The device limit per user. The existing on-boarded devices will still be able to sign in even if the limit is exceeded. Deprecated. Use root level field instead.
* `onboarding2_fa`: (Optional) On-boarding two-factor authentication settings. Leave it empty keep it disabled.
  * `mfa_provider_id`: (Required) MFA provider ID to use for the authentication.
  * `message`:  (Optional) On-boarding MFA message to be displayed on the Client UI during the second-factor authentication. Example: Please use your multi factor authentication device to on-board..
  * `claim_suffix`:  (Optional)  default value `onBoarding` Upon successful on-boarding, the claim will be added as if MFA remedy action is fulfilled.
  * `always_required`:  (Optional) If enabled, MFA will be required on every authentication.
* `inactivity_timeout_minutes`: (Optional) (Desktop) clients will sign out automatically after the user has been inactive on the device for the configured duration. Set it to 0 to disable.
* `network_inactivity_timeout_enabled`: (Optional) Whether or not to take network inactivity into account when measuring client inactivity timeout.
* `enforce_windows_network_profile_as_domain`: (Optional) If enabled, Windows Client will configure the network profile as "DomainAuthenticated".
* `block_local_dns_requests`: (Optional) Whether the Windows Client will block local DNS requests or not.
* `on_demand_claim_mappings`: (Optional) The mapping of Identity Provider on demand attributes to claims.
  * `command`: (Required)  Enum values: `fileSize,fileExists,fileCreated,fileUpdated,fileVersion,fileSha512,processRunning,processList,serviceRunning,serviceList,regExists,regQuery,runScript`The name of the command.
  * `claim_name`: (Required) The name of the device claim to be used in Appgate SDP.
  * `parameters`:  (Optional) Depending on the command type, extra parameters to pass to the on-demand claim.
  * `platform`: (Required)  Enum values: `desktop.windows.all,desktop.macos.all,desktop.linux.all,desktop.all,mobile.android.all,mobile.ios.all,mobile.all,all`The platform(s) to run the on-demand claim.

### Oidc Identity Provider Specific Arguments
* `issuer`: (Optional) OIDC issuer URL.
* `audience`: (Optional) Audience/Client ID to make sure the recipient of the Token is this Controller.
* `scope`: (Optional) Scope to use for tokens.
* `google`: (Optional): Google specific OIDC settings
  * `enabled`: (Optional) Whether to enable Google specific OIDC settings
  * `client_secret`: (Optional) Client secret.
  * `refresh_token`: (Optional) Whether to enable refresh token with Google OIDC or not

## Import
Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_oidc_identity_provider.example d3131f83-10d1-4abc-ac0b-7349538e8300
```
