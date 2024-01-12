---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_admin_mfa_settings"
sidebar_current: "docs-appgate-resource-admin_mfa_settings"
description: |-
   Update Admin MFA settings.
---

# appgatesdp_admin_mfa_settings

Update Admin MFA settings.

~> **NOTE:**  The resource documentation is based on the latest available appgate sdp appliance version, which currently is 6.2
Some attributes may not be available if you are running an older version, if you try to use an attribute block that is not permitted in your current version, you will be prompted by an error message.


## Example Usage

```hcl


resource "appgatesdp_admin_mfa_settings" "mfa_settings" {
  exempted_users = [
    "CN=Dan,OU=local",
  ]
}

```
## Example with data source
```hcl
## Example with data source

data "appgatesdp_mfa_provider" "fido" {
  mfa_provider_name = "Default FIDO2 Provider"
}
data "appgatesdp_identity_provider" "local" {
  identity_provider_name = "local"
}

resource "appgatesdp_admin_mfa_settings" "mfa_settings" {
  provider_id = data.appgatesdp_mfa_provider.fido.mfa_provider_id
  exempted_users = [
    "CN=Jane,OU=local",
    format("CN=Joe,OU=%s", data.appgatesdp_identity_provider.local.identity_provider_name),
  ]
}


```


## Argument Reference

The following arguments are supported:


* `provider_id`: (Optional) The MFA provider ID to use during Multi-Factor Authentication. If null, Admin MFA is disabled.
* `exempted_users`: (Optional) List of users to be excluded from MFA during admin login.


### exempted_users
List of users to be excluded from MFA during admin login.




## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_admin_mfa_settings.example d3131f83-10d1-4abc-ac0b-7349538e8300
```
