---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_appliance_customization"
sidebar_current: "docs-appgate-datasource-appliance_customization"
description: |-
  The appliance_customization data source provides details about a specific appliance_customization.
---

# appgatesdp_appliance_customization

The appliance_customization data source provides details about a specific appliance_customization.


## Example Usage

```hcl

variable "appliance_customization_id" {}

data "appgatesdp_appliance_customization" "default_appliance_customization" {
    appliance_customization_id = "${var.appliance_customization_id}"
}

```

## Argument Reference

* appliance_customization_id - (Optional) ID of appliance_customization
* appliance_customization_name - (Optional) Name of appliance_customization
