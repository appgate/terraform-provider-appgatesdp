---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_entitlement"
sidebar_current: "docs-appgate-resource-entitlement"
description: |-
   Create a new Entitlement.
---

# appgatesdp_entitlement

Create a new Entitlement.

~> **NOTE:**  The resource documentation is based on the latest available appgate sdp appliance version, which currently is 6.2
Some attributes may not be available if you are running an older version, if you try to use an attribute block that is not permitted in your current version, you will be prompted by an error message.


## Example Usage

```hcl


resource "appgatesdp_entitlement" "ping_entitlement" {
  name = "test entitlement"
  site = data.appgatesdp_site.default_site.id
  conditions = [
    data.appgatesdp_condition.always.id
  ]

  tags = [
    "terraform",
    "api-created"
  ]
  disabled = true

  condition_logic = "and"
  actions {
    action  = "allow"
	subtype = "tcp_up"
	hosts = [
	  "103.15.3.254/32",
	  "172.17.3.255/32",
	  "192.168.2.255/32",
	]
	ports   = ["53"]
  }
  actions {
    action  = "allow"
	subtype = "udp_up"
	hosts   = ["192.168.2.255/32"]
	ports   = ["53"]
  }

  app_shortcuts {
    name       = "ping"
    url        = "https://www.google.com"
    color_code = 5
  }

  app_shortcut_scripts = [
    "313464a6-9dcb-4c6e-90fc-28dceaecb0a1"
  ]

}



```

## Argument Reference

The following arguments are supported:


* `disabled`: (Optional) If true, the Entitlement will be disregarded during authorization.
* `site`: (Required) ID of the Site for this Entitlement.
* `site_name`: (Optional) Name of the Site for this Entitlement. For convenience only.
* `risk_sensitivity`: (Optional) Generate Conditions for the Entitlement based on the Risk Model. Cannot be combined with other Conditions.
* `condition_logic`: (Optional) Whether all the Conditions must succeed to have access to this Entitlement or just one.
* `conditions`: (Required) List of Condition IDs applies to this Entitlement.
* `actions`: (Required) List of all IP Access actions in this Entitlement.
* `app_shortcuts`: (Optional) Array of App Shortcuts.
* `app_shortcut_scripts`: (Optional) List of Entitlement Script IDs used for creating App Shortcuts dynamically.
* `entitlement_id`: (Optional) Computed if empty -  ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.


### conditions
List of Condition IDs applies to this Entitlement.

### actions
List of all IP Access actions in this Entitlement.

* `subtype`:  (Optional)  Enum values: `icmp_up,icmp_down,icmpv6_up,icmpv6_down,udp_up,udp_down,tcp_up,tcp_down,ah_up,ah_down,esp_up,esp_down,gre_up,gre_down,http_up`Type of the IP Access action. Required the action is exclude.
* `action`: (Required)  Enum values: `allow,block,alert,exclude`Applied action to the traffic.
* `hosts`: (Required) Hosts to apply the action to. See admin manual for possible values.
* `ports`:  (Optional) Destination port. Multiple ports can be entered comma separated. Port ranges can be entered dash separated. Only valid for tcp and udp subtypes
* `types`:  (Optional) ICMP type. Only valid for icmp subtypes.
* `methods`:  (Optional) HTTP method. Only valid for http subtypes. Leave it empty to allow all types.
* `monitor`:  (Optional) Only available for tcp_up and http_up subtypes. If enabled, Gateways will monitor this action for responsiveness and act accordingly. See admin manual for more details.

### app_shortcuts
Array of App Shortcuts.

* `name`: (Required) Name for the App Shortcut which will be visible on the Client UI.
* `description`:  (Optional) Description for the App Shortcut which will be visible on the Client UI.
* `url`: (Required) The URL that will be triggered on the OS to be handled. For example, an HTTPS URL will start the browser for the given URL.
* `color_code`:  (Optional)  default value `1` The code of the published app on the client.
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

### app_shortcut_scripts
List of Entitlement Script IDs used for creating App Shortcuts dynamically.

### tags
Array of tags.

## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_entitlement.example d3131f83-10d1-4abc-ac0b-7349538e8300
```
