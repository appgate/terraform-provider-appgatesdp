---
layout: "appgate"
page_title: "APPGATE: appgate_entitlement"
sidebar_current: "docs-appgate-resource-entitlement"
description: |-
   Create a new Entitlement.
---

# appgate_entitlement

Create a new Entitlement..

## Example Usage

```hcl


data "appgate_site" "default_site" {
  site_name = "Default site"
}




resource "appgate_entitlement" "ping_entitlement" {
  name = "test entitlement"
  site = data.appgate_site.default_site.id
  # site = appgate_site.gbg_site.id
  conditions = [
    data.appgate_condition.always.id
  ]

  actions {
    subtype = "icmp_up"
    action  = "allow"
    # https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml#icmp-parameters-types
    types = ["0-16"]
    hosts = [
      "10.0.0.1",
      "10.0.0.0/24",
      "hostname.company.com",
      "dns://hostname.company.com",
      "aws://security-group:accounting"
    ]
  }

}


```

## Argument Reference

The following arguments are supported:


* `display_name`: (Required) This field is deprecated as of 5.1 in favor of &#39;appShortcut.name&#39;. For backwards compatibility, it will set &#39;appShortcut.name&#39; if it does not exist.
* `disabled`: (Optional) If true, the Entitlement will be disregarded during authorization.
* `site`: (Required) ID of the site for this Entitlement.
* `condition_logic`: (Optional) Whether all the Conditions must succeed to have access to this Entitlement or just one.
* `conditions`: (Required) List of Condition IDs applies to this Entitlement.
* `actions`: (Required) List of all IP Access actions in this Entitlement.
* `app_shortcut`: (Optional) Publishes the configured URL as an app on the client using the display name as the app name.


### app_shortcut
Publishes the configured URL as an app on the client using the display name as the app name.

* `name`: Name for the App Shortcut which will be visible on the Client UI. Example: Accounting.
* `url`: The URL that will be triggered on the OS to be handled. For example, an HTTPS URL will start the browser for the given URL. Example: https:&#x2F;&#x2F;service.company.com.
* `color_code`: The code of the published app on the client.
- 1: Light Green
- 2: Green
- 3: Indigo
- 4: Deep Purple
- 5: Yellow
- 6: Lime
- 7: Light Blue
- 8: Blue
- 9: Amber
- 10: Orange
- 11: Cyan
- 12: Teal
- 13: Deep Orange
- 14: Red
- 15: Gray
- 16: Brown
- 17: Pink
- 18: Purple
- 19: Blue Gray
- 20: Near Black





## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgate_entitlement d3131f83-10d1-4abc-ac0b-7349538e8300
```
