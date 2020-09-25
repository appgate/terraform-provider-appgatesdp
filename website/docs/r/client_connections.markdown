---
layout: "appgate"
page_title: "APPGATE: appgate_client_connections"
sidebar_current: "docs-appgate-resource-client_connections"
description: |-
   Update Client Connection settings. With API version 12, this API has changed significantly in order to manage client profiles. It is still possible to use the older APIs using older Accept headers.
---

# appgate_client_connections

Update Client Connection settings.

## Example Usage

```hcl

resource "appgate_client_connections" "cc" {
  spa_mode = "TCP"
  profiles {
    name                   = "Company Employee"
    spa_key_name           = "test_key"
    identity_provider_name = "local"
  }
}


```

## Argument Reference

The following arguments are supported:


* `spa_mode`: (Optional) SPA mode.
* `profiles`: (Optional) Client Profiles.


### profiles
Client Profiles.

* `name`:  (Optional) A name to identify the client profile. It will appear on the client UI.
* `spa_key_name`:  (Optional) SPA key name to be used in the profile. Same key names in different profiles will have the same SPA key. SPA key is used by the client to connect to the controllers.
* `identity_provider_name`:  (Optional) Name of the Identity Provider to be used to authenticate.
* `url`:  (Computed) Connection URL for the profile.



## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgate_client_connections d3131f83-10d1-4abc-ac0b-7349538e8300
```
