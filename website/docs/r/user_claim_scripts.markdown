---
layout: "appgate"
page_title: "APPGATE: appgate_user_claim_scripts"
sidebar_current: "docs-appgate-resource-user_claim_scripts"
description: |-
   Create a new User Claim Script.
---

# appgate_user_claim_scripts

Create a new User Claim Script.

## Example Usage

```hcl
resource "appgatesdp_user_claim_script" "custom_script" {
  name       = "updated claim name"
  notes      = "This object has been created for test purposes."
  expression = "return {'posture': 25};"
  tags = [
    "developer",
    "api-created"
  ]
}

```

## Argument Reference

The following arguments are supported:


* `expression`: (Required) A JavaScript expression that returns an object.
* `id`: (Optional) Computed if empty -  ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.


### tags
Array of tags.




## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgate_user_claim_scripts.example d3131f83-10d1-4abc-ac0b-7349538e8300
```
