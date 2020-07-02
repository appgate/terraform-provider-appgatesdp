---
layout: "appgate"
page_title: "APPGATE: appgate_ringfence_rule"
sidebar_current: "docs-appgate-datasource-ringfence_rule"
description: |-
  The ringfence-rule data source provides details about a specific ringfence-rule.
---

# appgate_ringfence_rule

The ringfence-rule data source provides details about a specific ringfence-rule.


## Example Usage

```hcl

variable "ringfence_rule_id" {}

data "appgate_ringfence_rule" "default_ringfence_rule" {
    ringfence-rule_id = "${var.ringfence_rule_id}"
}

```

## Argument Reference

* ringfence_rule_id - (Optional) ID of ringfence-rule
* ringfence_rule_name - (Optional) Name of ringfence-rule
