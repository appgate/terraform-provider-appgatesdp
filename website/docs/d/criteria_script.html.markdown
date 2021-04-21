---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_criteria_script"
sidebar_current: "docs-appgate-datasource-criteria_script"
description: |-
  The criteria_script data source provides details about a specific criteria_script.
---

# appgatesdp_criteria_script

The criteria_script data source provides details about a specific criteria_script.


## Example Usage

```hcl

variable "criteria_script_id" {}

data "appgatesdp_criteria_script" "default_criteria_script" {
    criteria_script_id = "${var.criteria_script_id}"
}

```

## Argument Reference

* criteria_script_id - (Optional) ID of criteria_script
* criteria_script_name - (Optional) Name of criteria_script
