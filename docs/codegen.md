# Resource Code Generation Tool — DSL Reference

The `codegen` tool reads a YAML definition file and generates all boilerplate Go source files for an aws-nuke resource. This page documents every field in the YAML DSL.

DSL files are stored directly in `resources/` alongside the generated `.go` files, named `resources/<service>-<resource>.yaml`. Go ignores `.yaml` files during compilation, so they don't interfere — but they're right next to the code they produce, making it obvious which definition generated which file. To regenerate:

```console
go run ./tools/codegen/ resources/<service>-<resource>.yaml
```

---

## Top-Level Fields

| Field | Type | Required | Description |
|---|---|---|---|
| `service` | string | yes | AWS service name used in file naming. Becomes the prefix of generated filenames (e.g. `mailmanager` → `resources/mailmanager-relay.go`). |
| `resource` | string | yes | Resource name used in file naming. Becomes the suffix (e.g. `relay` → `resources/mailmanager-relay.go`). |
| `resourceName` | string | yes | Go constant name used for the resource registration and struct name (e.g. `MailManagerRelay`). Generates `const MailManagerRelayResource = "MailManagerRelay"`. |
| `scope` | string | yes | Either `account` or `region`. Maps to `nuke.Account` or `nuke.Region` in the registration. |
| `sdkPackage` | string | yes | AWS SDK v2 Go package name (e.g. `mailmanager`, `eks`, `mediapackagev2`). Used in imports and client construction. |
| `svcType` | string | no | Set to `concrete` to use `*<sdkPackage>.Client` directly as the `svc` field type. When omitted or empty, the tool generates a client interface (e.g. `MailmanagerClient`) and uses that as the `svc` type. Concrete types skip interface, mock, and mock test file generation. |

---

## List Strategies

The `list` section defines how resources are enumerated. Exactly one of three strategies must be specified. If none of the declarative strategies fit your use case, use `list.override` to provide raw Go code for the entire `List` method body.

### List Fields (all strategies)

| Field | Type | Required | Description |
|---|---|---|---|
| `strategy` | enum | yes | One of: `flat`, `nested`, `singleton`. |

### Flat List

A single API call, optionally paginated. For complex listing logic that the declarative fields can't express, use `list.override` to provide raw Go code for the entire `List` method body.

