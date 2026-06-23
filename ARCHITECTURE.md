# Architecture: terraform-provider-powermax

## Metadata

<!-- yaml-metadata-start -->
scope_paths: ["./"]
capture_git_sha: "4ee9e5f81d27df1681427850ddd3f948837e69b7"
status: "current"
auto_update: false
preview_before_apply: true
scaffold_version: "1.0"
<!-- yaml-metadata-end -->

---

## Purpose and Structure

Terraform provider for Dell PowerMax enterprise storage arrays.
Implements 8 managed resources and 9 data sources
using HashiCorp's Terraform Plugin Framework, enabling
infrastructure-as-code management via REST API.

The provider is a standalone Go binary that communicates with Terraform
Core over gRPC (go-plugin protocol).

**SDK strategy:** Vendored SDK — lives inside the provider repo as `powermax-go-client-100/`. `go.mod` declares `require dell/powermax-go-client v0.0.0` with `replace => ./powermax-go-client-100`. SDK and provider release together.

---

## Components

| Component | Path | Responsibility |
|-----------|------|---------------|
| Entry point | `main.go` | `providerserver.Serve` — starts gRPC server |
| Provider | `powermax/provider/provider.go` | Schema, Configure, resource/datasource registration |
| Resources | `powermax/provider/*_resource.go` | CRUD lifecycle for 8 managed resources |
| Data sources | `powermax/provider/*_datasource.go` | Read-only queries for 9 data sources |
| Vendored SDK | `powermax-go-client-100/` | Local PowerMax Go SDK |
| SDK archives | `goClientZip/` | SDK distribution archives |
| Client wrapper | `client/` | Wraps vendored SDK |
| Models | `powermax/models/` | Terraform state model structs |
| Helper | `powermax/helper/` | Type mapping functions |
| Examples | `examples/` | HCL configurations for resources and data sources |
| Docs | `docs/` | Generated provider documentation |

---

## Key Behaviors

### Authentication

**GIVEN** a user configures the provider with endpoint, username,
and password (via HCL block or environment variables)
**WHEN** `Configure()` runs
**THEN** (1) env vars `POWERMAX_ENDPOINT`, `POWERMAX_USERNAME`,
`POWERMAX_PASSWORD`, `POWERMAX_INSECURE`, `POWERMAX_TIMEOUT`
override HCL values, (2) SDK client is initialized, (3) authentication
is validated before any resource operations proceed

### Resource CRUD Lifecycle

**GIVEN** a resource definition in HCL
**WHEN** `terraform apply` runs
**THEN** the resource's `Create()` reads the plan into a model struct,
calls the SDK/client to create the resource, maps the API response
back to Terraform state, and sets `resp.State`

### Drift Detection

**GIVEN** a resource exists in Terraform state
**WHEN** `terraform plan` or `terraform refresh` runs
**THEN** `Read()` calls the SDK/client to fetch current state,
compares it with stored state, and updates the state if drifted

### Import

**GIVEN** a resource exists on the hardware but not in Terraform state
**WHEN** `terraform import` runs
**THEN** `ImportState()` fetches the resource by ID and populates state

---

## Interfaces

### Provider Configuration Schema

| Attribute | Type | Env Var | Description |
|-----------|------|---------|-------------|
| `endpoint` | string | `POWERMAX_ENDPOINT` | Unisphere management IP or FQDN |
| `username` | string | `POWERMAX_USERNAME` | API username |
| `password` | string (sensitive) | `POWERMAX_PASSWORD` | API password |
| `insecure` | bool | `POWERMAX_INSECURE` | Skip TLS verification (lab only) |
| `timeout` | int64 | `POWERMAX_TIMEOUT` | Request timeout in seconds |

---

## Dependencies

| Depends On | For |
|------------|-----|
| `dell/powermax-go-client` (vendored, local) | Platform API SDK/client |
| `hashicorp/terraform-plugin-framework` v1.19.0 | Core provider interfaces |
| `hashicorp/terraform-plugin-framework-validators` | Attribute validation |
| `hashicorp/terraform-plugin-log` | Structured logging |
| `hashicorp/terraform-plugin-testing` | Acceptance test harness |
| `bytedance/mockey` | Unit test function-level mocking |
| `stretchr/testify` | Test assertions |

---

## Known Constraints

1. **Terraform Plugin Framework only** — no SDK v2 code.
2. **CGO_ENABLED=0** — static binaries for all platforms.
3. **Sensitive attributes marked** — credentials never in plan output.
4. **ImportState required** — all resources support `terraform import`.
5. **Environment variable fallback** — all credentials support env vars.
6. **Acceptance tests gated** — never run without `TF_ACC=1`.
7. **Endpoint format** — Unisphere management IP or FQDN.
8. **Vendored SDK co-versioned** — SDK changes require
   commits in the same repo.

---

## Change History

| Date | Feature | What Changed | Author |
|------|---------|-------------|--------|
| 2026-06-10 | Initial architecture | Provider-specific architecture extracted from generic multi-provider doc | architecture-agent |
