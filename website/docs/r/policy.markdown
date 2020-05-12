---
layout: "appgate"
page_title: "APPGATE: appgate_policy"
sidebar_current: "docs-appgate-resource-policy"
description: |-
   Create a new Policy.
---

# appgate_policy

Create a new Policy..

## Example Usage

```hcl

resource "appgate_policy" "test_policy" {

}

```

## Argument Reference

The following arguments are supported:


* `disabled`: (Optional) If true, the Policy will be disregarded during authorization.
* `expression`: (Required) A JavaScript expression that returns boolean. Criteria Scripts may be used by calling them as functions.
* `entitlements`: (Optional) List of Entitlement IDs in this Policy.
* `entitlement_links`: (Optional) List of Entitlement tags in this Policy.
* `ringfence_rules`: (Optional) List of Ringfence Rule IDs in this Policy.
* `ringfence_rule_links`: (Optional) List of Ringfence Rule tags in this Policy.
* `tamper_proofing`: (Optional) Will enable Tamper Proofing on desktop clients which will make sure the routes and ringfence configurations are not changed.
* `override_site`: (Optional) Site ID where all the Entitlements of this Policy must be deployed. This overrides Entitlement&#39;s own Site and to be used only in specific network layouts. Otherwise the assigned site on individual Entitlements will be used.
* `administrative_roles`: (Optional) List of Administrative Role IDs in this Policy.





## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgate_policy d3131f83-10d1-4abc-ac0b-7349538e8300
```
