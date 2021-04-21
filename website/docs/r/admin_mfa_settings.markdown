---
layout: "appgate"
page_title: "APPGATE: appgatesdp_admin_mfa_settings"
sidebar_current: "docs-appgate-resource-admin_mfa_settings"
description: |-
   Update Admin MFA settings.
---

# appgatesdp_admin_mfa_settings

Update Admin MFA settings.

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
$ terraform import appgatesdp_admin_mfa_settings.mfa_settings admin_mfa_settings
```
