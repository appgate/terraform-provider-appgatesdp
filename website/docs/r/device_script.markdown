---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_device_script"
sidebar_current: "docs-appgate-resource-device_script"
description: |-
   Create a new Device Claim Script.
---

# appgatesdp_device_script

Create a new Device Claim Script.

~> **NOTE:**  The resource documentation is based on the latest available appgate sdp appliance version, which currently is 5.5.0
Some attributes may not be available if you are running an older version, if you try to use an attribute block that is not permitted in your current version, you will be prompted by an error message.


## Example Usage

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
## Example with data source
```hcl
## Example with data source

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
$ terraform import appgatesdp_device_script.example d3131f83-10d1-4abc-ac0b-7349538e8300
```
