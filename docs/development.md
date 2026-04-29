# Development

## Building

The following will build the binary for the current platform and place it in the `releases` directory.

```console
goreleaser build --clean --snapshot --single-target
```

## Resource Code Generation Tool

The `codegen` tool at `tools/codegen/` generates all boilerplate Go files for a new aws-nuke resource from a YAML definition. Instead of manually writing resource files, client interfaces, mocks, and tests, you define the resource declaratively and let the tool produce everything.

### DSL File Location

YAML definition files are stored directly in `resources/` alongside the generated `.go` files, using the naming convention `resources/<service>-<resource>.yaml`. This keeps the definition co-located with the code it generates — easy to find, easy to regenerate. Go's build toolchain ignores `.yaml` files, so they don't interfere with compilation.

For example, a resource at `resources/mailmanager-relay.go` has its definition at `resources/mailmanager-relay.yaml`.

### Quick Start

1. Create a YAML DSL file in `resources/` describing your resource (see [DSL Reference](codegen.md) for the full schema):

    ```yaml
    service: mailmanager
    resource: relay
    resourceName: MailManagerRelay
    scope: account
    sdkPackage: mailmanager

    list:
      strategy: flat
      operation: ListRelays
      pagination: paginator
      itemsField: Relays

    delete:
      operation: DeleteRelay
      inputFields:
        - RelayId

    fields:
      - name: RelayId
        type: "*string"
        fromList: RelayId
      - name: RelayName
        type: "*string"
        fromList: RelayName

    stringRepresentation:
      field: RelayName
    ```

2. Validate the DSL without generating files:

    ```console
    go run ./tools/codegen/ --validate resources/mailmanager-relay.yaml
    ```

3. Preview what would be generated (dry run):

    ```console
    go run ./tools/codegen/ --dry-run resources/mailmanager-relay.yaml
    ```

4. Generate the files (overwrites existing, runs post-checks, and commits by default):

    ```console
    go run ./tools/codegen/ resources/mailmanager-relay.yaml
    ```

5. Generate without committing:

    ```console
    go run ./tools/codegen/ --no-commit resources/mailmanager-relay.yaml
    ```

6. Generate without build checks:

    ```console
    go run ./tools/codegen/ --no-build resources/mailmanager-relay.yaml
    ```

7. Generate without overwriting existing files:

    ```console
    go run ./tools/codegen/ --no-force resources/mailmanager-relay.yaml
    ```

8. Regenerate mock tests even if they already exist:

    ```console
    go run ./tools/codegen/ --force-mock-tests resources/mailmanager-relay.yaml
    ```

9. Regenerate integration tests even if they already exist:

    ```console
    go run ./tools/codegen/ --force-integration-tests resources/mailmanager-relay.yaml
    ```

### Generated Files

For a resource using an interface-based client, the tool produces up to 5 files:

| File | Description |
|---|---|
| `resources/<service>-<resource>.go` | Resource struct, lister, Remove, Properties, String, Filter, Settings |
| `resources/<service>.go` | Client interface (created or appended to) |
| `resources/<service>_mock_test.go` | Testify mock struct (created or appended to) |
| `resources/<service>-<resource>_mock_test.go` | Mock-based unit tests (List, Remove, Properties, String) |
| `resources/<service>-<resource>_test.go` | Integration test scaffold (only if `integrationTest` is defined) |

For concrete client types (`svcType: concrete`), only the resource file is generated.

Mock test files and integration test files are only generated when they do not already exist. This protects hand-edited tests from being overwritten during regeneration. Use `--force-mock-tests` or `--force-integration-tests` to regenerate them.

### CLI Flags

| Flag | Description |
|---|---|
| `--dry-run` | Print generated content to stdout without writing files |
| `--no-force` | Do not overwrite existing resource files (default: overwrite) |
| `--validate` | Check DSL for errors only, no generation |
| `--output-manifest` | Write a JSON manifest of generated file paths and statuses |
| `--no-build` | Skip running `goimports`, `go build`, `golangci-lint` after generation |
| `--no-commit` | Skip staging and committing generated files |
| `--force-mock-tests` | Regenerate mock test files even if they already exist |
| `--force-integration-tests` | Regenerate integration test files even if they already exist |
| `--output-dir` | Override output directory (default: current directory) |

### Running Tests

```console
go test -v ./tools/codegen/
```

### Sample DSL Files

The `tools/codegen/testdata/` directory contains sample YAML files for every supported pattern:

- `flat-simple.yaml` — minimal flat list with paginator
- `flat-with-describe.yaml` — flat list with describe enrichment
- `flat-with-tags.yaml` — flat list with tag fetching
- `nested-two-level.yaml` — two-level nested list
- `nested-three-level.yaml` — three-level nested list
- `singleton.yaml` — singleton resource
- `with-filters.yaml` — status-based filters
- `with-settings.yaml` — deletion protection settings
- `with-pre-deletion.yaml` — pre-deletion cleanup steps
- `with-overrides.yaml` — raw Go override blocks
- `with-integration-test.yaml` — integration test configuration

Use these as starting points for new resources. See the full [DSL Reference](codegen.md) for all available fields.

## Documentation

This is built using Material for MkDocs and can be run very easily locally providing you have docker available.

[Read more about it here.](documentation.md)

### Running Locally

```console
make docs-serve
```
