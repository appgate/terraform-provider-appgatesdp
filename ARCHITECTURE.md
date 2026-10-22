# Architecture

This provider is generated from the SDP Admin API OpenAPI specifications using
HashiCorp's provider code generation tools. One release supports every registered
Admin API version, and adding a new SDP release is a generation and validation task.

## Overview

```
internal/apiversion/versions.json        registry: API version → SDP version, branch, pinned commit
codegen/resources.yml                    explicit definitions, names and exclusions that tune discovery
            │
            ▼  make generate  (cmd/sdpgen → internal/codegen)
┌───────────────────────────────────────────────────────────────────────────────────────────┐
│ per API version                                                                           │
│  1. fetch      git: branch version-N exists, pinned commit is on it, tree declares API N  │
│  1b. discover  classify every path: resource, data source, action or excluded             │
│  2. normalize  resolve multi-file $refs, flatten allOf/oneOf/discriminators, split        │
│                request/response by readOnly/writeOnly, select variants, drop what the     │
│                generator cannot express → one synthetic generator-input spec              │
│  3. tfplugingen-openapi → codegen/ir/vN.json (provider code spec)                         │
│  4. verify IR (no skipped resources, every attribute maps to a JSON field)                │
│  5. emit       internal/apiversion/generated/vN.json (runtime metadata)                   │
│                codegen/reports/vN.md (provenance and lossy normalizations)                │
└───────────────────────────────────────────────────────────────────────────────────────────┘
            │
            ▼  all versions
  6. merge IRs → codegen/ir/merged.json (superset schema)
  7. tfplugingen-framework → internal/generated/{resource,datasource}_*/  (schema functions)
  8. prune generated models, write internal/provider/zz_generated_schemas.go and
     templates/guides/api-versions.md.tmpl (support matrix)
            │
            ▼  runtime
internal/provider     generic resource/data source driven by schema + version metadata
internal/sdp          HTTP client: TLS, login/token, version negotiation, Accept header, errors
```

| Path | Kind | Purpose |
|---|---|---|
| `internal/apiversion/versions.json` | maintained | Version registry. The generator and the provider both read it. |
| `codegen/resources.yml` | maintained | Resource definitions, shared by all versions. |
| `internal/codegen`, `cmd/sdpgen` | maintained | Generation pipeline. |
| `internal/provider`, `internal/sdp` | maintained | Provider runtime. |
| `internal/provider/hooks.go`, `appliance_seed.go`, `controller_activation.go` | maintained | SDP behavior the specification cannot express. |
| `codegen/ir/*.json`, `codegen/reports/*.md` | generated | Reviewable generator output. |
| `internal/apiversion/generated/*.json` | generated | Runtime metadata for each version. |
| `internal/generated/**`, `internal/provider/zz_generated_*.go` | generated | Go schema code. |
| `templates/guides/api-versions.md.tmpl`, `docs/**` | generated | Documentation. |

Generated files are committed so that changes to the API surface show up in review. CI
regenerates everything and fails on any diff. Never edit generated files by hand.

## Multi-version design

### Constraint

Terraform asks for resource schemas (`GetProviderSchema`) before it configures the provider.
The schema therefore cannot depend on `api_version`.

### Chosen approach: superset schema plus version adapters in data

* **One Terraform schema per resource**, merged from the IR of every version:
  * an attribute is required only if every version requires it;
  * attributes present in only some versions are optional, and their description lists the
    supporting versions;
  * a default is kept only if every version has the attribute with the same default;
  * validators are kept only if all versions agree, so a stricter old enum never rejects a
    newer value;
  * if a type changed between versions, the newest type wins.
* **Per-version metadata** (`internal/apiversion/generated/vN.json`) describes each resource's
  endpoints, discriminator, fields, JSON names, read-only/write-only flags, required fields
  and enums for that version.
* **One generic implementation** (`internal/provider/resource.go`, `datasource.go`) converts
  between Terraform values and JSON using the metadata of the configured version. There are
  no version conditionals in resource code.
* **Plan-time validation** (`ModifyPlan` → `validate.go`) rejects:
  * resources the version does not have;
  * attributes the version does not have, has as read-only, or has with a different type;
  * missing attributes the version requires.

  This runs before any API call.
* **The version is selected once**, in provider `Configure`, by explicit `api_version`, legacy
  `client_version`, or negotiation through the 406 response of `POST /login`. It is bound to
  one client and one metadata set that every resource uses, so authentication, Accept headers
  (`application/vnd.appgate.peer-vN+json`), serialization and validation always agree.

### Alternatives considered

