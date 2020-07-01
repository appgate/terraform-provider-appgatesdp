---
layout: "appgate"
page_title: "APPGATE: appgate_entitlement"
sidebar_current: "docs-appgate-datasource-entitlement"
description: |-
  The entitlement data source provides details about a specific entitlement.
---

# appgate_entitlement

The entitlement data source provides details about a specific entitlement.


## Example Usage

```hcl

variable "entitlement_id" {}

data "appgate_entitlement" "default_entitlement" {
    entitlement_id = "${var.entitlement_id}"
}

```

## Argument Reference

* entitlement_id - (Optional) ID of entitlement
* entitlement_name - (Optional) Name of entitlement
