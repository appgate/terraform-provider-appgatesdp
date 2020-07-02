---
layout: "appgate"
page_title: "APPGATE: appgate_appliance"
sidebar_current: "docs-appgate-datasource-appliance"
description: |-
  The appliance data source provides details about a specific appliance.
---

# appgate_appliance

The appliance data source provides details about a specific appliance.


## Example Usage

```hcl

variable "appliance_id" {}

data "appgate_appliance" "default_appliance" {
    appliance_id = "${var.appliance_id}"
}

```

## Argument Reference

* appliance_id - (Optional) ID of appliance
* appliance_name - (Optional) Name of appliance
