---
layout: "appgate"
page_title: "APPGATE: appgate_device_script"
sidebar_current: "docs-appgate-datasource-device_script"
description: |-
  The device_script data source provides details about a specific device_script.
---

# appgate_device_script

The device_script data source provides details about a specific device_script.


## Example Usage

```hcl

variable "device_script_id" {}

data "appgate_device_script" "default_device_script" {
    device_script_id = "${var.device_script_id}"
}

```

## Argument Reference

* device_script_id - (Optional) ID of device_script
* device_script_name - (Optional) Name of device_script
