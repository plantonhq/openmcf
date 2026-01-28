# Spec Proto Authoring Guide

This document describes how to author or update `spec.proto` for a Planton resource kind.

## Folder and Naming
- Kind (PascalCase): e.g., AwsCloudFront, GcpPubsubTopic, AzureKeyVault
- Kind keyword (snake_case): aws_cloudfront, gcp_pubsub_topic, azure_key_vault
- Folder name (lowercase, no underscores): awscloudfront, gcppubsubtopic, azurekeyvault
- Path: `apis/project/planton/provider/<provider>/<kindfolder>/v1/spec.proto`

## Syntax and Package
- `syntax = "proto3";`
- `package org.openmcf.provider.<provider>.<kindfolder>.v1;`
- Do NOT include `go_package` in new proto files.

## Imports
- No validations in this step (do not import `buf/validate/validate.proto`).
- Optional when needed for value-or-reference fields:
  - `import "org/openmcf/shared/foreignkey/v1/foreign_key.proto";`

## Message Structure
- Define a single top-level message named `<Kind>Spec`.
- Keep messages and enums minimal; prefer clarity over completeness.

## Field Guidelines (80/20)
- Include only essential fields most users need.
- For cross-resource identifiers (IDs/ARNs such as IAM role ARN, KMS key ARN, security group IDs, subnet IDs, Route53 zone IDs, etc.), prefer the shared foreign key wrappers:
  - `org.openmcf.shared.foreignkey.v1.StringValueOrRef`
  - `org.openmcf.shared.foreignkey.v1.Int32ValueOrRef`
- When using `StringValueOrRef` for well-known kinds, you may set default hints using field options to improve Canvas wiring later:
  - `(org.openmcf.shared.foreignkey.v1.default_kind) = <CloudResourceKind>`
  - `(org.openmcf.shared.foreignkey.v1.default_kind_field_path) = "status.outputs.<field>"`
- Common examples for defaults:
  - Subnets: default_kind = AwsVpc
  - Security groups: default_kind = AwsSecurityGroup, default_kind_field_path = "status.outputs.security_group_id"
  - IAM role: default_kind = AwsIamRole, default_kind_field_path = "status.outputs.role_arn"
  - KMS key: default_kind = AwsKmsKey, default_kind_field_path = "status.outputs.key_arn"
- If a referenced kind does not yet exist (e.g., a future Lambda Layer), use plain `StringValueOrRef` without defaults.

## Enum Guidelines

When defining enums in spec.proto, follow these conventions for better user experience and cleaner manifests:

### Nesting
- **Always nest enums** inside the message where they are used
- Keep the full enum name (e.g., `KubernetesNamespaceBuiltInProfile`) for clarity in code
- Protobuf automatically namespaces nested enums to prevent collisions
- Reference as `MessageName.EnumName` in field definitions

Example:
```proto
message KubernetesNamespaceResourceProfile {
  enum KubernetesNamespaceBuiltInProfile {
    built_in_profile_unspecified = 0;
    small = 1;
    medium = 2;
    large = 3;
  }
  
  KubernetesNamespaceBuiltInProfile preset = 1;
}
```

### Value Naming
- **UNSPECIFIED values**: Use `lower_snake_case` with full enum prefix
  - Pattern: `{enum_name_in_snake_case}_unspecified`
  - Example: `built_in_profile_unspecified`, `service_mesh_type_unspecified`, `pod_security_standard_unspecified`
  - Rationale: Makes the zero value explicit and searchable
  
- **Other values**: Use lowercase without prefixes, minimal underscores
  - Single words: `small`, `medium`, `large`, `istio`, `linkerd`, `baseline`, `restricted`
  - Multiple words: Use words directly where clear, underscores only when necessary for clarity
  - Rationale: Clean YAML manifests (`preset: small` vs `preset: BUILT_IN_PROFILE_SMALL`)

### When NOT to Follow This Pattern
- Enums that represent external standards where uppercase is conventional (e.g., DNS record types: `A`, `AAAA`, `CNAME`)
- In these cases, add a comment explaining the deviation from the standard pattern

### Benefits
- **Cleaner user experience**: Manifests use `preset: small` instead of `preset: BUILT_IN_PROFILE_SMALL`
- **Better readability**: Lowercase values are easier to read and type
- **No collisions**: Protobuf nesting provides automatic namespacing
- **Consistent patterns**: All components follow the same enum style

## What to Avoid
- Do not add provider credentials here (those belong in stack input later).
- Avoid deep nesting unless essential.
- Keep comments brief and helpful.

## Example Skeleton (adapt kind/package)
```proto
syntax = "proto3";
package org.openmcf.provider.aws.awscloudfront.v1;

// Optional if you need value-or-ref wrappers
// import "org/openmcf/shared/foreignkey/v1/foreign_key.proto";

message AwsCloudFrontSpec {
  // Add essential 80/20 fields here
}
```

## Default Field Options

When a field should have a default value that OpenMCF applies if the user doesn't specify it:

### Requirements

1. **Mark the field as `optional`**: This generates a pointer type (`*string`) in Go with presence tracking
2. **Add the `(org.openmcf.shared.options.default)` field option**: Specifies the default value

### Import Required

```proto
import "org/openmcf/shared/options/options.proto";
```

### Syntax

```proto
// Container image repository.
// Default: ghcr.io/actions/actions-runner
optional string repository = 1 [(org.openmcf.shared.options.default) = "ghcr.io/actions/actions-runner"];

// Port number for the service.
// Default: 443
optional int32 port = 2 [(org.openmcf.shared.options.default) = "443"];
```

**Note:** Default values are always specified as strings, regardless of field type.

### Why Both Are Required

1. **`optional` keyword**: Enables field presence tracking in Go (generates `*string` vs `string`)
2. **Default field option**: Tells OpenMCF middleware what value to apply when field isn't set

### Build Enforcement

The custom linter `DEFAULT_REQUIRES_OPTIONAL` in `buf/lint/optional-linter` fails builds if a field has `(org.openmcf.shared.options.default)` but is NOT marked as `optional`.

### What NOT to Do

```proto
// WRONG: just a comment, no enforcement!
// Runner group name (defaults to "default" if not specified)
string runner_group = 7;
```

### Correct Pattern

```proto
// Runner group name.
// Default: default
optional string runner_group = 7 [(org.openmcf.shared.options.default) = "default"];
```

### Impact on IaC Modules

When fields become `optional`:
- Generated Go code changes from `string` to `*string`
- IaC modules must use getter methods: `spec.GetFieldName()` instead of `spec.FieldName`
- OpenMCF middleware guarantees defaults are applied before IaC modules run
- **No defensive coding needed** in IaC modules

## Notes
- Use official provider docs as reference while keeping the draft minimal.
- Validations (buf/validate + CEL) are added by a later rule.
