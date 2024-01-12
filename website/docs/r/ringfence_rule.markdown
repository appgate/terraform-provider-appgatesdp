---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_ringfence_rule"
sidebar_current: "docs-appgate-resource-ringfence_rule"
description: |-
   Create a new Ringfence Rule.
---

# appgatesdp_ringfence_rule

Create a new Ringfence Rule.

~> **NOTE:**  The resource documentation is based on the latest available appgate sdp appliance version, which currently is 6.2
Some attributes may not be available if you are running an older version, if you try to use an attribute block that is not permitted in your current version, you will be prompted by an error message.


## Example Usage

```hcl


resource "appgatesdp_ringfence_rule" "basic_rule" {
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
* `ringfence_rule_id`: (Optional) ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.


### actions
List of all ringfence actions in this Ringfence Rule.

* `protocol`: (Required)  Enum values: `icmp,icmpv6,udp,tcp`Protocol of the ringfence action.
* `direction`: (Required)  Enum values: `up,down,out,in`The direction of the action
* `action`: (Required)  Enum values: `allow,block`Applied action to the traffic.
* `hosts`: (Required) Destination address. IP address or hostname.
* `ports`:  (Optional) Destination port. Multiple ports can be entered comma separated. Port ranges can be entered dash separated. Only valid for tcp and udp subtypes.
* `types`:  (Optional) ICMP type. Only valid for icmp protocol.
### tags
Array of tags.




## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_ringfence_rule.example d3131f83-10d1-4abc-ac0b-7349538e8300
```
