---
layout: "appgate"
page_title: "APPGATE: appgate_entitlement_script"
sidebar_current: "docs-appgate-resource-entitlement_script"
description: |-
   Create a new Entitlement Script.
---

# appgate_entitlement_script

Create a new Entitlement Script.

## Example Usage

```hcl

resource "appgate_entitlement_script" "test_entitlement_script" {
  name       = "app_shortcut_script"
  type       = "appShortcut"
  expression = <<-EOF
/* This should return an array of app shortcuts. See example below for details:

return [{
  name: '', // Optimally up to 12 characters (maximum 32)
  description: '', // Optimally up to 120 characters (maximum 300)
  url: '',
  colorCode: 1 // Number between 1-20. Optional, leave empty and a random color will be assigned.
}];

*/

return [];
EOF

  tags = [
    "terraform"
  ]
}
```

## Argument Reference

The following arguments are supported:


* `type`: (Optional) The type of the field to use the script for.
* `expression`: (Required) A JavaScript expression that returns a list of IPs and names.
* `id`: (Required) ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.


### tags
Array of tags.




## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgate_entitlement_script d3131f83-10d1-4abc-ac0b-7349538e8300
```