| Alternative | Why not |
|---|---|
| Separate resource types per version (`appgatesdp_v25_condition`) | Breaks every configuration on each SDP upgrade and multiplies the documentation. |
| A provider binary per version | Requirement: one release supports all versions. |
| Generated Go CRUD code per version and resource | Five copies of every resource, and it still needs a shared schema. The data-driven engine is one implementation tested against all versions. |
| Muxing the old SDK v2 provider with the new one | Mux requires identical provider schemas and disjoint resource names, so the old version logic would remain. |

### Unavoidable limitations

* Every version's users see the superset schema in documentation and editor completion.
  Attribute descriptions state which versions support each attribute.
* Defaults and enum validation of attributes that differ between versions move to the API,
  which reports a 422 error at apply time instead of plan time.
* An attribute whose type changed is only usable with the versions that have the newest type.

## Generation pipeline details

### Specification retrieval and pinning

`internal/codegen/fetch.go` uses the `git` CLI against `repository` from the registry:

1. `git ls-remote` confirms the branch `version-N` exists. A missing branch fails with the
   branch name.
2. It fetches the branch and checks that the pinned commit exists and is an ancestor of the
   branch head, so a commit from another version's branch is rejected.
3. It runs `git archive <commit>` into `.cache/specs/vN/<commit>/` and checks that
   `api_specs.yml` is OpenAPI 3 and declares `API version N`.
4. Extracted trees are cached by commit, so regeneration is offline and reproducible.

Generation never uses a branch head. `make pin [API_VERSION=vN]` moves the pins to the current
heads; the change to `versions.json` is reviewed like any other change. Each report in
`codegen/reports/vN.md` and the versions guide records the branch and commit used.

### Discovery (`internal/codegen/discover.go`)

Each version's specification decides what that version offers. Every path is classified, and
the classification is listed in `codegen/reports/vN.md`:

| Shape | Becomes |
|---|---|
| POST `/x` and GET `/x/{id}` | A collection resource. A discriminated request body gives one resource per variant, named `<variant>_<singular>` (e.g. `/secrets` → `credentials_secret`, `generic_secret`, `ssh_key_secret`). |
| GET and PUT `/x` | A singleton resource named after the path (`auto_update_settings`). |
| GET, POST and DELETE `/x`, without objects | A registration: an immutable singleton created with POST (`ztp_registration`). |
| Any other GET returning JSON | A read-only data source named after the path's literal segments (`/sites/{id}/resources` → `sites_resources`). Path and query parameters become attributes. |
| Everything else | An action: bulk upserts, reboots, upgrades, diagnostics, tests. Actions have no state and are not exposed. |

Every resource also gets a data source that looks objects up by `id` or exact `name`.

`codegen/resources.yml` only tunes discovery:
* explicit `resources` take precedence and keep established names and behavior (e.g. the
  typed policies, the adopted built-in identity providers);
* `names` renames discovered resources;
* `exclude` skips paths, with a reason.

New collections, variants and endpoints of a new API version therefore appear without any
change. Name collisions fail generation with instructions.

### Normalization (`internal/codegen/normalize.go`, `build.go`)

The SDP specification uses constructs `tfplugingen-openapi` v0.3.0 cannot handle. The
normalizer rewrites them:

| SDP construct | Generator limitation | Normalization |
|---|---|---|
| ~60 files; refs like `./common.yml#/definitions/X`, path items in other files | Cannot load multi-file references | Resolved by our own loader. Duplicate YAML keys (present in some branches) are tolerated, last one wins. |
| `allOf: [BaseEntity, {...}]` | Only one allOf member | Merged. Fragments such as description-only members are merged before validation. |
| `discriminator` / `oneOf` (identity providers, entitlements, client profiles) | Composition unsupported | Each Terraform resource selects one variant (`variant:` in `resources.yml`). Nested polymorphic fields become a union of all variants. |
| `readOnly` / `writeOnly` | Ignored | Read-only fields are removed from requests and write-only fields from responses. The generator then marks them Computed or input-only. |
| Secrets marked only `writeOnly` | Sensitivity comes only from `format: password` | Write-only and secret-named strings get `format: password`. |
| ECMA regex `pattern` | Generated code uses `regexp.MustCompile` | Patterns RE2 cannot compile are dropped and reported. |
| Recursive references | Infinite schema | Dropped at re-entry and reported. |
| Free-form objects (`additionalProperties: true`, no declared properties, untyped values) | Not representable | A JSON-encoded string attribute (`jsonencode`/`jsondecode`), compared by content. A body that is entirely free-form (`/my-preferences`) is exposed as one `value` attribute. |
| Object-valued query parameters (`filterBy`, `geolocationQuery`) | Not representable as parameters | Dropped from read-only data sources and reported. |
| `uniqueItems` primitive arrays | — | Become sets. |

