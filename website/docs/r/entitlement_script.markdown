---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_entitlement_script"
sidebar_current: "docs-appgate-resource-entitlement_script"
description: |-
   Create a new Entitlement Script.
---

# appgatesdp_entitlement_script

Create a new Entitlement Script.

~> **NOTE:**  The resource documentation is based on the latest available appgate sdp appliance version, which currently is 5.5.0
Some attributes may not be available if you are running an older version, if you try to use an attribute block that is not permitted in your current version, you will be prompted by an error message.


## Example Usage

```hcl


resource "appgatesdp_entitlement_script" "test_entitlement_script" {
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
$ terraform import appgatesdp_entitlement_script.example d3131f83-10d1-4abc-ac0b-7349538e8300
```
