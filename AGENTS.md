# AGENTS.md - Dell Terraform Provider for PowerMax

## Project Overview

This is the Terraform provider for Dell PowerMax enterprise storage arrays. It implements resources and data sources using HashiCorp's Terraform Plugin Framework, enabling infrastructure-as-code management of PowerMax arrays.

- **Language:** Go 1.25
- **Module path:** `terraform-provider-powermax`
- **Terraform Plugin Framework:** v1.19.0
- **SDK:** `dell/powermax-go-client` (vendored, local)
- **Registry address:** `registry.terraform.io/dell/powermax`
- **License:** Mozilla Public License 2.0

## Architecture

The provider follows the standard Terraform Plugin Framework architecture. It runs as a gRPC server that Terraform Core communicates with to manage PowerMax resources.

### Provider Configuration

The provider authenticates to a PowerMax Unisphere server using endpoint, username, password, serial number, and PMax version. Configuration can be supplied via HCL provider block or environment variables (`POWERMAX_ENDPOINT`, `POWERMAX_USERNAME`, `POWERMAX_PASSWORD`, `POWERMAX_INSECURE`, `POWERMAX_TIMEOUT`).

### SDK Strategy

Uses a **vendored SDK** — `dell/powermax-go-client` lives inside the provider repo as a local directory. The `go.mod` declares:

```go
require dell/powermax-go-client v0.0.0
replace dell/powermax-go-client => ./powermax-go-client-100
```

SDK and provider release together. Changes to the SDK require changes in the same repo.

### Resources and Data Sources

The provider exposes approximately 8 resources and 1 data source covering PowerMax entities such as storage groups, port groups, hosts, host groups, masking views, and volumes.

## Directory Structure

```
main.go                           Entry point (providerserver.Serve)
powermax/
  provider/
    provider.go                   Provider configuration, resource/datasource registration
    *_resource.go                 Resource implementations
    *_resource_schema.go          Resource schema definitions
    *_datasource.go               Data source implementations
    *_test.go                     Unit and acceptance tests
  helper/                         Shared helper functions
  models/                         Terraform state model structs
  constants/                      Shared constants
client/                           PowerMax client wrapper
goClientZip/                      SDK distribution archives
powermax-go-client-100/           Vendored PowerMax Go SDK (local dependency)
examples/                         Example HCL configurations
docs/                             Generated documentation
templates/                        Documentation templates
tools/                            Build and generation tools
about/                            Provider metadata
```

## Build Commands

| Command | Description |
|---------|-------------|
| `make build` | Compile the provider binary |
| `make install` | Build and install to `~/.terraform.d/plugins/` |
| `make test` | Run unit tests |
| `make testacc` | Run acceptance tests (`TF_ACC=1`, requires live hardware) |
| `make check` | Run `gofmt`, `golangci-lint`, `go vet` |
| `make gosec` | Run security scan with `gosec` |
| `make cover` | Generate HTML coverage report |
| `make generate` | Run `go generate` (docs generation) |

## Testing

### Unit Tests (mockey)

- Test files follow `*_test.go` convention in `powermax/provider/`.
- Frameworks: `github.com/stretchr/testify` (assertions), `github.com/bytedance/mockey` (function-level mocking).
- Run with `make test`.
- No hardware required.

### Acceptance Tests (terraform-plugin-testing)

- **Requires live PowerMax hardware** with credentials set via environment variables.
- Creates real resources — clean up after failures.
- Run with `make testacc`.

### Running Tests

```bash
# Unit tests (no hardware)
make test

# Acceptance tests (requires live hardware)
export POWERMAX_ENDPOINT="https://unisphere-ip"
export POWERMAX_USERNAME="admin"
export POWERMAX_PASSWORD="secret"
export POWERMAX_INSECURE="true"
make testacc
```

## Code Style and Conventions

### Code Organization Patterns

- **Resource pattern:** Each resource has up to three files: `<name>_resource.go`, `<name>_resource_schema.go`, plus helpers.
- **Models:** Terraform state structs in `powermax/models/` using `tfsdk` struct tags.
- **Helpers:** API-to-Terraform mapping functions in `powermax/helper/`.

### File Header

All source files must include the Dell copyright and MPL 2.0 license header.

## Common Development Tasks

### Adding a New Resource

1. Create resource, schema, and model files following existing patterns.
2. Add helper functions for API-to-Terraform mapping.
3. Register in `powermax/provider/provider.go`.
4. Add unit and acceptance tests.
5. Create example HCL in `examples/resources/powermax_<name>/`.
6. Run `make generate` to produce documentation.

### Updating the Vendored SDK

Edit files directly in `./powermax-go-client-100/`. No `go mod tidy` needed (local replace directive). Commit SDK and provider changes together.

## CI/CD

GitHub Actions workflows in `.github/workflows/`. GoReleaser configuration in `.goreleaser.yaml` builds cross-platform binaries.

## Code Ownership

All files are owned by the maintainers defined in `.github/CODEOWNERS`.
