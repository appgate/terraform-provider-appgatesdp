---
layout: "appgate"
page_title: "APPGATE: appgate_ringfence-rule"
sidebar_current: "docs-appgate-resource-ringfence-rule"
description: |-
   Create a new Ringfence Rule.
---

# appgate_ringfence-rule

Create a new Ringfence Rule..

## Example Usage

```hcl

resource "appgate_ringfence_rule" "basic_rule" {
  name = "basic"
  tags = [
    "terraform",
    "api-created"
  ]

  actions {
    protocol  = "icmp"
    direction = "out"
    action    = "allow"

    hosts = [
      "10.0.2.0/24"
    ]

    ports = [
      "80",
      "443",
      "1024-2048"
    ]

    types = [
      "0-255"
    ]

  }

  actions {
    protocol  = "tcp"
    direction = "in"
    action    = "allow"

    hosts = [
      "10.0.2.0/24"
    ]

    ports = [
      "22-25"
    ]
  }

}
```

## Argument Reference

The following arguments are supported:


* `actions`: (Required) List of all ringfence actions in this Ringfence Rule.
* `id`: (Required) ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `created`: (Optional) Create date.
* `updated`: (Optional) Last update date.
* `tags`: (Optional) Array of tags.


### actions
List of all ringfence actions in this Ringfence Rule.

* `protocol`: (Required)  Enum values: `icmp,icmpv6,udp,tcp`Protocol of the ringfence action.
* `direction`: (Required)  Enum values: `up,down`The direction of the action
* `action`: (Required)  Enum values: `allow,block`Applied action to the traffic.
* `hosts`: (Required) Destination address. IP address or hostname.
* `ports`:  (Optional) Destination port. Multiple ports can be entered comma separated. Port ranges can be entered dash separated. Only valid for tcp and udp subtypes. Example: 80,443,1024-2048.
* `types`:  (Optional) ICMP type. Only valid for icmp protocol.
#### hosts
Destination address. IP address or hostname.
#### ports
Destination port. Multiple ports can be entered comma separated. Port ranges can be entered dash separated. Only valid for tcp and udp subtypes.
#### types
ICMP type. Only valid for icmp protocol.



## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgate_ringfence-rule d3131f83-10d1-4abc-ac0b-7349538e8300
```
