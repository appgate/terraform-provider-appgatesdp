---
layout: "appgate"
page_title: "APPGATE: appgate_entitlement_script"
sidebar_current: "docs-appgate-resource-entitlement_script"
description: |-
   Create a new Entitlement Script.
---

# appgate_entitlement_script

Create a new Entitlement Script..

## Example Usage

```hcl

resource "appgate_entitlement_script" "test_entitlement_script" {

}

```

## Argument Reference

The following arguments are supported:


* `expression`: (Required) A JavaScript expression that returns a list of IPs and names.
* `id`: (Required) ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.





## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgate_entitlement_script d3131f83-10d1-4abc-ac0b-7349538e8300
```
