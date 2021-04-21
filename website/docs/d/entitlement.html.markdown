---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_entitlement"
sidebar_current: "docs-appgate-datasource-entitlement"
description: |-
  The entitlement data source provides details about a specific entitlement.
---

# appgatesdp_entitlement

The entitlement data source provides details about a specific entitlement.


## Example Usage

```hcl

variable "entitlement_id" {}

data "appgatesdp_entitlement" "default_entitlement" {
    entitlement_id = "${var.entitlement_id}"
}

```

## Argument Reference

* entitlement_id - (Optional) ID of entitlement
* entitlement_name - (Optional) Name of entitlement
