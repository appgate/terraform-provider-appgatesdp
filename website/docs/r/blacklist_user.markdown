---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_blacklist_user"
sidebar_current: "docs-appgate-resource-blacklist_user"
description: |-
   Blacklists a User.
---

# appgatesdp_blacklist_user

Blacklists a User.

## Example Usage

```hcl

resource "appgatesdp_blacklist_user" "user" {
  user_distinguished_name = "CN=JohnDoe,OU=ldap"
}

```

## Argument Reference

The following arguments are supported:


* `blacklisted_at`: (Optional) The date and time of the blacklisting.
* `reason`: (Optional) The reason for blacklisting. The value is stored and logged.
* `user_distinguished_name`: (Optional) Distinguished name of a user. Format: "CN=,OU="
* `username`: (Optional) The username, same as the one in the user Distinguished Name.
* `provider_name`: (Optional) The provider name of the user, same as the one in the user Distinguished Name.





## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_blacklist_user d3131f83-10d1-4abc-ac0b-7349538e8300
```
