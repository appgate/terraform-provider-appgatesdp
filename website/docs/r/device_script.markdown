---
layout: "appgate"
page_title: "APPGATE: appgate_device_script"
sidebar_current: "docs-appgate-resource-device_script"
description: |-
   Create a new Device Script.
---

# appgate_device_script

Create a new Device Script.

## Example Usage

```hcl

resource "appgate_device_script" "example_device_script" {
  name     = "device_script_name"
  filename = "script.sh"
  content  = <<-EOF
#!/usr/bin/env bash
echo "hello world"
EOF
  tags = [
    "terraform",
    "api-created"
  ]
}

```

## Argument Reference

The following arguments are supported:


* `filename`: (Required) The name of the file to be downloaded as to the client devices.
* `file`: (Optional) The Device Script binary path, conflicts with `content`.
* `content`: (Optional) The Device Script content, conflicts with `file`.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.





## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgate_device_script d3131f83-10d1-4abc-ac0b-7349538e8300
```
