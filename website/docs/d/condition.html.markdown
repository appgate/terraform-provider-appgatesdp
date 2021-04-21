---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_condition"
sidebar_current: "docs-appgate-datasource-condition"
description: |-
  The condition data source provides details about a specific condition.
---

# appgatesdp_condition

The condition data source provides details about a specific condition.


## Example Usage

```hcl

variable "condition_id" {}

data "appgatesdp_condition" "default_condition" {
    condition_id = "${var.condition_id}"
}

```

## Argument Reference

* condition_id - (Optional) ID of condition
* condition_name - (Optional) Name of condition
