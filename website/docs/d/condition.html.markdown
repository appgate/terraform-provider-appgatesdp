---
layout: "appgate"
page_title: "APPGATE: appgate_condition"
sidebar_current: "docs-appgate-datasource-condition"
description: |-
  The condition data source provides details about a specific condition.
---

# appgate_condition

The condition data source provides details about a specific condition.


## Example Usage

```hcl

variable "condition_id" {}

data "appgate_condition" "default_condition" {
    condition_id = "${var.condition_id}"
}

```

## Argument Reference

* condition_id - (Optional) ID of condition
* condition_name - (Optional) Name of condition
