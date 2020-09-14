---
layout: "appgate"
page_title: "APPGATE: appgate_administrative_role"
sidebar_current: "docs-appgate-datasource-administrative_role"
description: |-
  The administrative_role data source provides details about a specific administrative_role.
---

# appgate_administrative_role

The administrative_role data source provides details about a specific administrative_role.


## Example Usage

```hcl

variable "administrative_role_id" {}

data "appgate_administrative_role" "default_administrative_role" {
    administrative_role_id = "${var.administrative_role_id}"
}

```

## Argument Reference

* administrative_role_id - (Optional) ID of administrative_role
* administrative_role_name - (Optional) Name of administrative_role
