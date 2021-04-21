---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_device_script"
sidebar_current: "docs-appgate-resource-device_script"
description: |-
   Create a new Device Script.
---

# appgatesdp_device_script

Create a new Device Script.

## Example Usage


### Inline content script
```hcl

resource "appgatesdp_device_script" "example_device_script" {
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

### Upload device script from file path

```hcl

resource "appgatesdp_device_script" "example_device_script" {
  name     = "device_script_name"
  filename = "script.sh"
  file     = "/path/to/file/script.sh"
}

```

## Argument Reference

The following arguments are supported:


* `filename`: (Required) The name of the file to be downloaded as to the client devices.
* `file`: (Optional) The Device Claim Script binary in Base64 format.
* `checksum`: (Optional) MD5 checksum of the file. It's used by the Client to decide whether to download the script again or not. Deprecated as of 5.0. Use checksumSha256 field.
* `checksum_sha256`: (Optional) SHA256 checksum of the file. It's used by the Client to decide whether to download the script again or not.
* `id`: (Required) ID of the object.
* `name`: (Required) Name of the object.
* `notes`: (Optional) Notes for the object. Used for documentation purposes.
* `tags`: (Optional) Array of tags.


### tags
Array of tags.




## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_device_script d3131f83-10d1-4abc-ac0b-7349538e8300
```
