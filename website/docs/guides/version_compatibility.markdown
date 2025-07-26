---
layout: "appgatesdp"
page_title: "Version compatibility"
sidebar_current: "docs-appgate-guide-version-compatibility"
description: |-
  Version compatibility
---

## Version compatibility

### Support Policy
Appgate SDP maintains full support for the three most recent major versions of the platform. Each supported version includes corresponding compatibility within the Terraform provider.

Always refer to the official Appgate support matrix to verify supported versions:
[SDP Supported Versions](https://www.appgate.com/support/software-defined-perimeter-support)


### Choosing the Right Terraform Provider Version
Each release of the provider follows the format:
```
v1.<API_VERSION>.<PATCH_VERSION>
```
Where:
- `API_VERSION` corresponds to the Appgate SDP API version (e.g., 22 for SDP 6.5),
- `PATCH_VERSION` is the patch release for the provider.

To select the correct Terraform provider version for your environment:
1. Identify your Appgate SDP API version
2. Choose a provider version in the format `v1.<API_VERSION>.<PATCH_VERSION>`

Example:
If your Appgate SDP API version is 22, use a provider version like `v1.22.0`

You can use the following version constraints in your Terraform configuration to pin the API version:

```terraform
terraform {
  required_providers {
    appgatesdp = {
      source  = "appgate/appgatesdp"
      version = ">= 1.22.0, < 1.23.0"
    }
  }
}
```

For additional configuration options, see [example usage](https://registry.terraform.io/providers/appgate/appgatesdp/latest/docs#example-usage).

### Current Support Versions
### ✅ Supported Versions

We currently maintain patches for the latest **three** SDP versions which are the following:

| Appgate SDP Version | API Version | Terraform Provider |
| ------------------- | ----------- | ------------------ |
| 6.5 (latest)        | v22         | `v1.22.x`          |
| 6.4                 | v21         | `v1.21.x`          |
| 6.3                 | v20         | `v1.20.x`          |
