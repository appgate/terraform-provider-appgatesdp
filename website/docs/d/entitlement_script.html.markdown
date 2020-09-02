---
layout: "appgate"
page_title: "APPGATE: appgate_entitlement_script"
sidebar_current: "docs-appgate-datasource-entitlement_script"
description: |-
  The entitlement_script data source provides details about a specific entitlement_script.
---

# appgate_entitlement_script

The entitlement_script data source provides details about a specific entitlement_script.


## Example Usage

```hcl

variable "entitlement_script_id" {}

data "appgate_entitlement_script" "default_entitlement_script" {
    entitlement_script_id = "${var.entitlement_script_id}"
}

```

## Argument Reference

* entitlement_script_id - (Optional) ID of entitlement_script
* entitlement_script_name - (Optional) Name of entitlement_script
