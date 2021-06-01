---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_client_profile"
sidebar_current: "docs-appgate-resource-client-profile"
description: |-
   Update Client Connection settings. With API version 12, this API has changed significantly in order to manage client profiles. It is still possible to use the older APIs using older Accept headers.
---

# appgatesdp_client_profile

Create or update client profile.


## Example Usage

```hcl

resource "appgatesdp_client_profile" "test_uppercase" {
  name                   = "development"
  spa_key_name           = "development"
  identity_provider_name = "local"
}


```

## Argument Reference

The following arguments are supported:

* `name`: (Required) A name to identify the client profile. It will appear on the client UI.
* `spa_key_name`: (Required) SPA key name to be used in the profile. Same key names in different profiles will have the same SPA key. SPA key is used by the client to connect to the controllers.
* `identity_provider_name`: (Required) Name of the Identity Provider to be used to authenticate.
* `url`:  (Computed) Connection URL for the profile.



## Import

Instances can be imported using the `name`, e.g.

```
$ terraform import appgatesdp_client_profile development
```
