---
layout: "appgatesdp"
page_title: "APPGATE: appgatesdp_license"
sidebar_current: "docs-appgate-resource-license"
description: |-
   Upload a new License and override the existing one.
---

# appgatesdp_license

Upload a new License.

## Example Usage

```hcl

resource "appgatesdp_license" "the_license" {
    license  = <<-EOF
....
....
....
EOF
}

```

## Argument Reference

The following arguments are supported:

* `license`: (Required) The license file contents for this Controller (with the matching request code).





## Import

Instances can be imported using the `id`, e.g.

```
$ terraform import appgatesdp_license d3131f83-10d1-4abc-ac0b-7349538e8300
```
