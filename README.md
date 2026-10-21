# Appgate SDP Terraform Provider

Terraform provider for [Appgate SDP](https://www.appgate.com/products/zero-trust-network-access).
It is generated from the official
[SDP Admin API OpenAPI specifications](https://github.com/appgate/sdp-api-specification) with
HashiCorp's [provider code generation](https://developer.hashicorp.com/terraform/plugin/code-generation)
tools. One release supports every listed SDP version.

| SDP version | Admin API | Specification branch |
|---|---|---|
| 7.0 | v25 | `version-25` |
| 6.7 | v24 | `version-24` |
| 6.6 | v23 | `version-23` |
| 6.5 | v22 | `version-22` |
| 6.4 | v21 | `version-21` |

Pinned specification commits are listed in
[`internal/apiversion/versions.json`](internal/apiversion/versions.json) and in the
[API versions guide](docs/guides/api-versions.md).

```hcl
provider "appgatesdp" {
  url         = "https://controller.example.com:8443/admin"
  username    = var.username
  password    = var.password
  api_version = "v25" # optional: negotiated with the controller when unset
}
```

* User documentation: [`docs/`](docs/index.md)
* Upgrading from 1.x: [migration guide](docs/guides/migrating-to-v2.md)
* Design, generation pipeline, adding SDP versions: [ARCHITECTURE.md](ARCHITECTURE.md)

## Development

Requirements:

* Go (version in `go.mod`) and `git`
* the [Terraform CLI](https://developer.hashicorp.com/terraform/install), used by unit tests to
  drive plan/apply against an in-memory controller and by `make docs`
* [golangci-lint](https://golangci-lint.run/) v2 for `make lint`

The code generators (`tfplugingen-openapi`, `tfplugingen-framework`) and `tfplugindocs` are
pinned in `go.mod` and run through `go tool`. Nothing else needs to be installed.

| Command | Purpose |
|---|---|
| `make generate` | Fetch the pinned specifications and regenerate all versions |
| `make generate API_VERSION=v25` | Regenerate one version; the other versions are not touched |
| `make pin API_VERSION=v25` | Pin a version to the current head of its specification branch |
| `make fmt` / `make fmt-check` | Format with `gofmt` / fail on unformatted files without changing them |
| `make vet` | `go vet ./...` |
| `make lint` | `golangci-lint run ./...` with [`.golangci.yml`](.golangci.yml); generated code included |
| `make build` | `go build ./...` and the provider binary |
| `make test` / `make test-race` | Unit and mocked lifecycle tests, with or without the race detector |
| `make tidy-check` | `go mod tidy -diff` |
| `make generate-check` | Regenerate and fail if committed files change |
| `make docs` / `make docs-check` | Generate the registry documentation with `tfplugindocs` |
| `make check` | All of the above quality gates; used by CI |
| `make testacc` | Acceptance tests against live controllers (see below) |
| `make dev` | Build and install into `~/.terraform.d/plugins` |

### Tests

`make test` needs no SDP environment. Provider tests start a version-aware fake controller
(`internal/sdptest`) and run real Terraform plans, applies, imports and refreshes against it.
They run against every registered API version.

### Acceptance tests against live controllers

Configure one or more environments by API version:

```bash
export SDP_ACC_V25_URL=https://controller-70.example.com:8443/admin
export SDP_ACC_V25_USERNAME=admin
export SDP_ACC_V25_PASSWORD=...
export SDP_ACC_V25_PEM=/path/to/ca.pem   # optional; without it TLS is not verified
# the same for SDP_ACC_V21_* to SDP_ACC_V24_*

make testacc                                   # read-only: compares live API responses with the specification
SDP_ACC_ALLOW_DESTRUCTIVE=1 make testacc       # also creates, updates and deletes test objects
```

Tests only create objects named `tf-acc-*` and remove them afterwards. Run
`SDP_ACC_ALLOW_DESTRUCTIVE=1 make testacc TESTARGS='-run TestAccSweep'` to clean up after an
interrupted run. Destructive tests never run unless `SDP_ACC_ALLOW_DESTRUCTIVE=1` is set.

## Contributing

Never edit generated files by hand: `codegen/ir`, `codegen/reports`,
`internal/apiversion/generated`, `internal/generated`, `internal/provider/zz_generated_*.go` and
`docs/`. Change `codegen/resources.yml`, the normalization in `internal/codegen`, or the runtime
in `internal/provider`, then run `make generate docs check`.