Each resource gets synthetic paths (`/condition`, `/condition/{id}`) in the generator input, so
variants of one collection (for example the five policy types) map independently. The real
endpoints live in the metadata.

### Generator output post-processing

* **IR merge and overlay** (`merge.go`):
  * nested attributes have no defaults. The framework applies nested defaults
    wherever configuration is null, including elements the controller created.
    For example, a site's built-in DNS resolver uses `update_interval` 300 while
    the specification default is 60, which would be a permanent diff;
  * every attribute below a computed parent is computed. Terraform core treats
    a non-computed child value in prior state as proof the parent was
    configured, and plans to remove server-created elements. Required children
    are still enforced per version at plan time;
  * every resource gets a stable `id` with `UseStateForUnknown`;
  * user-chosen ids also get `RequiresReplaceIfConfigured`;
  * data sources take `id` or `name` as lookup arguments;
  * reserved Terraform names are rejected.
* **Pruning** (`prune.go`): `tfplugingen-framework` v0.4.1 always generates model and custom
  type code for nested attributes, named after the attribute alone. SDP reuses names such as
  `allow_sources` at different depths, so that code does not compile. The runtime works on
  `tftypes` values and never uses the models. An AST pass keeps only the generated schema
  functions and removes `CustomType` references. The schema itself is unchanged.

All output is formatted with `go/format`. The pipeline is deterministic: maps are sorted
before output, and the tests run it twice and compare the results.

### Commands

```bash
make generate                    # all versions
make generate API_VERSION=v25    # one version; other versions' files are not touched
make pin API_VERSION=v25         # move the pin to the head of version-25
make generate-check              # regenerate and fail on any diff (CI)
```

The pipeline needs Go and `git`. `tfplugingen-openapi`, `tfplugingen-framework` and
`tfplugindocs` are pinned with `tool` directives in `go.mod` and run via `go tool`.

## Runtime behavior

The generic engine (`convert.go`) applies SDP-specific reconciliation rules:

| Rule | Reason |
|---|---|
| Write-only fields keep the prior or planned value | The API never returns secrets. |
| After apply, fields the API does not echo keep the planned value | Terraform requires state to match the plan. |
| On refresh, fields that disappeared become null | Drift detection. |
| Optional, non-computed attributes stay null when the API returns a zero value | The API fills unset fields with zero values. |
| Lists of strings keep the configured order and case if the API returns the same elements | The API lowercases and sorts values such as tags. |
| Update is GET, then overlay the plan, then PUT the whole object | The SDP API requires the complete object on PUT; fields Terraform does not manage are kept. |
| After create and update, state is built from a fresh GET, not from the write response | Write responses are not canonical: timestamps have nanosecond rather than microsecond precision, and server-set fields (e.g. `blacklistedAt`) are missing. |
| 404 on read removes the resource from state | Out-of-band deletion leads to re-creation. |
| Discriminator mismatch on read removes the resource from state | The object's type changed outside Terraform. |

| JSON attributes keep their configured text if the API returns the same JSON | Key order and formatting are not changes. |
| Values outside the configured version's enum are rejected at plan time | The shared schema only validates values every version accepts. |

Resource kinds (discovered, or set in `codegen/resources.yml`):
* **collection**: POST, then GET/PUT/DELETE on `/{id}`;
* **singleton**: always exists. Create works like update: GET, overlay the plan, PUT, so
  settings the configuration does not mention keep their values. DELETE resets to defaults if
  supported;
* **registration**: POST to create, DELETE to remove, any change replaces it;
* **variant**: the discriminator is injected and checked;
* **adopt**: a built-in object is looked up by name; destroy only removes it from state;
* **non-collection APIs** (`create_method`, `read_list`, `id_property`, `delete_path`,
  `immutable`): for example licenses are created with POST `/license`, read from the
  `licenses` list of GET `/license` and deleted with DELETE `/license/{id}`; blacklist entries
  are identified by `userDistinguishedName`. Immutable objects have no update, so every input
  change is planned as a replacement. Schemas stay generated per version.

Data source kinds:
* **lookup**: by `id`, or by exact `name` within the resource's discriminator;
* **singleton**: GET of the settings;
* **read-only**: GET of an endpoint with path parameters substituted and query parameters
  added from the configuration.

