---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_client_connections"
sidebar_current: "docs-appgate-resource-client_connections"
description: |-
   Update Client Connection settings.
---

# appgatesdp_client_connections

Update Client Connection settings.

~> **NOTE:**  The resource documentation is based on the latest available appgate sdp appliance version, which currently is 5.5.0
Some attributes may not be available if you are running an older version, if you try to use an attribute block that is not permitted in your current version, you will be prompted by an error message.


## Example Usage

```hcl


resource "appgatesdp_client_connections" "cc" {
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


* `spa_mode`: (Optional) SPA mode. Deprecated as of 5.4. Use global-settings API instead.
* `profile_hostname`: (Optional) The hostname to use for generating profile URLs.
* `profiles`: (Optional) Client Profiles.


### profiles
Client Profiles.

* `name`: (Required) A name to identify the client profile. It will appear on the client UI.
* `spa_key_name`: (Required) SPA key name to be used in the profile. Same key names in different profiles will have the same SPA key. SPA key is used by the client to connect to the controllers.
* `identity_provider_name`: (Required) Name of the Identity Provider to be used to authenticate.
* `url`:  (Optional) Connection URL for the profile.



## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_client_connections d3131f83-10d1-4abc-ac0b-7349538e8300
```
