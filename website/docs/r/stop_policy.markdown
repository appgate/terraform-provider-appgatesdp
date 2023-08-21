---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_stop_policy"
sidebar_current: "docs-appgate-resource-stop-policy"
description: |-
   Create a new Stop Policy.
---

# appgatesdp_stop_policy
Stop Policy 

## Example Usage
```hcl
resource "appgatesdp_stop_policy" "test_stop_policy" {
	name = "test stop policy"
	tags = ["aa", "bb", "cc"]
	client_profile_settings {
	  enabled = false
	  force = false
	  profiles = []
	}
}
```

## Argument Reference
The following arguments are supported:
* `disabled`: (Optional) If true, the Policy will be disregarded during authorization.
* `expression`: (Required) A JavaScript expression that returns boolean. Criteria Scripts may be used by calling them as functions.
* `type`: (Computed) Type of the Policy. It is informational and not enforced.
* `policy_id`: (Computed) ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.

### tags
Array of tags.

### client_profile_settings
* `force`: (Optional) Makes the client skip the user prompt and apply the profiles immediately. Required to be true to apply the settings when authorization fails, such as in case of Stop Policies.
* `enabled`: (Optional) Enable Client Profile Settings for this Policy.

## Import
Instances can be imported using the `id`, e.g.
```
$ terraform import appgatesdp_admin_policy.example d3131f83-10d1-4abc-ac0b-7349538e8300
```
