---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_appliance_customization"
sidebar_current: "docs-appgate-resource-appliance_customization"
description: |-
   Create a new Appliance Customization.
---

# appgatesdp_appliance_customization

Create a new Appliance Customization.

https://sdphelp.appgate.com/adminguide/v5.1/appliance-customizations-configure.html

## Example Usage

```hcl

data "archive_file" "customization" {
  type        = "zip"
  output_path = "${path.module}/customization/package.zip"

  source {
    content  = <<-EOF
#!/usr/bin/env bash
echo "startup script"
EOF
    filename = "start"
  }

  source {
    content  = <<-EOF
#!/usr/bin/env bash
echo "stop script"
EOF
    filename = "stop"
  }
}

resource "appgatesdp_appliance_customization" "test_customization" {
  name = "test customization"
  file = data.archive_file.customization.output_path

  tags = [
    "terraform",
    "api-created"
  ]
}

```

## Argument Reference

The following arguments are supported:


* `file`: (Optional) The Appliance Customization binary in Base64 format.
* `checksum`: (Optional) SHA256 checksum of the file.
* `size`: (Optional) Binary file's size in bytes.
* `id`: (Required) ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.


### tags
Array of tags.




## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_appliance_customization d3131f83-10d1-4abc-ac0b-7349538e8300
```
