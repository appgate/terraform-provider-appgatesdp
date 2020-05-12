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

resource "appgate_ringfence-rule" "test_ringfence-rule" {

}

```

## Argument Reference

The following arguments are supported:


* `actions`: (Required) List of all ringfence actions in this Ringfence Rule.


### actions
List of all ringfence actions in this Ringfence Rule.

* `protocol`: Protocol of the ringfence action.
* `direction`: The direction of the action
* `action`: Applied action to the traffic.
* `hosts`: Destination address. IP address or hostname.
* `ports`: Destination port. Multiple ports can be entered comma separated. Port ranges can be entered dash separated. Only valid for tcp and udp subtypes. Example: 80,443,1024-2048.
* `types`: ICMP type. Only valid for icmp protocol.

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
