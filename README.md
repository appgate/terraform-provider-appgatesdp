# Appgate SDP Terraform Provider

This repository contains the official Terraform provider for [Appgate SDP](https://www.appgate.com/software-defined-perimeter), enabling you to manage your SDP infrastructure as code.

## 🔧 Purpose

Our goal is to provide first-class support for **the latest version of Appgate SDP**, with compatibility and maintenance extending to the **two most recent versions** as well.

## 📦 Versioning Strategy

Each release of the provider follows the format:

```
v1.<API_VERSION>.<PATCH_VERSION>
```

Where:

* `API_VERSION` corresponds to the Appgate SDP API version (e.g., `22` for SDP 6.5),
* `PATCH_VERSION` is the patch release for the provider.

### 🔍 Finding the Right Version

To use a version of the provider that matches your Appgate SDP deployment:

1. Visit the [GitHub Releases page](https://github.com/appgate/terraform-provider-appgatesdp/releases).
2. Look for a version with the appropriate `API_VERSION`.

For example, if you're using Appgate SDP 6.4 (API v21), you should use:

```
v1.21.x
```

### ✅ Supported Versions

We actively maintain patches for the latest **three** SDP versions:

| Appgate SDP Version | API Version | Terraform Provider |
| ------------------- | ----------- | ------------------ |
| 6.5 (latest)        | v22         | `v1.22.x`          |
| 6.4                 | v21         | `v1.21.x`          |
| 6.3                 | v20         | `v1.20.x`          |

> Earlier versions may still be available, but they are not guaranteed to receive further updates or support.

---

## 🤝 Contributing

We welcome contributions from the community!

### Bug Fixes

* All bug fixes should be made against the `master` branch first.
* Once merged, they will be backported to each of the currently supported release branches (`release-1.22`, `release-1.21`, `release-1.20`).

### Adding Support for a New API Version

* All changes to support a new Appgate SDP version should be implemented in `master`.
* After full support is verified and merged, create a new release branch in the form of `release-1.<API_VERSION>` (e.g., `release-1.23` for a future API version 23).

Please open a pull request and include relevant testing or examples when possible.