Hand-written code is limited to:
* `hooks.go`: an activated appliance is deactivated before delete;
* `appliance_seed.go`: the seed export action;
* `controller_activation.go`: updates `controller`/`admin_interface` of an appliance and
  waits for `/appliances/status` to report `controller_ready` or `appliance_ready`. Its
  `controller` and `admin_interface` attributes are taken from the generated appliance
  schema, and conversion and version validation use the generated appliance metadata.

## Adding a new SDP / API version (example: SDP 7.1, API v26)

1. Make sure the branch `version-26` exists in `appgate/sdp-api-specification`.
2. Register the version in `internal/apiversion/versions.json`:
   ```json
   {"api_version": "v26", "sdp_version": "7.1", "branch": "version-26", "commit": "0000000000000000000000000000000000000000"}
   ```
3. Pin it and generate it. Existing versions are not modified:
   ```bash
   make pin API_VERSION=v26
   make generate API_VERSION=v26
   ```
4. Review the diff:
   * `codegen/reports/v26.md` (path classification, unsupported resources, normalization
     warnings);
   * `codegen/ir/v26.json`;
   * `templates/guides/api-versions.md.tmpl` (new and removed attributes);
   * the merged schema.
5. Run `make docs check`. Every table-driven test runs against v26 automatically.
6. Optionally, run acceptance tests against a v26 controller with `SDP_ACC_V26_URL`.
7. Release.

New collections, variants, singletons and read-only endpoints are discovered. Only add
`names` or `exclude` entries to `codegen/resources.yml` where the defaults do not fit. No Go
code changes are needed for new versions, attributes, resources or data sources.

## Testing

| Layer | Location | Live SDP needed |
|---|---|---|
| Registry mappings and validation | `internal/apiversion` | no |
| Fetching (branches, pins, wrong commits, wrong versions), normalization, merge, pruning, determinism, single-version isolation, full pipeline on fixture branches | `internal/codegen` | no |
| Committed artifacts match the registry, and the cached spec when present | `internal/codegen/artifacts_test.go` | no |
| HTTP client: login, token cache, re-login, retries, TLS, errors, negotiation | `internal/sdp` | no |
| Conversion and validation units | `internal/provider/*_test.go` | no |
| Terraform lifecycle (plan/apply/import/refresh/drift/destroy) with the Terraform CLI against the in-memory controller `internal/sdptest`, table-driven over all versions; every resource of every version round-trips | `internal/provider` | no |
| Acceptance: schema conformance (read-only), lifecycle, references, a resource matrix (create, clean re-plan, import, destroy), built-in data sources, replication source status and token rejection, sweep | `internal/acctest` | yes (gated) |

`internal/sdptest` is version-aware. It serves the collections of the version's metadata,
checks the Accept header, rejects unknown fields and missing required fields like a strict
controller, never returns write-only fields, and can inject failures. The acceptance suite is
also run against it (`TestAcceptanceSuiteAgainstMock`), so the suite itself is tested.

Run `make testacc` against a controller of every registered version before a
release; adding a version only needs its `SDP_ACC_V<N>_*` variables. Features
the collective is not licensed for (HTTP 402) are skipped with the reason, not
failed.

## Known limitations and open items

* Action endpoints (bulk upserts, reboots, upgrades, diagnostics, connection tests, CA
  switch) are intentionally not exposed: they have no state to manage.
* Endpoints with non-JSON responses (appliance metrics, recorded session downloads, session
  observation) are not exposed.

* Attributes that come from additional endpoints are not available
  (`client_profile.url`, `replication_target.registration_token`, `certificate_authority.pem`).
* Licenses are not covered by destructive acceptance tests, to avoid replacing the license of a
  test collective. A real replication source registration is not covered either, because it
  turns the test collective into a replica of another one.
* Controllers return some fields their specification branch does not list. For example, SDP
  6.6 returns `clientInterface.stunPort`, `connectionBroker`, `geoIpSettings` and `vpn.quic`,
  which `version-23` only adds in v24. The provider follows the specification, so these are
  rejected at plan time for v23. `TestAccSchemaConformance` lists such gaps.
* The generator output still needs validation against real controllers. Attributes that the
  API fills with server defaults but the specification does not mark as defaults may show
  diffs. `TestAccSchemaConformance` and the acceptance lifecycle tests are designed to find
  these.
* Lists of objects (such as entitlement `actions`) are ordered. If the controller reorders
  them, plans will show a diff.
* `tfplugingen-openapi` and `tfplugingen-framework` are tech-preview tools. Upgrading them
  should be done with `make generate`, followed by review of the IR and schema diff.
