---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_criteria_script"
sidebar_current: "docs-appgate-resource-criteria_script"
description: |-
   Create a new Criteria Script.
---

# appgatesdp_criteria_script

Create a new Criteria Script.

~> **NOTE:**  The resource documentation is based on the latest available appgate sdp appliance version, which currently is 6.0.0
Some attributes may not be available if you are running an older version, if you try to use an attribute block that is not permitted in your current version, you will be prompted by an error message.


## Example Usage

```hcl


resource "appgatesdp_criteria_script" "test_criteria_script" {
  name       = "Test"
  expression = "return claims.user.username === 'admin';"
}


```


## Argument Reference

The following arguments are supported:


* `expression`: (Required) A JavaScript expression that returns boolean.
* `criteria_script_id`: (Optional) ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.


### tags
Array of tags.




## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_criteria_script.example d3131f83-10d1-4abc-ac0b-7349538e8300
```
