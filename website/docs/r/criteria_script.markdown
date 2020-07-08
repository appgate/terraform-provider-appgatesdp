---
layout: "appgate"
page_title: "APPGATE: appgate_criteria_script"
sidebar_current: "docs-appgate-resource-criteria_script"
description: |-
   Create a new Criteria Script.
---

# appgate_criteria_script

Create a new Criteria Script..

## Example Usage

```hcl

resource "appgate_criteria_script" "test_criteria_script" {

}

```

## Argument Reference

The following arguments are supported:


* `expression`: (Required) A JavaScript expression that returns boolean.
* `id`: (Required) ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.





## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgate_criteria_script d3131f83-10d1-4abc-ac0b-7349538e8300
```
