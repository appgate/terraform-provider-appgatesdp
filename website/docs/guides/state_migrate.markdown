---
layout: "appgatesdp"
page_title: "State migration"
sidebar_current: "docs-appgate-guide-state_migrate"
description: |-
  Migrate state from <= 0.4.0 to >= 0.5.0
---

## Migrate state



Prior to version 0.5.0, the project was named `terraform-provider-appgate-sdp`, and the provider was not yet published to registry.terraform.io.
When we published it, we noticed problems with using kebab-case in the name, which forced us to re-name the project to `terraform-provider-appgatesdp`

We can still use the old .tf files and state files, but we need to update the names, we can use the state-migrate tool to do that:

```sh
$ ./state-migrate migrate -dir /path/to/terraform-resources

```
