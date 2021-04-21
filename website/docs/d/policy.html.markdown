---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_policy"
sidebar_current: "docs-appgate-datasource-policy"
description: |-
  The policy data source provides details about a specific policy.
---

# appgatesdp_policy

The policy data source provides details about a specific policy.


## Example Usage

```hcl

variable "policy_id" {}

data "appgatesdp_policy" "default_policy" {
    policy_id = "${var.policy_id}"
}

```

## Argument Reference

* policy_id - (Optional) ID of policy
* policy_name - (Optional) Name of policy
