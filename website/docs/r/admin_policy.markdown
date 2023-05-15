---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_admin_policy"
sidebar_current: "docs-appgate-resource-admin-policy"
description: |-
   Create a new Admin Policy.
---

# appgatesdp_admin_policy

Admin policy rights to perform one or more administrative roles (admin UI & APIs).


## Example Usage

```hcl


resource "appgatesdp_admin_policy" "test_admin_policy" {
	name = "test admin policy"
	tags = ["aa", "bb", "cc"]
	administrative_roles = [
		appgatesdp_administrative_role.first.id,
		appgatesdp_administrative_role.second.id,
	]
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


### administrative_roles
List of Administrative Role IDs in this Policy.



### tags
Array of tags.




## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_admin_policy.example d3131f83-10d1-4abc-ac0b-7349538e8300
```
