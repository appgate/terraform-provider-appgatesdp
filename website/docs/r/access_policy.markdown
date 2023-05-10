---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_access_policy"
sidebar_current: "docs-appgate-resource-access-policy"
description: |-
   Create a new Access Policy.
---

# appgatesdp_access_policy

Create a new Access Policy.



## Example Usage

```hcl


resource "appgatesdp_access_policy" "test_access_policy" {
	name = "Test access policy"
	tags = ["aa", "bb", "cc"]
	entitlements = [
		appgatesdp_entitlement.one.id,
		appgatesdp_entitlement.two.id,
	]
	override_site = data.appgatesdp_site.default_site.id
}


```


## Argument Reference

The following arguments are supported:


* `disabled`: (Optional) If true, the Policy will be disregarded during authorization.
* `expression`: (Required) A JavaScript expression that returns boolean. Criteria Scripts may be used by calling them as functions.
* `type`: (Computed) Type of the Policy. It is informational and not enforced.
* `entitlements`: (Optional) List of Entitlement IDs in this Policy.
* `entitlement_links`: (Optional) List of Entitlement tags in this Policy.
* `policy_id`: (Computed) ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.
* `override_site`: (Optional) Site ID where all the Entitlements of this Policy must be deployed. This overrides Entitlement's own Site and to be used only in specific network layouts. Otherwise the assigned site on individual Entitlements will be used.
* `override_site_claim`: (Optional) The path of a claim that contains the UUID of an override site. It should be defined as "claims.xxx.xxx" or "claims.xxx.xxx.xxx".

### entitlements
List of Entitlement IDs in this Policy.

### entitlement_links
List of Entitlement tags in this Policy.



### tags
Array of tags.




## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_access_policy.example d3131f83-10d1-4abc-ac0b-7349538e8300
```
