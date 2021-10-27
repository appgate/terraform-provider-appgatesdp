---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_license"
sidebar_current: "docs-appgate-resource-license"
description: |-
   Upload a new License and override the existing one.
---

# appgatesdp_license

Upload a new License.

~> **NOTE:**  The resource documentation is based on the latest available appgate sdp appliance version, which currently is 5.5.0
Some attributes may not be available if you are running an older version, if you try to use an attribute block that is not permitted in your current version, you will be prompted by an error message.


## Example Usage

```hcl


resource "appgatesdp_license" "the_license" {
  license = <<-EOF
....
....
....
EOF
}


```


## Argument Reference

The following arguments are supported:







## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_license d3131f83-10d1-4abc-ac0b-7349538e8300
```
