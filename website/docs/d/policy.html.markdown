---
layout: "appgate"
page_title: "APPGATE: appgate_policy"
sidebar_current: "docs-appgate-datasource-policy"
description: |-
  The policy data source provides details about a specific policy.
---

# appgate_policy

The policy data source provides details about a specific policy.


## Example Usage

```hcl

variable "policy_id" {}

data "appgate_policy" "default_policy" {
    policy_id = "${var.policy_id}"
}

```

## Argument Reference

* policy_id - (Optional) ID of policy
* policy_name - (Optional) Name of policy