| Field | Type | Required | Description |
|---|---|---|---|
| `operation` | string | yes (flat) | SDK operation name (e.g. `ListRelays`). |
| `pagination` | enum | yes (flat) | Controls how pagination is handled. One of: `paginator` — uses the SDK's built-in paginator (`<sdkPackage>.New<Operation>Paginator`). `nextToken` — generates a manual for-loop that checks `resp.NextToken` each iteration. `none` — single API call with no pagination. |
| `itemsField` | string | yes (flat) | Field name on the response struct containing the slice of items (e.g. `Relays`, `Clusters`). |
| `describe` | object | no | Optional describe call to enrich each listed item. See [Describe](#describe). |
| `tags` | object | no | Optional tag-fetching call per item. See [Tags](#tags). |
| `override` | string | no | Raw Go code that replaces the entire `List` method body. When set, the declarative list fields are ignored. |

#### Pagination Modes

`paginator` generates:
```go
paginator := mailmanager.NewListRelaysPaginator(svc, &mailmanager.ListRelaysInput{})
for paginator.HasMorePages() {
    resp, err := paginator.NextPage(ctx)
    // ...
}
```

`nextToken` generates:
```go
params := &mailmanager.ListRelaysInput{}
for {
    resp, err := svc.ListRelays(ctx, params)
    // ... process items ...
    if resp.NextToken == nil { break }
    params.NextToken = resp.NextToken
}
```

`none` generates:
```go
resp, err := svc.ListRelays(ctx, &mailmanager.ListRelaysInput{})
```

#### Describe

Enrich each listed item with an additional API call.

| Field | Type | Required | Description |
|---|---|---|---|
| `operation` | string | yes | SDK describe operation (e.g. `DescribeCluster`). |
| `inputMapping` | map[string]string | yes | Maps describe input parameter names to source fields on the list item. Key = input param, value = field on the iterator variable. |
| `responseField` | string | no | If the describe response nests the data under a field (e.g. `Cluster`), specify it here. Fields with `fromDescribe` will read from `dcResp.Cluster.<field>` instead of `dcResp.<field>`. |

```yaml
describe:
  operation: DescribeCluster
  inputMapping:
    Name: cluster          # passes item.cluster as the Name param
  responseField: Cluster   # reads fields from dcResp.Cluster
```

#### Tags

Fetch tags per item using a separate API call.

| Field | Type | Required | Description |
|---|---|---|---|
| `operation` | string | yes | SDK tag-fetching operation (e.g. `ListTagsForResource`). |
| `arnField` | string | yes | Field on the list item containing the ARN to pass to the tag call. |

```yaml
tags:
  operation: ListTagsForResource
  arnField: Arn
```

Generated code logs a warning (instead of failing) if the tag call errors, so a single untaggable resource doesn't break the entire list.

### Nested List

Multiple levels of API calls where child resources are listed per parent. Requires at least 2 levels.

#### Level Fields

| Field | Type | Required | Description |
|---|---|---|---|
| `operation` | string | yes | SDK list operation for this level. |
| `pagination` | enum | yes | One of: `paginator`, `nextToken`, `none`. Same behavior as flat list pagination. |
| `itemsField` | string | yes | Response field containing the slice of items at this level. |
| `iteratorVar` | string | yes | Variable name for the current item at this level. Used in `fromList` field references (e.g. `cg.ChannelGroupName`) and in child levels' `parentInputMapping`. |
| `parentInputMapping` | map[string]string | no | Maps input parameter names to parent-level fields. Key = input param, value = `<parentIteratorVar>.<field>`. Not needed on the first level. |

```yaml
list:
  strategy: nested
  levels:
    - operation: ListChannelGroups
      pagination: paginator
      itemsField: Items
      iteratorVar: cg
    - operation: ListChannels
      pagination: paginator
      itemsField: Items
      iteratorVar: ch
      parentInputMapping:
        ChannelGroupName: cg.ChannelGroupName
```

The tool generates nested for-loops, one per level. Each level uses its own pagination strategy independently.

### Singleton

A single Get/Describe call returning at most one resource.

| Field | Type | Required | Description |
|---|---|---|---|
| `operation` | string | yes (singleton) | SDK get/describe operation (e.g. `GetBlockPublicAccessConfiguration`). |
| `responseField` | string | no | If the response nests the data under a field, specify it here. Fields with `fromList` will read from `resp.<responseField>.<field>`. |
| `nilCheck` | bool | no | When `true` and `responseField` is set, generates a nil check that returns an empty list if the response field is nil. |

```yaml
list:
  strategy: singleton
  operation: GetBlockPublicAccessConfiguration
  responseField: BlockPublicAccessConfiguration
  nilCheck: true
```

Generated code:
```go
resp, err := svc.GetBlockPublicAccessConfiguration(ctx, &emr.GetBlockPublicAccessConfigurationInput{})
if err != nil { return nil, err }
if resp.BlockPublicAccessConfiguration == nil { return nil, nil }
resources = append(resources, &EMRBlockPublicAccessConfiguration{
    svc:       svc,
    AccountID: resp.BlockPublicAccessConfiguration.AccountID,
})
```

---

## Delete

Defines how the resource is removed.

| Field | Type | Required | Description |
|---|---|---|---|
| `operation` | string | yes | SDK delete operation name (e.g. `DeleteRelay`). |
| `inputFields` | []string | yes | List of resource struct field names passed to the delete input. Each generates `<FieldName>: r.<FieldName>` in the input struct. |
| `override` | string | no | Raw Go code that replaces the entire generated `Remove` method body. When set, `inputFields` is still required for validation but the generated delete call is replaced by this code. |

```yaml
delete:
  operation: DeleteRelay
  inputFields:
    - RelayId
```

Generated code:
```go
func (r *MailManagerRelay) Remove(ctx context.Context) error {
    _, err := r.svc.DeleteRelay(ctx, &mailmanager.DeleteRelayInput{
        RelayId: r.RelayId,
    })
    return err
}
```

With `delete.override`:
```yaml
delete:
  operation: PutBlockPublicAccessConfiguration
  inputFields:
    - BlockPublicAccessConfiguration
  override: |
    _, err := r.svc.PutBlockPublicAccessConfiguration(ctx, &emr.PutBlockPublicAccessConfigurationInput{
        BlockPublicAccessConfiguration: &emrtypes.BlockPublicAccessConfiguration{
            BlockPublicSecurityGroupRules: aws.Bool(false),
        },
    })
    return err
```

Note: When `override` is set, the generated settings-based deletion protection and pre-deletion steps are also skipped — the override is the complete method body.

---

## Fields

Each entry becomes a struct field on the generated resource. At least one field is required.

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | Go field name. Use PascalCase for exported fields, camelCase for unexported. |
| `type` | string | yes | Go type string (e.g. `*string`, `*time.Time`, `map[string]string`, `*bool`, `apigatewayv2types.PublishStatus`). |
| `fromList` | string | no | Source field on the list response item. For flat lists, this is a field on the iterator variable (e.g. `RelayId`). For nested lists, prefix with the iterator var (e.g. `cg.ChannelGroupName`). |
| `fromDescribe` | string | no | Source field on the describe response. Read from `dcResp.<field>` or `dcResp.<responseField>.<field>` if `describe.responseField` is set. |
| `fromTags` | bool | no | When `true`, this field is populated from the tag-fetching call result. Type should be `map[string]string`. |
| `exported` | *bool | no | Defaults to `true` (nil = exported). Set to `false` to generate a lowercase field name excluded from `Properties()` output. Used for internal fields like `protection`. |
| `propertyTag` | string | no | Custom `property:"..."` struct tag. Use `property:"-"` to exclude from properties, or `property:"name=CustomName"` for a custom property name. |

A field should have exactly one of `fromList`, `fromDescribe`, or `fromTags` to indicate its data source. Fields without any source (like `svc`) are not populated during listing.

```yaml
fields:
  - name: Name
    type: "*string"
    fromList: cluster
  - name: CreatedAt
    type: "*time.Time"
    fromDescribe: CreatedAt
  - name: Tags
    type: "map[string]string"
    fromTags: true
  - name: protection
    type: "*bool"
    fromDescribe: DeletionProtection
    exported: false
  - name: PublishStatus
    type: "apigatewayv2types.PublishStatus"
    fromList: PublishStatus
    propertyTag: "name=PublishStatus"
```

---

## Filters

When filters are defined, the tool generates a `Filter()` method on the resource. If no filters are defined, the method is omitted entirely. Each filter rule becomes a condition that returns an error (skipping the resource) when matched.

### Filter Fields

| Field | Type | Required | Description |
|---|---|---|---|
| `field` | string | yes | Name of the struct field to check (e.g. `Status`, `IsDefault`). |
| `operator` | enum | yes | One of: `equals`, `boolFalse`, `boolTrue`. |
| `values` | []string | for `equals` | List of values that trigger the filter. |
| `message` | string | yes | Error message returned when the filter matches. |

### Operators

#### `equals`

Checks if a field matches any of the listed values. Generates one `if` branch per value. The field is cast to `string()` to handle SDK enum types.

```yaml
filters:
  - field: Status
    operator: equals
    values:
      - DeleteInProgress
      - Cancelled
      - ReleaseInProgress
    message: "already being deleted or cancelled"
```

Generated code:

```go
func (r *MyResource) Filter() error {
    if string(r.Status) == "DeleteInProgress" {
        return fmt.Errorf("already being deleted or cancelled")
    }
    if string(r.Status) == "Cancelled" {
        return fmt.Errorf("already being deleted or cancelled")
    }
    if string(r.Status) == "ReleaseInProgress" {
        return fmt.Errorf("already being deleted or cancelled")
    }
    return nil
}
```

#### `boolFalse`

Checks if a `*bool` field is non-nil and `false`. Useful for filtering resources where a flag being false means they shouldn't be deleted.

```yaml
filters:
  - field: Enabled
    operator: boolFalse
    message: "resource is not enabled"
```

Generated code:

```go
if r.Enabled != nil && !*r.Enabled {
    return fmt.Errorf("resource is not enabled")
}
```

#### `boolTrue`

Checks if a `*bool` field is non-nil and `true`. Useful for filtering AWS-managed or default resources.

```yaml
filters:
  - field: IsDefault
    operator: boolTrue
    message: "cannot delete default resource"
```

Generated code:

```go
if r.IsDefault != nil && *r.IsDefault {
    return fmt.Errorf("cannot delete default resource")
}
```

### Multiple Filters

Multiple filter rules are checked in order; the first match returns the error.

```yaml
filters:
  - field: Status
    operator: equals
    values: [DeleteInProgress]
    message: "already being deleted"
  - field: IsDefault
    operator: boolTrue
    message: "cannot delete default resource"
```

### Override

If the declarative operators aren't sufficient, use `filterOverride` (a top-level field) to provide raw Go code. When set, the declarative `filters` list is ignored.

```yaml
filterOverride: |
  if strings.HasPrefix(*r.Path, "/aws-service-role/") {
      return fmt.Errorf("cannot delete service-linked roles")
  }
  return nil
```

---

## Settings

When settings are defined, the tool generates a `Settings(*libsettings.Setting)` method on the resource, adds `settings *libsettings.Setting` to the struct, registers the setting names in `registry.Registration.Settings`, and (for deletion protection settings) generates conditional disable logic in the `Remove` method.

### Setting Fields

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | Setting name as it appears in the user's nuke config (e.g. `DisableDeletionProtection`). Registered in `registry.Registration.Settings`. |
| `protectionField` | string | no | Name of a `*bool` field on the resource struct that indicates whether protection is enabled. When set, the Remove method checks this field. |
| `disableOperation` | string | no | SDK operation to call to disable protection (e.g. `UpdateClusterConfig`). Required if `protectionField` is set. |
| `disableInput` | map[string]any | no | Input parameters for the disable operation. String values are treated as resource field references (`r.<value>`). Boolean values become `aws.Bool(<value>)`. |

```yaml
settings:
  - name: DisableDeletionProtection
    protectionField: protection
    disableOperation: UpdateClusterConfig
    disableInput:
      Name: Name                    # becomes: Name: r.Name
      DeletionProtection: false     # becomes: DeletionProtection: aws.Bool(false)
```

Generated code in `Remove`:
```go
if ptr.ToBool(r.protection) && r.settings.GetBool("DisableDeletionProtection") {
    _, err := r.svc.UpdateClusterConfig(ctx, &eks.UpdateClusterConfigInput{
        Name:               r.Name,
        DeletionProtection: aws.Bool(false),
    })
    if err != nil { return err }
}
```

Generated registration:
```go
registry.Register(&registry.Registration{
    // ...
    Settings: []string{"DisableDeletionProtection"},
})
```

A setting without `protectionField` is still registered and generates the `Settings()` method, but doesn't add any conditional logic to `Remove`.

---

## Dependencies

A list of resource constant names that must be deleted before this resource. Generates a `DependsOn` field in the registration.

| Field | Type | Required | Description |
|---|---|---|---|
| `dependencies` | []string | no | List of resource constant names (e.g. `MailManagerRuleSet`). Each is referenced as `<Name>Resource` in the generated code. |

```yaml
dependencies:
  - MailManagerRuleSet
  - MailManagerEndpoint
```

Generated code:
```go
registry.Register(&registry.Registration{
    // ...
    DependsOn: []string{
        MailManagerRuleSetResource,
        MailManagerEndpointResource,
    },
})
```

---

## Pre-Deletion Steps

API calls that execute before the main delete call in the `Remove` method. Two types are supported.

### `listAndBatchDelete`

Lists dependent items and batch-removes them before the main delete.

| Field | Type | Required | Description |
|---|---|---|---|
| `type` | enum | yes | Must be `listAndBatchDelete`. |
| `listOperation` | string | yes | SDK operation to list dependent items (e.g. `ListTargets`). |
| `listInput` | map[string]string | yes | Input params for the list call. Key = param name, value = resource struct field name. |
| `listItemsField` | string | yes | Response field containing the slice of dependent items. |
| `deleteOperation` | string | yes | SDK operation to batch-delete the items (e.g. `DeregisterTargets`). |
| `deleteInput` | map[string]string | yes | Input params for the delete call (excluding the items slice). Key = param name, value = resource struct field name. |
| `deleteItemsField` | string | yes | Input field name for the items slice in the delete call (e.g. `Targets`). The singular form is used as the SDK types struct name. |
| `itemMapping` | map[string]string | yes | Maps fields from each listed item to the delete item struct. Key = target field, value = source field on the listed item. |

```yaml
preDeletion:
  - type: listAndBatchDelete
    listOperation: ListTargets
    listInput:
      TargetGroupIdentifier: ARN
    listItemsField: Items
    deleteOperation: DeregisterTargets
    deleteInput:
      TargetGroupIdentifier: ARN
    deleteItemsField: Targets
    itemMapping:
      Id: Id
      Port: Port
```

Generated code:
```go
listResp, err := r.svc.ListTargets(ctx, &vpclattice.ListTargetsInput{
    TargetGroupIdentifier: r.ARN,
})
if err == nil && len(listResp.Items) > 0 {
    var items []vpclatticetypes.Target
    for _, t := range listResp.Items {
        items = append(items, vpclatticetypes.Target{
            Id:   t.Id,
            Port: t.Port,
        })
    }
    _, err = r.svc.DeregisterTargets(ctx, &vpclattice.DeregisterTargetsInput{
        TargetGroupIdentifier: r.ARN,
        Targets:               items,
    })
    if err != nil { return err }
}
```

### `conditional`

Calls an API if a Go condition evaluates to true.

| Field | Type | Required | Description |
|---|---|---|---|
| `type` | enum | yes | Must be `conditional`. |
| `condition` | string | yes | Raw Go boolean expression (e.g. `r.Status == Active`). Inserted directly into an `if` statement. |
| `operation` | string | yes | SDK operation to call when the condition is true. |
| `input` | map[string]string | yes | Input params. Key = param name, value = resource struct field name. |

```yaml
preDeletion:
  - type: conditional
    condition: "r.PublishStatus == apigatewayv2types.PublishStatusPublished"
    operation: DisablePortal
    input:
      PortalId: PortalID
```

Generated code:
```go
if r.PublishStatus == apigatewayv2types.PublishStatusPublished {
    _, err := r.svc.DisablePortal(ctx, &apigatewayv2.DisablePortalInput{
        PortalId: r.PortalID,
    })
    if err != nil { return err }
}
```

---

## String Representation

Controls the generated `String()` method. Exactly one of `field`, `format`+`fields`, or `conditional` should be specified.

### Modes

| Field | Type | Required | Description |
|---|---|---|---|
| `field` | string | no | Single field name. Generates `return *r.<field>`. |
| `format` | string | no | `fmt.Sprintf` format string (e.g. `"%s (%s)"`). Must be paired with `fields`. |
| `fields` | []string | no | List of field names passed to `fmt.Sprintf`. Each is dereferenced as `*r.<field>`. |
| `conditional` | object | no | Conditional nil-check logic. See below. |

#### Simple field
```yaml
stringRepresentation:
  field: RelayName
```
Generates: `return *r.RelayName`

#### Format expression
```yaml
stringRepresentation:
  format: "%s (%s)"
  fields: [PortalID, Name]
```
Generates: `return fmt.Sprintf("%s (%s)", *r.PortalID, *r.Name)`

#### Conditional

| Field | Type | Required | Description |
|---|---|---|---|
| `field` | string | yes | Field to nil-check. |
| `ifNil` | string | yes | Field to return when `field` is nil. Generates `return *r.<ifNil>`. |
| `ifNotNil` | object | no | A nested `StringRepDef` (with `field` or `format`+`fields`) used when the field is non-nil. If omitted, returns `*r.<field>`. |

```yaml
stringRepresentation:
  conditional:
    field: Name
    ifNil: PortalID
    ifNotNil:
      format: "%s (%s)"
      fields: [PortalID, Name]
```

Generated code:
```go
func (r *MyResource) String() string {
    if r.Name != nil {
        return fmt.Sprintf("%s (%s)", *r.PortalID, *r.Name)
    }
    return *r.PortalID
}
```

### Override

Use `stringRepresentation.override` to replace the entire `String()` method body with raw Go code. When set, the declarative fields are ignored.

```yaml
stringRepresentation:
  override: |
    if r.Name != nil {
        return fmt.Sprintf("%s [%s]", *r.ID, *r.Name)
    }
    return *r.ID
```

---

## Extra Imports

A top-level field for additional import lines needed by any override blocks. Merged into the generated resource file's import block.

| Field | Type | Required | Description |
|---|---|---|---|
| `extraImports` | []string | no | Additional Go import lines. Use the full import syntax (e.g. `aliasname "github.com/aws/..."` or just `"github.com/aws/..."`). |

```yaml
extraImports:
  - apigatewayv2types "github.com/aws/aws-sdk-go-v2/service/apigatewayv2/types"
  - "strings"
```

---

## Override Summary

Overrides live in the section they're specific to. All override blocks must have balanced braces (validated by the tool).

| Override | YAML Location | Replaces |
|---|---|---|
| `list.override` | Under `list:` | Entire `List` method body |
| `delete.override` | Under `delete:` | Entire `Remove` method body (skips settings/pre-deletion logic) |
| `filterOverride` | Top-level | Entire `Filter` method body (ignores declarative `filters` list) |
| `stringRepresentation.override` | Under `stringRepresentation:` | Entire `String` method body |
| `extraImports` | Top-level | Adds imports needed by any override |

---

## Integration Test

Generates a `//go:build integration` test file with a testify suite containing `SetupSuite`, `TearDownSuite`, `TestList`, and `TestRemove`.

### Fields

| Field | Type | Required | Description |
|---|---|---|---|
| `integrationTest` | object | no | Omit entirely to skip integration test generation. |
| `integrationTest.create` | object | no | Defines the resource creation call for `SetupSuite`. |
| `integrationTest.create.operation` | string | yes | SDK create operation (e.g. `CreateRelay`). |
| `integrationTest.create.input` | map[string]any | yes | Input parameters. String values become `ptr.String("...")`. Strings containing `{{timestamp}}` become `ptr.String(fmt.Sprintf("...", time.Now().UnixNano()))` for unique naming. Boolean values become `aws.Bool(...)`. |

```yaml
integrationTest:
  create:
    operation: CreateRelay
    input:
      RelayName: "aws-nuke-test-{{timestamp}}"
```

Generated file structure:
```go
//go:build integration

type TestMailManagerRelaySuite struct {
    suite.Suite
    svc     *mailmanager.Client
    relayId *string
}

func (suite *TestMailManagerRelaySuite) SetupSuite()    { /* creates resource */ }
func (suite *TestMailManagerRelaySuite) TearDownSuite()  { /* deletes resource */ }
func (suite *TestMailManagerRelaySuite) TestList()        { /* lists and asserts > 0 */ }
func (suite *TestMailManagerRelaySuite) TestRemove()      { /* removes and asserts no error */ }

func TestMailManagerRelayIntegration(t *testing.T) {
    suite.Run(t, new(TestMailManagerRelaySuite))
}
```

---

## Agentic Workflow

The tool is designed to work in an automated loop where an agent (like Kiro) writes the YAML, generates code, checks for errors, and iterates until the output compiles and lints cleanly.

### Expected Agent Loop

1. The agent writes a `.yaml` DSL file to `resources/<service>-<resource>.yaml`.
2. The agent runs validation to catch structural problems early:
    ```console
    go run ./tools/codegen/ --validate resources/<service>-<resource>.yaml
    ```
    If validation fails, the agent reads the JSON error output from stdout, fixes the YAML, and retries.

3. The agent generates files with post-run checks:
    ```console
    go run ./tools/codegen/ --no-commit --output-manifest resources/<service>-<resource>.yaml
    ```
    Post-run steps (`goimports`, `go build`, `golangci-lint`) run automatically by default. Use `--no-build` to skip them.

4. If post-run fails (non-zero exit), the agent reads the error output from stderr. Common fixes:
    - Build errors from incorrect field types or missing imports → fix the YAML fields or add `overrides.extraImports`
    - Lint errors from unused variables or wrong patterns → adjust the DSL or add an `overrides` block
    - The agent edits the YAML and re-runs from step 2.

5. Once generation and post-run succeed, the agent generates and commits:
    ```console
    go run ./tools/codegen/ resources/<service>-<resource>.yaml
    ```

6. If the declarative DSL can't express the needed logic (complex conditionals, unusual API patterns), the agent adds `overrides` blocks with raw Go code and regenerates.

    Mock test files and integration test files are preserved by default during regeneration. Use `--force-mock-tests` or `--force-integration-tests` to regenerate them if the DSL changes require updated tests.

### Machine-Readable Output

All output is designed for programmatic consumption:

| Stream | Content |
|---|---|
| stdout | Generated file content (dry-run), JSON validation results (`--validate`), JSON manifest (`--output-manifest`), post-generation checklist |
| stderr | Diagnostic messages, warnings, build/lint error output |
| Exit code | 0 on success, non-zero on any error |

The `--output-manifest` flag produces JSON listing every file and what happened to it:

```json
{
  "files": [
    {"path": "resources/mailmanager-relay.go", "status": "created"},
    {"path": "resources/mailmanager.go", "status": "updated"},
    {"path": "resources/mailmanager_mock_test.go", "status": "skipped"}
  ]
}
```

The `--validate` flag produces JSON without generating any files:

```json
{
  "valid": false,
  "errors": [
    {"Field": "list.strategy", "Message": "list.strategy is required"}
  ]
}
```
